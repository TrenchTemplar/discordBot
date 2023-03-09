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

func GetCountryCode(requestedCountry string) CountryCode {
	fmt.Println("request", "!"+requestedCountry)
	path := fmt.Sprintf("https://countrycode.dev/api/countries/%s", requestedCountry)

	response, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var codes CountryCode
	err = json.NewDecoder(response.Body).Decode(&codes)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()

	countryCode := CountryCode{
		Name: codes.Name,
		ISO2: codes.ISO2,
		ISO3: codes.ISO3,
	}

	return countryCode
}

// type ReturnCode struct {
// 	Name    string  `json:"country_name"`
// 	ISO2    string `json:"ISO2"`
// 	ISO3    string `json:"ISO3"`

// }
