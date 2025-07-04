package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/aquasecurity/table"
)

type Config struct {
	url                string
	totalRequests      int
	concurrentRequests int
}

type HTTPResponse struct {
	Time          time.Duration
	StatusCode    int
	UnableToReach bool
}

type Response struct {
	TotalHits int
	Success   int
	Failed    int
}

func main() {
	args := os.Args

	// ./stress -u https://beyondthesyntax.substack.com -n 100000 -c 30

	config, err := parseArgs(args)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		os.Exit(-1)
	}

	req := buildRequest(config)
	if req == nil {
		fmt.Println("Unable to build HTTP request")
		os.Exit(1)
	}

	var wg sync.WaitGroup

	semaphore := make(chan struct{}, config.concurrentRequests)
	resp := &Response{}

	for i := 0; i < config.totalRequests; i++ {
		wg.Add(1)

		semaphore <- struct{}{}

		go func() {
			defer wg.Done()
			defer func() {
				<-semaphore
			}()

			httpResp := makeHTTPRequest(req)

			if httpResp.UnableToReach {
				return
			}

			if httpResp.StatusCode == http.StatusOK {
				resp.Success++
			} else {
				resp.Failed++
			}

			resp.TotalHits++
		}()
	}

	wg.Wait()

	close(semaphore)

	t := table.New(os.Stdout)
	t.SetHeaders("Stat", "Count")
	t.AddRow("Total Hits", strconv.Itoa(resp.TotalHits))
	t.AddRow("Success", strconv.Itoa(resp.Success))
	t.AddRow("Failed", strconv.Itoa(resp.Failed))

	t.Render()
}

func parseArgs(args []string) (*Config, error) {
	URL := args[2]
	tRequests := args[4]
	cRequests := args[6]

	totalRequests, err := strconv.Atoi(tRequests)
	if err != nil || totalRequests <= 0 {
		return nil, errors.New("invalid value of -n. Must be a positive integer")
	}

	concurrentRequests, err := strconv.Atoi(cRequests)
	if err != nil || totalRequests <= 0 {
		return nil, errors.New("invalid value of -c. Must be a positive integer")
	}

	return &Config{
		url:                URL,
		totalRequests:      totalRequests,
		concurrentRequests: concurrentRequests,
	}, nil
}

func buildRequest(config *Config) *http.Request {
	req, err := http.NewRequest(http.MethodGet, config.url, nil)
	if err != nil {
		log.Println("unable to make HTTP request: " + err.Error())
		return nil
	}

	return req
}

func makeHTTPRequest(req *http.Request) *HTTPResponse {
	resp := &HTTPResponse{}

	now := time.Now()

	httpResp, err := http.DefaultClient.Do(req)

	spent := time.Since(now)
	resp.Time = spent

	if err != nil {
		resp.UnableToReach = true
		return resp
	}

	resp.StatusCode = httpResp.StatusCode
	return resp
}
