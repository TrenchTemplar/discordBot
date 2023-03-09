package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type MainWeather struct {
	Main    TempVariables     `json:"main"`
	Weather []TempDescription `json:"weather"`
}

type TempDescription struct {
	Primary     string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type TempVariables struct {
	Temperature float32 `json:"temp"`
	RealFeel    float32 `json:"feels_like"`
	MinTemp     float32 `json:"temp_min"`
	MaxTemp     float32 `json:"temp_max"`
}

type ReturnData struct {
	Primary     string
	Description string
	Icon        string
	Temp        string
	RealFeel    string
	MinTemp     string
	MaxTemp     string
}

func GetWeather(lat float64, lon float64) ReturnData {

	godotenv.Load("./weather/.env")
	key := os.Getenv("WEATHER_API")

	path := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", lat, lon, key)

	response, err := http.Get(path)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	var tempWeather MainWeather
	err = json.NewDecoder(response.Body).Decode(&tempWeather)
	if err != nil {
		panic(err)
	}

	returnData := ReturnData{
		Primary:     tempWeather.Weather[0].Primary,
		Description: tempWeather.Weather[0].Description,
		Icon:        tempWeather.Weather[0].Icon,
		Temp:        strconv.FormatFloat(float64(tempWeather.Main.Temperature), 'f', 0, 32),
		RealFeel:    strconv.FormatFloat(float64(tempWeather.Main.RealFeel), 'f', 1, 32),
		MinTemp:     strconv.FormatFloat(float64(tempWeather.Main.MinTemp), 'f', 1, 32),
		MaxTemp:     strconv.FormatFloat(float64(tempWeather.Main.MaxTemp), 'f', 1, 32),
	}

	return returnData

}
