package main

import (
	"errors"
	"regexp"
	"strings"
)

// Formatting Default formats for the response
var Formatting = map[int]string{
	0: "There's a %DIR %NAME %TYPE coming to %STOP %TIME.",
	1: "%DIR %NAME %TYPE at %STOP %TIME.",
	2: "%SHORT_DIR %COLOR %TYPE: %SHORT_STOP %TIME.",
	3: "%COLOR %SHORT_DIR @ %SHORT_STOP %TIME",
}

// CircleEmojis Emojis that represent various Metro lines.
var CircleEmojis = map[string]string{
	"gold":   "🟡",
	"orange": "🟠",
	"red":    "🔴",
	"green":  "🟢",
	"blue":   "🔵",
}

// Variables Possible variables that can be in a given format.
var Variables = map[string]*regexp.Regexp{
	"name":       regexp.MustCompile("(?i)%NAME"),
	"dir":        regexp.MustCompile("(?i)%DIR"),
	"type":       regexp.MustCompile("(?i)%TYPE"),
	"stop":       regexp.MustCompile("(?i)%STOP"),
	"time":       regexp.MustCompile("(?i)%time"),
	"short_dir":  regexp.MustCompile("(?i)%SHORT_DIR"),
	"short_stop": regexp.MustCompile("(?i)%SHORT_STOP"),
	"color":      regexp.MustCompile("(?i)%COLOR"),
}

// IsProperFormat Find if a given format number is proper by
// checking if it's in Formatting.
//
// formatType — An ID for a format.
//
// Returns: If the included formatType is in Formatting
func IsProperFormat(formatType int) bool {
	_, ok := Formatting[formatType]

	return ok
}

// Format a departure based on a given type ID.
//
// formatType — The type of format to take from Formatting
// data — The Departure pointer.
//
// Returns:
// - a pointer to a string with the formatted response
// - an error if there's an issue finding the format
func Format(formatType int, data *Departure) (*string, error) {
	var format, foundFormat = Formatting[formatType]
	if !foundFormat {
		return nil, errors.New("unknown format type")
	}

	for key, re := range Variables {
		var value string

		switch key {
		// the name of the metro line
		case "name":
			value = data.Name
			break

		// the short ID of the stop
		case "short_stop":
			value = data.ShortStopName
			break

		// nb or sb
		case "short_dir":
			value = data.Direction
			break

		// the type of metro (train, bus)
		case "type":
			lowerName := strings.ToLower(data.Name)

			if lowerName == "blue" || lowerName == "green" {
				value = "Train"
			} else {
				value = "Bus"
			}
			break

		// the stop that it's finding the next for.
		case "stop":
			value = strings.ToLower(data.StopName)
			break

		// the color of the metro line (only for lrt or brt)
		case "color":
			emoji, ok := CircleEmojis[strings.ToLower(data.Name)]

			if ok {
				value = emoji
			} else {
				value = "N/A"
			}

			break

		// the amount of time til the metro leaves
		case "time":
			if strings.Contains(data.Text, "Min") {
				value = "in " + data.Text
			} else {
				value = "at " + data.Text
			}
			break

		// direction of the metro
		case "dir":
			switch data.Direction {
			case "NB":
				value = "northbound"
				break
			case "SB":
				value = "southbound"
				break
			case "EB":
				value = "eastbound"
				break
			case "WB":
				value = "westbound"
				break
			}

			break
		}

		format = re.ReplaceAllString(format, value)
	}

	return &format, nil
}
