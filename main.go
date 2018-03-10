package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var OpenWeatherApiKey = ""

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/", weatherHandler)
	fmt.Println("Starting server...")
	fmt.Println("The Server lives! [CTRL+C to exit]")
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func weatherIn(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + OpenWeatherApiKey + "&q=" + city)
	if err != nil {
		fmt.Println(err)
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var responseWeather weatherData

	if err := json.NewDecoder(resp.Body).Decode(&responseWeather); err != nil {
		fmt.Println(err)
		return weatherData{}, err
	}

	return responseWeather, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := strings.SplitN(r.URL.Path, "/", 3)[2]

	data, err := weatherIn(city)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}
