package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
)

func (client *PokeClient) GetLocationAreas(url string) (PokeResponse[LocationArea], error) {
	if val, ok := client.cache.Get(url); ok {
		resp := PokeResponse[LocationArea]{}
		err := json.Unmarshal(val, &resp)
		if err != nil {
			log.Fatal(err)
		}
		return resp, nil
	}

	res, err := client.httpClient.Get(url)
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

	response := PokeResponse[LocationArea]{}
	marshalErr := json.Unmarshal(body, &response)
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	client.cache.Add(url, body)

	return response, nil
}

func (client *PokeClient) GetLocationArea(location string) (PokeArea, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + location

	if val, ok := client.cache.Get(location); ok {
		resp := PokeArea{}
		err := json.Unmarshal(val, &resp)
		if err != nil {
			log.Fatal(err)
		}
		return resp, nil
	}

	res, err := client.httpClient.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode == 404 {
		return PokeArea{}, errors.New("This area does not exist!")
	}

	if res.StatusCode > 299 {
		log.Fatal("Response has a non 2xx status code")
	}
	if err != nil {
		log.Fatal(err)
	}

	response := PokeArea{}
	marshalErr := json.Unmarshal(body, &response)
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	client.cache.Add(location, body)

	return response, nil
}
