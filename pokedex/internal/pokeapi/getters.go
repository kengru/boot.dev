package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func (client *PokeClient) GetLocationAreas(url string) (PokeResponse[LocationArea], error) {
	if val, ok := client.cache.Get(url); ok {
		resp := PokeResponse[LocationArea]{}
		err := json.Unmarshal(val, &resp)
		if err != nil {
			fmt.Println("it was here")
			fmt.Println(val)
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
