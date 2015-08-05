package main

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	Port    int                     `json:"port"`
	Sensors map[string]sensorConfig `json:"sensors"`
}

type sensorConfig struct {
	ID string `json:"id"`
}

func loadConfig(filename string) (*config, error) {
	c := &config{
		Port: defaultPort,
	}

	file, err := fs.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Config file '%s' not found, ignoring", filename)
			return c, nil
		}
		return nil, err
	}

	err = json.NewDecoder(file).Decode(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
