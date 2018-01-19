package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	resty "gopkg.in/resty.v1"

	"github.com/spf13/viper"
)

type Measurement struct {
	AirQualityIndex float32 `json:"airQualityIndex"`
	Humidity        float32 `json:"humidity"`
	Pm1             float32 `json:"pm1"`
	Pm10            float32 `json:"pm10"`
	Pm25            float32 `json:"pm25"`
	PollutionLevel  float32 `json:"pollutionLevel"`
	Pressure        float32 `json:"pressure"`
	Temperature     float32 `json:"temperature"`
	WindDirection   float32 `json:"windDirection"`
	WindSpeed       float32 `json:"windSpeed"`
}

type History struct {
	Measurement Measurement `json:"measurements"`
}

type Response struct {
	Measurement Measurement `json:"currentMeasurements"`
	History     []History   `json:"history"`
}

func main() {
	viper.SetDefault("neopixel.port", "/dev/ttyUSB0")
	viper.SetDefault("neopixel.baud", 9600)
	viper.SetDefault("sleep", 300)

	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.airly-neopixel")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	neopixel := newNeopixel(viper.GetString("neopixel.port"),
		viper.GetInt("neopixel.baud"))

	neopixel.open()

	time.Sleep(3 * time.Second)

	colorProvider := NewColorProvider(viper.GetFloat64("neopixel.brightness"))

	for {
		resp, err := resty.R().
			SetQueryParams(map[string]string{
				"latitude":  viper.GetString("latitude"),
				"longitude": viper.GetString("longitude"),
			}).
			SetHeader("apikey", viper.GetString("airly.apikey")).
			Get("https://airapi.airly.eu/v1/mapPoint/measurements")

		if err != nil {
			log.Fatal(err)
		}

		response := Response{}
		err = json.Unmarshal(resp.Body(), &response)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Current air quality: %f\n",
			response.Measurement.Pm25)

		neopixel.setColor(7,
			colorProvider.getColor(response.Measurement.Pm25))

		fmt.Println("Previous measurements:")

		historyLen := len(response.History)
		historyPixels := viper.GetInt("neopixel.pixels") - 1

		historyPixel := historyPixels - 1
		for i := historyLen - 1; i >= 0 && historyPixel >= 0; i-- {
			fmt.Printf("%f\n", response.History[i].Measurement.Pm25)
			color := colorProvider.getColor(response.History[i].Measurement.Pm25)
			neopixel.setColor(historyPixel, color)
			historyPixel--
		}

		fmt.Println()
		time.Sleep(time.Duration(viper.GetInt("sleep")) * time.Second)
	}
}
