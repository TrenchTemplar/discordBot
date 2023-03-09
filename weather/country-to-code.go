package weather

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CountryCode struct {
	Name string `json:"country_name"`
	ISO2 string `json:"ISO2"`
	ISO3 string `json:"ISO3"`
}

type CountryNameOnly struct {
	Name string `json:"country_name"`
}

func GetCountryCode(requestedCountry string) CountryCode {

	path := fmt.Sprintf("https://countrycode.dev/api/countries/%s", requestedCountry)

	response, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var codes []CountryCode
	err = json.NewDecoder(response.Body).Decode(&codes)
	if err != nil {
		log.Fatal(err)
	}

	countryCode := CountryCode{
		Name: codes[0].Name,
		ISO2: codes[0].ISO2,
		ISO3: codes[0].ISO3,
	}

	return countryCode
}

func GetCountryName(givenCode string) string {

	path := fmt.Sprintf("https://countrycode.dev/api/countries/iso2/%s", givenCode)

	response, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var countryResponse []CountryNameOnly
	err = json.NewDecoder(response.Body).Decode(&countryResponse)
	if err != nil {
		log.Fatal(err)
	}

	countryNameOnly := CountryNameOnly{
		Name: countryResponse[0].Name,
	}

	return countryNameOnly.Name
}
