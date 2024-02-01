package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetLocationAreas(url string) (string, string) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode > 299 {
		log.Fatal("Response has a non 2xx status code")
	}
	if err != nil {
		log.Fatal(err)
	}

	results := PokeResponse[LocationArea]{}
	marshalErr := json.Unmarshal(body, &results)
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	for _, r := range results.Results {
		fmt.Println(r.Name)
	}

	if results.Previous == nil {
		return *results.Next, ""
	}
	if results.Next == nil {
		return "", *results.Previous
	}

	return *results.Next, *results.Previous
}
