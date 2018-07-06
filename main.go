package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kr/pretty"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"io"
	"log"
	"os"
)

func EnvLoad() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Spot struct {
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

func main() {

	EnvLoad()

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GEO_CODING_API_KEY")))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	flag.Parse()

	inputFile, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal("Error:", err)
	}

	defer inputFile.Close()
	reader := csv.NewReader(inputFile)

	var spots []Spot

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else {
			if err != nil {
				log.Fatal("Error:", err)
			}
		}

		var searchName = record[0]

		var r = createGeocodingRequest(searchName, "ja")

		lat, lng, address, shortName := searchAddress(c, r)

		spot := Spot{Name: searchName, Address: address, Lat: lat, Lng: lng}
		spots = append(spots, spot)

		pretty.Println(lat, lng, address, shortName)

	}

	writeFile("spot.json", toJson(spots))
}

func createGeocodingRequest(address string, language string) *maps.GeocodingRequest {
	return &maps.GeocodingRequest{
		Address:  address,
		Language: language,
	}
}

func searchAddress(client *maps.Client, request *maps.GeocodingRequest) (lat float64, lng float64, address string, shortName string) {

	result, err := client.Geocode(context.Background(), request)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	lat = result[0].Geometry.Location.Lat
	lng = result[0].Geometry.Location.Lng
	address = result[0].FormattedAddress
	shortName = result[0].AddressComponents[0].ShortName

	return lat, lng, address, shortName
}

func toJson(spots []Spot) []byte {

	// Convert Strict To Json
	b, err := json.Marshal(&spots)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(100)
	}

	var out bytes.Buffer
	json.Indent(&out, b, "    ", "    ")

	return b
}

func writeFile(path string, b []byte) {

	// Create New File
	file, err := os.Create(path)

	if err != nil {
		os.Exit(101)
	}

	// File Close
	defer file.Close()

	// Write Bytes
	_, err = file.Write(b)
	if err != nil {
		os.Exit(102)
	}

}
