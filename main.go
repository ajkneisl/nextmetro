package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// main Starts HTTP server.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		// assuming /metro/name/stop/dir
		if len(parts) != 4 {
			http.Error(w, "Usage: /{name}/{stop}/{north|south|east|west}", http.StatusBadRequest)
			return
		}

		// the format of the departure
		format, err := strconv.Atoi(r.URL.Query().Get("format"))
		if err != nil || !IsProperFormat(format) {
			format = 0
		}

		// the amount of departures to return
		amount, err := strconv.Atoi(r.URL.Query().Get("amount"))
		if err != nil {
			amount = 1
		}

		name := parts[1]

		// check if it's one of the possible shortened names
		metroName, foundName := metroNames[strings.ToUpper(parts[1])]
		if foundName {
			name = metroName
		}

		stop := parts[2]
		directionStr := parts[3]

		// the direction
		direction := SouthBound
		if strings.EqualFold(strings.ToLower(directionStr), "north") || strings.EqualFold(strings.ToLower(directionStr), "east") {
			direction = NorthBound
		}

		departures, err := NextMetro(name, direction, stop, amount)
		if err != nil && !strings.Contains(err.Error(), "departures") {
			http.Error(w, fmt.Sprintf("There was an issue. Please double check the name and direction."), http.StatusInternalServerError)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("No upcoming departures."), http.StatusOK)
			return
		}

		w.Header().Set("Content-Type", "text/plain")

		var response = ""

		for _, dep := range departures {
			var departureResponse, formatError = Format(format, dep)

			if formatError != nil {
				http.Error(w, fmt.Sprintf("There was an issue with your format."), http.StatusInternalServerError)
				return
			}

			response += *departureResponse + "\n"
		}

		_, err = fmt.Fprint(w, response)
		if err != nil {
			return
		}
	})

	fmt.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
