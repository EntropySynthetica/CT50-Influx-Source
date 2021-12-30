package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("HOST")
	apiKey := os.Getenv("APIKEY")
	org := os.Getenv("ORG")
	bucket := os.Getenv("BUCKET")

	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(host, apiKey)

	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(org, bucket)

	// Create point using fluent style
	p := influxdb2.NewPointWithMeasurement("stat").
		AddTag("unit", "temperature").
		AddField("avg", 23.2).
		AddField("max", 45.0).
		SetTime(time.Now())

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		fmt.Println(err)
	}

	// Ensures background processes finishes
	client.Close()
}
