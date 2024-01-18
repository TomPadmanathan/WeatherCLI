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
				fmt.Println("Flag used multiple times")
				return
			}

			flags = append(flags, 'h')
			continue
		}

		if arg == "--temp" || arg == "-t" || arg == "--temperature" {
			if isFlagPresent('t', flags) {
				fmt.Println("Flag used multiple times")
				return
			}

			flags = append(flags, 't')
			continue
		}

		// Check if flag is forcast
		if len(arg) >= 11 {
			if arg[:11] == "--forecast=" {
				if isFlagPresent('f', flags) {
					fmt.Println("Flag used multiple times")
					return
				}
				flags = append(flags, 'f')
				number, err := strconv.Atoi(arg[11:])
				forecast = number
				if err != nil {
					fmt.Println("Invalid forcast flag")
					return
				}
				continue
			}
		}
		if len(arg) >= 3 {
			if arg[:3] == "-f=" {
				if isFlagPresent('f', flags) {
					fmt.Println("Flag used multiple times")
					return
				}
				flags = append(flags, 'f')
				number, err := strconv.Atoi(arg[3:])
				forecast = number
				if err != nil {
					fmt.Println("Flag used multiple times")
					return
				}
				continue
			}
		}

		fmt.Println("Invalid flag")
		return
	}

	if isFlagPresent('h', flags) {
		fmt.Print("Get Weather CLI\n\n")
		fmt.Print("To use this CLI use \"getweather [flags] [location]\"\n\n")

		fmt.Println("Flags:")
		fmt.Print("	--help\n	-h\n	Used to open this help screen.\n\n")
		fmt.Print("	--temp\n	-t\n	Used to get just the temperature for a location.\n\n")
		fmt.Print("	--forecast\n	-f\n	Used to get the forecast for up to the next 3 days. E.g: --forecast=2\n\n")
		fmt.Println("View the source code for this CLI:")
		fmt.Println("https://github.com/TomPadmanathan/WeatherCLI")
		return
	}

	// If no location provided set to England
	if len(location) == 0 {
		location = "London"
	}

	if forecast > 3 {
		fmt.Println("Forecast is larger than limit (3)")
		return
	}

	var forecastString string = "current"

	if forecast > 0 {
		forecastString = "forecast"
	}

	var url string = fmt.Sprintf("https://api.weatherapi.com/v1/%s.json?key=%s&q=%s&days=%v&aqi=no&alerts=no", forecastString, os.Getenv("WeatherApiKey"), location, forecast)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Something went wrong fetching data")
		return
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
		fmt.Println("Something went wrong fetching data")
		return
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println("Something went wrong fetching data")
		return
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
