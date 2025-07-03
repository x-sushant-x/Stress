package main

import (
	"errors"
	"log"
	"os"
	"strconv"
)

func main() {
	args := os.Args

	// stress -u https://beyondthesyntax.substack.com -n 100000 -c 30

	config, err := parseArgs(args)
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	url                string
	totalRequests      int
	concurrentRequests int
}

func parseArgs(args []string) (*Config, error) {
	URL := args[2]
	tRequests := args[4]
	cRequests := args[6]

	totalRequests, err := strconv.Atoi(tRequests)
	if err != nil || totalRequests <= 0 {
		return nil, errors.New("invalid value of -n. Must be a positive integer.")
	}

	concurrentRequests, err := strconv.Atoi(cRequests)
	if err != nil || totalRequests <= 0 {
		return nil, errors.New("invalid value of -c. Must be a positive integer.")
	}

	return &Config{
		url:                URL,
		totalRequests:      totalRequests,
		concurrentRequests: concurrentRequests,
	}, nil
}
