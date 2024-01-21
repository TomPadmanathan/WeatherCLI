# Weather CLI
### Overview

Get Weather CLI is a simple command-line interface (CLI) written in Golang that allows you to retrieve weather information for a specified location. Whether you need the current temperature or a forecast for the upcoming days, this CLI has you covered.
### Installation

To use Get Weather CLI, you need to have Golang installed on your machine. Once Golang is installed, you can install the CLI with the following command:


```bash
go get -u github.com/TomPadmanathan/WeatherCLI
```

You will also need to create an api key from [weatherapi.com](https://www.weatherapi.com/) and set it in your system environment variables as: `WeatherApiKey`.

### Usage

After installation, you can run the CLI using the following command:

```bash
getweather [flags] [location]
```
### Flags

    --help or -h: Opens the help screen, providing information on how to use the CLI.
    --temp or -t: Retrieves only the current temperature for the specified location.
    --forecast or -f: Retrieves the forecast for up to the next 3 days. You can specify the number of days, for example, --forecast=2.

### Examples

To get the current temperature for a location:

```bash
getweather -t New York
```

To get the forecast for the next 2 days for a location:

```bash
getweather -f=2 London
```
### Source Code

You can view the source code for Get Weather CLI on GitHub. Feel free to explore the code and contribute if you find ways to enhance the CLI.
Issues and Contributions

If you encounter any issues or have suggestions for improvement, please open an issue on the GitHub repository. Contributions are welcome!

Enjoy using Get Weather CLI to stay updated on the weather conditions in your desired locations!
