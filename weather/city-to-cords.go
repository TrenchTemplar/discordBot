package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Cordinates struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

func GetCords(requestedLocation string) (string, float64, float64, string) {

	godotenv.Load("./weather/.env")
	key := os.Getenv("WEATHER_API")

	path := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", requestedLocation, key)

	response, err := http.Get(path)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	var cords []Cordinates
	err = json.NewDecoder(response.Body).Decode(&cords)
	if err != nil {
		panic(err)
	}

	name := cords[0].Name
	lat := cords[0].Lat
	lon := cords[0].Lon
	country := cords[0].Country

	return name, lat, lon, country
}
