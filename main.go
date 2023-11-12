package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
)

type ResponseStruct struct {
	Data []Breeds `json:"data"`
}

type Breeds struct {
	Breed   string `json:"breed"`
	Country string `json:"country"`
	Origin  string `json:"origin"`
	Coat    string `json:"coat"`
	Pattern string `json:"pattern"`
}

func main() {
	resp, err := http.Get("https://catfact.ninja/breeds")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	var res ResponseStruct
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		log.Fatal(err)
	}

	breedsByCountry := make(map[string][]Breeds)
	for _, breed := range res.Data {
		breedsByCountry[breed.Country] = append(breedsByCountry[breed.Country], breed)
	}

	for _, breeds := range breedsByCountry {
		sort.Slice(breeds, func(i, j int) bool {
			return len(breeds[i].Breed) < len(breeds[j].Breed)
		})
	}

	out, err := json.Marshal(breedsByCountry)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("out.json", out, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
