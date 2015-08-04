package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alext/temperature-monitor/sensor"
)

func main() {
	sensorID := flag.String("sensorid", "", "The ID of the sensor device")
	logDest := flag.String("log", "STDERR", "Where to log to - STDOUT, STDERR or a filename")

	flag.Parse()

	setupLogging(*logDest)

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

func setupLogging(destination string) {
	switch destination {
	case "STDERR":
		log.SetOutput(os.Stderr)
	case "STDOUT":
		log.SetOutput(os.Stdout)
	default:
		file, err := os.OpenFile(destination, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			log.Fatalf("Error opening log %s: %s", destination, err.Error())
		}
		log.SetOutput(file)
	}
}
