package main

import (
	"flag"
	"fmt"
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
	configFile := flag.String("config-file", defaultConfigFile, "Path to the config file")
	returnVersion := flag.Bool("version", false, "return version information and exit")

	flag.Parse()

	if *returnVersion {
		fmt.Printf("temperature-monitor %s\n", versionInfo())
		os.Exit(0)
	}

	err := setupLogging(*logDest)
	if err != nil {
		log.Fatal(err)
	}

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

func setupLogging(destination string) error {
	switch destination {
	case "STDERR":
		log.SetOutput(os.Stderr)
	case "STDOUT":
		log.SetOutput(os.Stdout)
	default:
		file, err := fs.OpenFile(destination, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return fmt.Errorf("Error opening log %s: %s", destination, err.Error())
		}
		log.SetOutput(file)
	}
	return nil
}

func addSensors(config *config, srv *webserver.Webserver) error {
	for name, sensorConfig := range config.Sensors {
		s, err := sensor.NewW1Sensor(sensorConfig.ID)
		if err != nil {
			return err
		}
		srv.AddSensor(name, s)
	}
	return nil
}
