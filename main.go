package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/alext/temperature-monitor/sensor"
)

func main() {
	sensorID := flag.String("sensorid", "", "The ID of the sensor device")

	flag.Parse()

	if *sensorID == "" {
		log.Fatal("no sensorid provided")
	}

	sensor, err := sensor.New(*sensorID)
	if err != nil {
		log.Fatal("Error opening sensor : ", err)
	}

	for {
		temp := sensor.Temperature()
		fmt.Printf("Temperature: %.3f\n", temp)

		time.Sleep(10 * time.Second)
	}
}
