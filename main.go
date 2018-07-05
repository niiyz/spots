package main

import (
	"encoding/csv"
	"flag"
	"github.com/joho/godotenv"
	"github.com/kr/pretty"
	"golang.org/x/net/context"
	"googlemaps.github.io/maps"
	"io"
	"log"
	"os"
	"bytes"
	"encoding/json"
	"fmt"
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

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {

	EnvLoad()

	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GEO_CODING_API_KEY")))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	flag.Parse()

	inputFile, err := os.Open(flag.Arg(0))
	failOnError(err)
	defer inputFile.Close()
	reader := csv.NewReader(inputFile)

	var spots []Spot

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}

		var address = record[0]

		r := &maps.GeocodingRequest{
			Address:  address,
			Language: "ja",
		}
		result, err := c.Geocode(context.Background(), r)
		if err != nil {
			log.Fatalf("fatal error: %s", err)
		}

		var lat = result[0].Geometry.Location.Lat
		var lng = result[0].Geometry.Location.Lng

		spot := Spot{Name: address, Address: result[0].FormattedAddress, Lat: lat, Lng: lng}
		spots = append(spots, spot)

		pretty.Println(result[0].AddressComponents[0].ShortName, result[0].FormattedAddress, lat, lng)

	}

	writeFile("spot.json", toJson(spots))
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