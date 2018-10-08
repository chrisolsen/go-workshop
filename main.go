package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const kelvin = 273.15
const weatherAPIKey = "97a5f4ee1627098ed2f5df0c53f57568"
const weatherAPIUrl = "https://api.openweathermap.org/data/2.5/weather?q=%s,%s&appid=%s"

var port string

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
}

func main() {
	http.HandleFunc("/", weatherHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("failed to start server: %v", err.Error())
	}
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf(weatherAPIUrl, "Edmonton", "CA", weatherAPIKey)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed when fetching from weather API: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read API data: %v", err), http.StatusInternalServerError)
		return
	}

	var raw weatherData
	json.Unmarshal(data, &raw)

	weather := newMyWeather(raw)

	b, err := json.Marshal(weather)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal data: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(b))
}

// JSON Schema
// {
//   "coord": {
//      "lon":145.77,
// 	    "lat":-16.92
// 	 },
//   "weather":[
// 		{"id":803,"main":"Clouds","description":"broken clouds","icon":"04n"}
// 	 ],
//   "base":"cmc stations",
//   "main":{
// 		"temp":293.25,
//		"pressure":1019,
// 		"humidity":83,
// 		"temp_min":289.82,
// 		"temp_max":295.37
// 	 },
//   "wind":{
// 		"speed":5.1,"deg":150
// 	 },
//   "clouds":{
//		"all":75
// 	 },
//   "rain":{"3h":3},
//   "dt":1435658272,
//   "sys":{
// 		"type":1,"id":8166,"message":0.0166,"country":"AU","sunrise":1435610796,"sunset":1435650870
//   },
//   "id":2172797,
//   "name":"Cairns",
//   "cod":200
// }
type weatherData struct {
	Coord struct {
		Longitude float32 `json:"lon"`
		Latitude  float32 `json:"lat"`
	} `json:"coord"`

	Main struct {
		Temp     float32 `json:"temp"`
		Pressure float32 `json:"pressure"`
		Humidity float32 `json:"humidity"`
		TempMin  float32 `json:"temp_min"`
		TempMax  float32 `json:"temp_max"`
	}

	Weather []weatherDataWeather `json:"weather"`
}

type weatherDataWeather struct {
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type myWeather struct {
	Longitude   float32 `json:"lon"`
	Latitude    float32 `json:"lat"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Temp        float32 `json:"temp"`
	Pressure    float32 `json:"pressure"`
	Humidity    float32 `json:"humidity"`
	TempMin     float32 `json:"temp_min"`
	TempMax     float32 `json:"temp_max"`
	Icon        string  `json:"icon"`
}

func newMyWeather(data weatherData) myWeather {
	w := data.Weather[0]
	iconURL := fmt.Sprintf("http://openweathermap.org/img/w/%s.png", w.Icon)
	return myWeather{
		Longitude:   data.Coord.Longitude,
		Latitude:    data.Coord.Latitude,
		Title:       w.Main,
		Description: w.Description,
		Temp:        data.Main.Temp - kelvin,
		Pressure:    data.Main.Pressure,
		Humidity:    data.Main.Humidity,
		TempMin:     data.Main.TempMin - kelvin,
		TempMax:     data.Main.TempMax - kelvin,
		Icon:        iconURL,
	}
}
