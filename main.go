package main

import (
	"flag"
	"log"
	"os"

	"github.com/spf13/afero"

	"github.com/alext/temperature-monitor/sensor"
	"github.com/alext/temperature-monitor/webserver"
)

var fs afero.Fs = &afero.OsFs{}

const (
	defaultConfigFile = "./config.json"
	defaultPort       = 8081
)

func main() {
	logDest := flag.String("log", "STDERR", "Where to log to - STDOUT, STDERR or a filename")
	configFile := flag.String("config-file", "./config.json", "Path to the config file")

	flag.Parse()

	setupLogging(*logDest)

	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatal("Error opening config file : ", err)
	}

	srv := webserver.New(config.Port)
	err = addSensors(config, srv)
	if err != nil {
		log.Fatal(err)
	}

	err = srv.Run()
	if err != nil {
		log.Fatal(err)
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

func addSensors(config *config, srv *webserver.Webserver) error {
	for name, sensorConfig := range config.Sensors {
		s, err := sensor.New(sensorConfig.ID)
		if err != nil {
			return err
		}
		srv.AddSensor(name, s)
	}
	return nil
}
