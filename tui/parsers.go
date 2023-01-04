package tui

import (
	"regexp"
	"strconv"
)

func defaultParser(s string) string {
	return s
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
	isGram, _ := regexp.MatchString("g$", s)
	if len(s) == 1 || (len(s) == 3 && isMol) {
		if isGram {
			return 0.001, "kg" // kg - not gram - is SI
		}
		return 1, s
	} else if len(s) == 2 || isMol {
		hasUnit, _ := regexp.MatchString("[kcmunpf]", s[0:1])
		if hasUnit && !isMol && !isGram {
			return units[s[0:1]], s[1:2]
		} else if isMol {
			return units[s[0:1]], "mol"
		} else if isGram {
			return units[s[0:1]] / 1000, "kg" // kg - not gram - is SI
		}
	}
	// Checks that should be done
	// Should start with f p n u m c k if multi char
	// UNLESS the multichar ends with 'mol', which is a single unit (not to be confused with M, molar, which is mol/L)
	// If not multi char, should be m, M, g, L
	return 1, s
}

func parseTime(s string) int32 {
	timeUnits := make(map[string]int32)
	timeUnits["s"] = 1
	timeUnits["m"] = timeUnits["s"] * 60
	timeUnits["h"] = timeUnits["m"] * 60
	timeUnits["d"] = timeUnits["h"] * 24
	timeUnits["w"] = timeUnits["d"] * 7
	timeUnits["y"] = timeUnits["d"] * 365

	times := splitBefore(s, regexp.MustCompile("[0-9]*"))
	valR := regexp.MustCompile("^[0-9]*")
	unitR := regexp.MustCompile("[a-z]$")
	secs := int32(0)
	for _, v := range times {
		// Need to make sure something like d12m is INVALID
		// 1dm should also be invalid
		// 1d12m is fine tho
		val, _ := strconv.Atoi(valR.FindString(v))
		unit := unitR.FindString(v)
		secs += int32(val) * timeUnits[unit]
	}
	return secs
}
