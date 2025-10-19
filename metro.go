package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// metroNames Common names for route IDs.
var metroNames = map[string]string{
	"BLUE":   "901",
	"GREEN":  "902",
	"RED":    "903",
	"ORANGE": "904",
	"GOLD":   "905",
	"A":      "921",
	"B":      "922",
	"C":      "923",
	"D":      "924",
	"E":      "923",
	"F":      "923",
}

// Departure A specific departure for a train.
type Departure struct {
	Name          string
	StopName      string
	ShortStopName string
	When          string
	Direction     string
	Text          string
}

const (
	// NorthBound The direction ID for north.
	NorthBound = "0"

	// SouthBound the direction ID for south.
	SouthBound = "1"
)

// NextMetro Finds the next Departure for a specific route.
//
// route 		The route ID. Like "901" for the Blue Line.
// direction	The direction of the train, either "1" for South or "0" for North.
// stopId		The ID of the stop. Like "TF2" for Target Field Station 2.
// amount		The amount of departures to take, between 1 and 3.
func NextMetro(routeId string, direction string, stopId string, amount int) ([]*Departure, error) {
	if stopId == "" {
		return nil, errors.New("invalid stop ID")
	}

	if amount > 3 || amount < 1 {
		return nil, errors.New("invalid amount")
	}

	url := fmt.Sprintf("https://svc.metrotransit.org/nextrip/%s/%s/%s", routeId, direction, stopId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 7 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request NexTrip: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NexTrip API returned status %s", res.Status)
	}

	// nextrip response
	var payload struct {
		Departures []struct {
			RouteID        string `json:"route_id"`
			RouteShortName string `json:"route_short_name"`
			DirectionText  string `json:"direction_text"`
			DepartureText  string `json:"departure_text"`
			DepartureTime  int64  `json:"departure_time"`
			Description    string `json:"description"`
		} `json:"departures"`

		Stops []struct {
			Description string `json:"description"`
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// find the full stop name
	var stopName string
	for _, stop := range payload.Stops {
		stopName = stop.Description
	}

	// collect up to the next 3 upcoming departures, earliest first
	now := time.Now().Unix()
	var results []*Departure

	for _, departure := range payload.Departures {
		if departure.DepartureTime-now < 0 {
			continue
		}
		t := time.Unix(departure.DepartureTime, 0)
		formattedString := t.Format("2006-01-02 15:04 MST")

		d := &Departure{
			Name:          departure.RouteShortName,
			StopName:      stopName,
			ShortStopName: strings.ToUpper(stopId),
			When:          formattedString,
			Direction:     departure.DirectionText,
			Text:          departure.DepartureText,
		}
		results = append(results, d)
		if len(results) == amount {
			break
		}
	}

	if len(results) == 0 {
		return nil, errors.New("no upcoming Blue Line departures found at this stop")
	}

	return results, nil
}
