package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	var args []string = os.Args
	var location string
	var flags []rune
	for index, arg := range args {
		if index == 0 {
			continue
		}

		// location
		if index == len(args)-1 {
			var isFlag bool

			// Check if first arg has flag
			for _, char := range arg {
				if char == '-' {
					isFlag = true
				}
			}

			if !isFlag {
				location = args[len(args)-1]
				continue
			}
		}

		if arg == "--help" || arg == "-h" {
			for _, flag := range flags {
				if flag == 'h' {
					panic("flag used multiple times")
				}
			}

			flags = append(flags, 'h')
			continue
		}

		if arg == "--temp" || arg == "-t" || arg == "--temperature" {
			for _, flag := range flags {
				if flag == 't' {
					panic("flag used multiple times")
				}
			}

			flags = append(flags, 't')
			continue
		}

		panic("Invalid flag")
	}

	// If no location provided set to England
	if len(location) == 0 {
		location = "England"
	}

	var url string = fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&days=1&aqi=no&alerts=no", os.Getenv("WeatherApiKey"), location)

	res, err := http.Get(url)
	if err != nil {
		panic("Weather API not available")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	if isFlagPresent('t', flags) {
		fmt.Printf("Current temperature in %s is: %f°c \n", weather.Location.Name, weather.Current.TempC)
	}
}

func isFlagPresent(flag rune, flags []rune) bool {
	for _, flag2 := range flags {
		if flag == flag2 {
			return true
		}
	}
	return false
}
