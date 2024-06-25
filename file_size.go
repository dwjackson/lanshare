package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func parseFileSize(input string) (int64, error) {
	runes := []rune(input)
	n, i, err := parseNumber(runes)
	if err != nil {
		return 0, err
	}
	if i < 0 || i >= len(runes) {
		return n, nil
	}
	unit := string(runes[i:])
	shift, err := parseShift(unit)
	if err != nil {
		return 0, err
	}
	bytes := n << shift
	return bytes, nil
}

func parseNumber(runes []rune) (int64, int, error) {
	var digits []rune
	var index int = -1
	for i, r := range runes {
		if unicode.IsDigit(r) {
			digits = append(digits, r)
		} else if r == '.' {
			return 0, -1, errors.New("Fractional sizes not allowed")
		} else {
			index = i
			break
		}
	}
	if len(digits) == 0 {
		return 0, -1, errors.New("No digits in file size number")
	}
	s := string(digits)
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, -1, err
	}
	return int64(n), index, err
}

func parseShift(unit string) (int, error) {
	units := map[string]int{
		"KiB": 10,
		"MiB": 20,
		"GiB": 30,
	}
	shift, exists := units[unit]
	if !exists {
		var allowedUnits []string
		for key := range units {
			allowedUnits = append(allowedUnits, key)
		}
		allowedUnitsStr := strings.Join(allowedUnits, ",")
		msg := fmt.Sprintf("Unrecognized file size unit: %s; allowed units: %s", unit, allowedUnitsStr)
		return 0, errors.New(msg)
	}
	return shift, nil
}
