package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
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
			Hour struct {
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
	err := godotenv.Load()
	apiKey := os.Getenv("API_KEY")

	response, err := http.Get("http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=London")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		fmt.Println("Error: Received non-OK HTTP status:", response.Status)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// convert slice of byte to string
	// fmt.Println(string(body))

	var weather Weather
	json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}
	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("%s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)
	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		fmt.Printf("%s: %.0fC, %s (%.0f%%)\n",
			date.Format("15:04"),
			hour.TempC,
			hour.Condition.Text,
			hour.ChanceOfRain,
		)
	}
}
