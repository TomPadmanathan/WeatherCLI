package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
	var forecast int
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
			if isFlagPresent('h', flags) {
				panic("flag used multiple times")
			}

			flags = append(flags, 'h')
			continue
		}

		if arg == "--temp" || arg == "-t" || arg == "--temperature" {
			if isFlagPresent('t', flags) {
				panic("flag used multiple times")
			}

			flags = append(flags, 't')
			continue
		}

		// Check if flag is forcast
		if len(arg) >= 11 {
			if arg[:11] == "--forecast=" {
				if isFlagPresent('f', flags) {
					panic("flag used multiple times")
				}
				flags = append(flags, 'f')
				number, err := strconv.Atoi(arg[11:])
				forecast = number
				if err != nil {
					panic("Invalid Forcast Flag")
				}
				continue
			}
		}
		if len(arg) >= 3 {
			if arg[:3] == "-f=" {
				if isFlagPresent('f', flags) {
					panic("flag used multiple times")
				}
				flags = append(flags, 'f')
				number, err := strconv.Atoi(arg[3:])
				forecast = number
				if err != nil {
					panic("Invalid Forcast Flag")
				}
				continue
			}
		}

		panic("Invalid flag")
	}

	if isFlagPresent('h', flags) {
		fmt.Println("help")
		return
	}

	// If no location provided set to England
	if len(location) == 0 {
		location = "England"
	}

	if forecast > 3 {
		panic("Forcast is larger 3")
	}

	var forecastString string = "current"

	if forecast > 0 {
		forecastString = "forecast"
	}

	var url string = fmt.Sprintf("https://api.weatherapi.com/v1/%s.json?key=%s&q=%s&days=%v&aqi=no&alerts=no", forecastString, os.Getenv("WeatherApiKey"), location, forecast)

	res, err := http.Get(url)
	if err != nil {
		panic("Weather API not available")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 400 {
			fmt.Println("Location provided is invalid")
			return
		}
		if res.StatusCode == 403 {
			fmt.Println("Api key is invalid")
			return
		}
		fmt.Println("Something went wrong fetching data")
		return
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
		fmt.Printf("Current temperature in %s, %s is: %s°c \n\n", weather.Location.Name, weather.Location.Country, fmt.Sprintf("%.1f", weather.Current.TempC))

		if isFlagPresent('f', flags) {
			for index, day := range weather.Forecast.Forecastday {
				fmt.Printf("\nDay %v:\n", index+1)
				for index, hour := range day.Hour {
					fmt.Printf("Hour %v temperature: %s°c\n", index+1, fmt.Sprintf("%.1f", hour.TempC))
				}
			}
		}
		return
	}

	fmt.Printf("Weather in %s, %s:\n\nCurrent Temperature: %s°c\nCurrent Weather Condition: %s\n", weather.Location.Name, weather.Location.Country, fmt.Sprintf("%.1f", weather.Current.TempC), weather.Current.Condition.Text)

	if forecast > 0 {
		fmt.Print("\nForcast:\n\n")
	}
	for index, day := range weather.Forecast.Forecastday {
		fmt.Printf("Day %v:\n", index+1)
		for index, hour := range day.Hour {
			fmt.Printf("Hour %v:\n", index+1)
			fmt.Printf("Temperature: %s°c\nWeather Condition: %s\nChance of rain: %s\n", fmt.Sprintf("%.1f", hour.TempC), hour.Condition.Text, fmt.Sprintf("%.1f", hour.ChanceOfRain))
			fmt.Println()
		}
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
