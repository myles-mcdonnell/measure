package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var errorCount int64
var lastError error
var requestsComplete int64

type requestStat struct {
	responseTimeMs int64
	startTime      time.Time
}

func remove(slice []requestStat, s int) []requestStat {
	return append(slice[:s], slice[s+1:]...)
}

func main() {

	url := flag.String("url", "https://tour.golang.org", "target url")
	concurrency := flag.Int("concurrency", 5, "concurrency request count")
	averageWindow := flag.Int("averageWindow", 5, "window average seconds")

	flag.Parse()

	averageWindowDuration := time.Duration(*averageWindow)

	requestWindow := make([]requestStat, 5, 5)
	requestsCompleteChannel := make(chan requestStat, 4*(*concurrency))

	for i := 0; i < *concurrency; i++ {
		go func(channel chan requestStat) {
			for {
				makeRequest(channel, *url)
			}

		}(requestsCompleteChannel)
	}

	go func(channel chan requestStat) {
		for stat := range channel {
			requestWindow = append(requestWindow, stat)
		}
	}(requestsCompleteChannel)

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				var windowRequestCount int64
				var windowRequestResponseTimeTotal int64

				for _, stat := range requestWindow {
					if (time.Since(stat.startTime) > (averageWindowDuration * time.Second)) && (len(requestWindow) > 1) {
						requestWindow = remove(requestWindow, 1)
					} else {
						windowRequestCount++
						windowRequestResponseTimeTotal += stat.responseTimeMs
					}
				}

				var windowsAverageMs int64
				if windowRequestCount == 0 {
					windowsAverageMs = 0
				} else {
					windowsAverageMs = windowRequestResponseTimeTotal / windowRequestCount
				}

				var output = fmt.Sprintf("\rTotal Reqs: %d - Window Reqs: %v - Window Average (ms): %v",
					requestsComplete,
					windowRequestCount,
					windowsAverageMs)

				if lastError != nil {
					output = fmt.Sprintf("%s - Error Count: %d - Last Error: %s",
						output,
						errorCount,
						lastError)
				}

				fmt.Printf(output)
			}
		}
	}()

	bufio.NewReader(os.Stdin).ReadString('\n')
}

func makeRequest(channel chan requestStat, url string) {
	timeStart := time.Now()
	resp, err := http.Get(url)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	atomic.AddInt64(&requestsComplete, 1)
	if err != nil {
		atomic.AddInt64(&errorCount, 1)
		lastError = err
	} else {

		channel <- requestStat{responseTimeMs: int64(time.Since(timeStart).Nanoseconds()) / 1000000, startTime: timeStart}
	}
}
