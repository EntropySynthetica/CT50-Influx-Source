package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

type thermo_stats struct {
	Temp     float64 `json:"temp"`
	Tmode    int     `json:"tmode"`
	Fmode    int     `json:"fmode"`
	Override int     `json:"override"`
	Hold     int     `json:"hold"`
	THeat    float64 `json:"t_heat"`
	TCool    float64 `json:"t_cool"`
	Tstate   int     `json:"tstate"`
	Fstate   int     `json:"fstate"`
	Time     struct {
		Day    int `json:"day"`
		Hour   int `json:"hour"`
		Minute int `json:"minute"`
	} `json:"time"`
	TTypePost int `json:"t_type_post"`
}

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("HOST")
	apiKey := os.Getenv("APIKEY")
	org := os.Getenv("ORG")
	bucket := os.Getenv("BUCKET")
	ct50IP := os.Getenv("CT50IP")

	// Poll the CT50
	poll_url := "http://" + ct50IP + "/tstat"

	response, err_conn := http.Get(poll_url)

	if err_conn != nil {
		fmt.Print(err_conn.Error())
		os.Exit(1)
	}

	response_data, err_http := ioutil.ReadAll(response.Body)
	if err_http != nil {
		log.Fatal(err_http)
	}

	// Parse the JSON from the CT50
	var response_stats thermo_stats
	json.Unmarshal([]byte(response_data), &response_stats)

	// Create a new client using an InfluxDB server base URL and an authentication token
	client := influxdb2.NewClient(host, apiKey)

	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking(org, bucket)

	// Create point using fluent style
	p := influxdb2.NewPointWithMeasurement("ct50").
		AddTag("field", "value").
		AddField("temp", response_stats.Temp).
		AddField("theat", response_stats.THeat).
		AddField("tcool", response_stats.TCool).
		AddField("override", response_stats.Override).
		AddField("fmode", response_stats.Fmode).
		AddField("tmode", response_stats.Tmode).
		AddField("hold", response_stats.Hold).
		AddField("tstate", response_stats.Tstate).
		AddField("fstate", response_stats.Fstate).
		SetTime(time.Now())

	err := writeAPI.WritePoint(context.Background(), p)
	if err != nil {
		fmt.Println(err)
	}

	// Ensures background processes finishes
	client.Close()
}
