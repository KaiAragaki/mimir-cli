package tui

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/KaiAragaki/mimir-cli/shared"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

// --- STRUCTS ---

// Entry
type Entry struct {
	fields  []field // The fields
	focused int     // Which field is focused
	ok      bool    // Are all entries valid?
	repo    *gorm.DB
	subErr  string // What error (if any) came from submitting to DB?
}

// Field - a single unit of an entry
type field struct {
	displayName string // What the header of the field will be displayed as
	input       textarea.Model
	hasErr      bool
	errMsg      string
	vfuns       []func(s string) (string, bool)
}

// Sensible defaults for fields
func NewDefaultField() field {
	ta := textarea.New()
	ta.FocusedStyle = textAreaFocusedStyle
	ta.ShowLineNumbers = false
	ta.Prompt = " "
	ta.BlurredStyle = textAreaBlurredStyle
	//func(s string) (string, bool) { return "", false }
	var fns []func(s string) (string, bool)

	return field{
		input:  ta,
		hasErr: true,
		errMsg: "",
		vfuns:  fns,
	}
}

// --- STYLING ---

const (
	white     = lipgloss.Color("#FFFFFF")
	purple    = lipgloss.Color("#7f12c7")
	darkGray  = lipgloss.Color("#767676")
	red       = lipgloss.Color("#FF0000")
	green     = lipgloss.Color("#00FF00")
	lightBlue = lipgloss.Color("#5C8DFF")
	blue      = lipgloss.Color("#3772FF")
	yellow    = lipgloss.Color("#FDCA40")
	black     = lipgloss.Color("#000000")
)

var (
	activeInputStyle     = lipgloss.NewStyle().Foreground(white).Background(purple)
	inactiveInputStyle   = lipgloss.NewStyle().Foreground(purple)
	continueStyle        = lipgloss.NewStyle().Foreground(darkGray)
	cursorStyle          = lipgloss.NewStyle().Foreground(white)
	cursorLineStyle      = lipgloss.NewStyle().Background(lipgloss.Color("57")).Foreground(lipgloss.Color("230"))
	errorStyle           = lipgloss.NewStyle().Foreground(darkGray).Italic(true)
	okStyle              = lipgloss.NewStyle().Foreground(green)
	textAreaFocusedStyle = textarea.Style{
		Base: lipgloss.
			NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(yellow).
			BorderLeft(true).
			Foreground(yellow),
	}
	textAreaBlurredStyle = textarea.Style{
		Base: lipgloss.
			NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(white).
			BorderLeft(true).
			Foreground(white),
	}

	headerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(white).
			BorderLeft(true)

	activeHeaderStyle = headerStyle.Copy().
				Foreground(yellow).Bold(true).
				BorderForeground(yellow)

	titleStyle = lipgloss.NewStyle().
			Background(yellow).
			Foreground(black).
			Margin(0, 2, 3, 2)
	docStyle = lipgloss.NewStyle().Margin(1)
)

func newTextarea() textarea.Model {
	t := textarea.New()
	return t
}

// Both Action and Table share the same item structure, so it's defined here
type item struct {
	title, desc string
}

func (i item) Title() string {
	return i.title
}

func (i item) Description() string {
	return i.desc
}

func (i item) FilterValue() string {
	return i.title
}

// A function calls the correct Init* function based on the table name selected
// I'm sure there's a better way to do this (generics?) but I'm too dumb
func InitForm(tableName string) tea.Model {
	switch tableName {
	case "Cell":
		return InitCell()
	}
	return InitTable(shared.Table)
}

func noFieldHasError(c Entry) bool {
	for _, v := range c.fields {
		if v.hasErr {
			return false
		}
	}
	return true
}

func getEntryStatus(c Entry) string {
	if noFieldHasError(c) {
		return okStyle.Render("Lookin' good!")
	} else {
		return errorStyle.Render("Entry not ready to be submitted.")
	}
}

func Validate(c *Entry) {
	for i, v := range c.fields {
		for _, w := range v.vfuns {
			c.fields[i].errMsg, c.fields[i].hasErr = w(v.input.Value())
			if c.fields[i].hasErr {
				break
			}
		}
	}
}

// Validators

func valIsBlank(s string) (string, bool) {
	if s == "" {
		return "Field must not be blank", true
	}
	return "", false
}

func valIsntLcAndNum(s string) (string, bool) {
	lcAndNum := regexp.MustCompile("^[a-z0-9]*$")
	if !lcAndNum.MatchString(s) {
		return "May only include numbers and lowercase letters", true
	}
	return "", false
}

func valIsntLcNumUnderDash(s string) (string, bool) {
	lcNumUnderScoreDash := regexp.MustCompile("^[a-z0-9_-]*$")
	if !lcNumUnderScoreDash.MatchString(s) {
		return "May only include numbers, lowercase letters, underscores, and dashes", true
	}
	return "", false
}

func valIsntNum(s string) (string, bool) {
	num := regexp.MustCompile("^[0-9]*$")
	if !num.MatchString(s) {
		return "May only include numbers", true
	}
	return "", false
}

// parseUnits should take a string (like, say, "40mg/mL") and convert it
func parseUnits(s string) (float32, string) {
	rAmt := regexp.MustCompile("^[0-9]*")
	rUnits := regexp.MustCompile("[a-zA-Z/]*$") // Slash allows for "mg/mL"

	amt := rAmt.FindString(s)
	amtFloat, _ := strconv.ParseFloat(amt, 32)
	units := rUnits.FindString(s)
	// Checks that should be done:
	// Numbers should preceed units
	// units should follow numbers

	isFrac, _ := regexp.MatchString("/", units)
	if isFrac {
		rNumerator := regexp.MustCompile("[a-zA-Z]+")
		rDenominator := regexp.MustCompile("[^/]*$")
		numerator := rNumerator.FindString(s)
		denominator := rDenominator.FindString(s)
		numVal, numUnit := makeSI(numerator)
		denomVal, denomUnit := makeSI(denominator)
		return numVal * float32(amtFloat) / denomVal, numUnit + "/" + denomUnit
	} else {
		val, unit := makeSI(units)
		return val * float32(amtFloat), unit
	}
}

func makeSI(s string) (float32, string) {
	units := make(map[string]float32)
	units["k"] = 1000
	units["c"] = 0.01
	units["m"] = 0.001
	units["u"] = units["m"] / 1000
	units["n"] = units["u"] / 1000
	units["p"] = units["n"] / 1000
	units["f"] = units["p"] / 1000

	isMol, _ := regexp.MatchString("mol$", s)
	if len(s) == 1 || (len(s) == 3 && isMol) {
		return 1, s
	} else if len(s) == 2 || isMol {
		hasUnit, _ := regexp.MatchString("[kcmunpf]", s[0:1])
		if hasUnit && !isMol {
			return units[s[0:1]], s[1:2]
		} else if isMol {
			return units[s[0:1]], "mol"
		} else {
			fmt.Errorf("Unknown unit multiplier")
		}
	} else {
		fmt.Errorf("Unknown unit")
	}
	// Checks that should be done
	// Should start with f p n u m c k if multi char
	// UNLESS the multichar ends with 'mol', which is a single unit (not to be confused with M, molar, which is mol/L)
	// If not multi char, should be m, M, g, L
	return 1, s
}
