package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/picop-rd/picop-go/contrib/net/http/picophttp"
	"github.com/picop-rd/picop-go/propagation"
	"github.com/picop-rd/picop-go/protocol/header"
)

var (
	errCount     = 0
	emu          sync.Mutex
	successCount = 0
	smu          sync.Mutex
)

func main() {
	url := flag.String("url", "", "URL")

	envID := flag.String("env-id", "", "PiCoP env-id")
	picop := flag.Bool("picop", false, "use PiCoP or not")
	reqPerSec := flag.Int("req-per-sec", 1000, "request per second")
	reqDuration := flag.Int("duration", 10, "duration second")
	payload := flag.Int("payload", 1000, "payload byte")
	envNum := flag.Int("env-num", 1, "env number")
	envTotal := flag.Int("env-total", 10, "env total")

	flag.Parse()

	reqTotal := *reqPerSec * *reqDuration / *envTotal

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

	client.Timeout = 1 * time.Second

	ctx, cancel := context.WithCancel(ctx)

	data := bytes.Repeat([]byte("0"), *payload)

	var wg sync.WaitGroup

	interval := time.Duration(*envTotal) * time.Second / time.Duration(*reqPerSec)
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
			if i >= reqTotal {
				goto WAIT
			}
			wg.Add(1)
			go func(count, envNum int) {
				defer wg.Done()
				fmt.Printf("% 8d:Start:% 6d;% 3d\n", successCount, count, envNum)

				req, err := http.NewRequestWithContext(ctx, http.MethodGet, *url, bytes.NewReader(data))
				if err != nil {
					raiseErr(fmt.Sprintf("% 6d;% 3d: making request struct", count, envNum), err)
					return
				}
				if !*picop {
					h := http.Header{}
					h.Add(propagation.EnvIDHeader, *envID)
					req.Header = h
				}

				resp, err := client.Do(req)
				if err != nil {
					raiseErr(fmt.Sprintf("% 6d;% 3d: send request", count, envNum), err)
					return
				}
				resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					raiseErr(fmt.Sprintf("% 6d;% 3d: status code not 200, %d", count, envNum, resp.StatusCode), err)
					return
				}

				smu.Lock()
				successCount += 1
				smu.Unlock()

				fmt.Printf("% 8d:End  :% 6d;% 3d\n", successCount, count, envNum)
			}(i, *envNum)
			i++
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
	fmt.Printf("Complete: err rate: %d / %d / %d\n", errCount, i, reqTotal)
}

// func divInt(total, div int) []int {
// 	base := total / div
// 	rest := total % div
// 	ret := make([]int, div)
// 	for i := 0; i < base; i++ {
// 		if i < rest {
// 			ret[i] = base + 1
// 		} else {
// 			ret[i] = base
// 		}
// 	}
// 	return ret
// }

func raiseErr(message string, err error) {
	emu.Lock()
	errCount++
	emu.Unlock()
	smu.Lock()
	successCount = 0
	smu.Unlock()
	fmt.Printf("% 8d:Error:%s: %s\n", successCount, message, err.Error())
}
