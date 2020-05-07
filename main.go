package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

type HtmlBody struct {
	Title string
	Body  []Place
}

type JsonBody struct {
	Places []Place
}

type Place struct {
	PlaceId     string
	PlaceName   string
	CountryId   string
	RegionId    string
	CityId      string
	CountryName string
}

func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Path[6:]
	url := "https://skyscanner-skyscanner-flight-search-v1.p.rapidapi.com/apiservices/autosuggest/v1.0/US/USD/en-US/?query=" + query

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("x-rapidapi-host", "skyscanner-skyscanner-flight-search-v1.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", goDotEnvVariable("API_KEY"))
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	jsonData, _ := ioutil.ReadAll(res.Body)
	var body *JsonBody
	err := json.Unmarshal(jsonData, &body)
	if err != nil {
		fmt.Println("error:", err)
	}

	displayHTML(w, query, body)
}

func displayHTML(w http.ResponseWriter, query string, body *JsonBody) {
	places := body.Places
	for i, place := range places {
		place.CountryId = place.CountryId[0 : len(place.CountryId)-4]
		place.CityId = place.CityId[0 : len(place.CityId)-4]
		places[i] = place
	}
	fmt.Println("places", places)
	page := &HtmlBody{
		Title: query,
		Body:  places,
	}
	t, _ := template.ParseFiles("view.html")
	t.Execute(w, page)
}

func main() {
	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
