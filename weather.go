package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WeatherData struct {
	Weather []Weather
	Main    mainWeather
	Wind    wind
	Sys     sys
	Name    string
	Cod     int
	Message string
}

type Weather struct {
	Description string
}

type mainWeather struct {
	Temp    float32
	TempMin float32 `json:"temp_min"`
	TempMax float32 `json:"temp_max"`
}

type wind struct {
	Speed float32
}

type sys struct {
	Country string
}

func GetWeatherData(cityName, ownToken string) (*WeatherData, error) {
	request := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?units=metric&lang=en&q=%s&APPID=%s", cityName, ownToken)
	_data := WeatherData{}

	resp, err := http.Get(request)
	if err != nil {
		return nil, err
	}

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	err = json.Unmarshal(buffer, &_data)
	if err != nil {
		return nil, err
	}

	return &_data, nil
}

func GetWeather(cityName, ownToken, lang string) (string, error) {
	data, err := GetWeatherData(cityName, ownToken)
	if err != nil {
		return "", err
	}

	if lang == FlagUa {
		return fmt.Sprintf("ім'я: %s\nкраїна: %s\nтемпература: %.3f\nопис: %s\nмін. температура: %.3f\nмакс. температура %.3f\nшвидкість вітру: %.3f",
			data.Name, data.Sys.Country, data.Main.Temp, data.Weather[0].Description, data.Main.TempMin, data.Main.TempMax, data.Wind.Speed), nil
	} else if lang == FlagRu {
		return fmt.Sprintf("имя: %s\nстрана: %s\nтемпература: %.3f\nописание: %s\nмин. температура: %.3f\nмакc. температура %.3f\nскорость ветра: %.3f",
			data.Name, data.Sys.Country, data.Main.Temp, data.Weather[0].Description, data.Main.TempMin, data.Main.TempMax, data.Wind.Speed), nil
	} else {
		return fmt.Sprintf("name: %s\ncountry: %s\ntemp: %.3f\ndescription: %s\nmin temp: %.3f\nmax temp %.3f\nwind speed: %.3f",
			data.Name, data.Sys.Country, data.Main.Temp, data.Weather[0].Description, data.Main.TempMin, data.Main.TempMax, data.Wind.Speed), nil
	}
}
