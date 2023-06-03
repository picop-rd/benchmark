package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/picop-rd/picop-go/contrib/net/http/picophttp"
	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
)

func main() {
	url := flag.String("url", "", "URL")

	prefix := flag.String("prefix", "", "output csv file prefix")
	envID := flag.String("env-id", "", "PiCoP env-id")
	picop := flag.Bool("picop", false, "use PiCoP")
	reqPerSec := flag.Int("req-per-sec", 1000, "request per second")
	reqDuration := flag.Int("duration", 10, "duration second")
	clientNum := flag.Int("client-num", 16, "the number of client connections")
	payload := flag.Int("payload", 1000, "payload byte")

	flag.Parse()

	now := time.Now().Local().Format(time.RFC3339)
	reqTotal := *reqPerSec * *reqDuration

	out := fmt.Sprintf("%s-rtt-%s-%t-%d-%d-%d-%d-%s.csv", *prefix, *envID, *picop, *reqPerSec, *reqDuration, *clientNum, *payload, now)
	fmt.Printf("Output file: %s\n", out)

	if _, err := os.Stat(out); err == nil {
		fmt.Println("output file already exists")
		return
	}

	client := http.DefaultClient
	ctx := context.Background()

	if *picop {
		client = &http.Client{
			Transport: picophttp.NewTransport(nil, propagation.EnvID{}),
		}

		h := header.NewV1()
		h.Set(propagation.EnvIDHeader, *envID)
		ctx = propagation.EnvID{}.Extract(ctx, propagation.NewPiCoPCarrier(h))
	}

	ctx, cancel := context.WithCancel(ctx)

	data := bytes.Repeat([]byte("0"), *payload)

	type latencyQueue struct {
		startTime time.Time
		endTime   time.Time
	}

	latencyList := make([]latencyQueue, reqTotal)
	var mu sync.Mutex

	var wg sync.WaitGroup

	interval := time.Duration(*clientNum) * time.Second / time.Duration(*reqPerSec)
	fmt.Printf("interval: %d", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt)

	i := 0
	// after := time.After(time.Duration(*reqDuration) * time.Second)
	for {
		select {
		case <-ticker.C:
			for j := 0; j < *clientNum; j++ {
				if i >= reqTotal {
					goto WAIT
				}
				wg.Add(1)
				go func(count int) {
					defer wg.Done()
					fmt.Printf("Start: %d\n", count)

					req, err := http.NewRequestWithContext(ctx, http.MethodGet, *url, bytes.NewReader(data))
					if err != nil {
						fmt.Printf("Error %d:making request struct error: %s\n", count, err.Error())
						return
					}
					if !*picop {
						h := http.Header{}
						h.Add(propagation.EnvIDHeader, *envID)
						req.Header = h
					}

					queue := latencyQueue{}

					queue.startTime = time.Now()
					resp, err := client.Do(req)
					queue.endTime = time.Now()
					if err != nil {
						fmt.Printf("Error %d: request error: %s\n", count, err.Error())
						return
					}
					resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						fmt.Printf("Error %d: status code not 200, %d\n", count, resp.StatusCode)
						return
					}

					mu.Lock()
					latencyList[count] = queue
					mu.Unlock()
					fmt.Printf("End: %d\n", count)
				}(i)
				i++
			}
		case <-stopper:
			goto END
			// case <-after:
			// 	goto END
		}
	}
WAIT:
	wg.Wait()
END:
	cancel()
	fmt.Printf("Output: %s\n", out)
	file, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("count, latency, start, end\n")
	for i, v := range latencyList {
		start := v.startTime.UnixNano()
		end := v.endTime.UnixNano()
		// Use monotonic clocks for latency
		latency := v.endTime.Sub(v.startTime)
		fmt.Fprintf(file, "%d, %d, %d, %d\n", i, latency, start, end)
	}
	fmt.Println("Complete")
}
