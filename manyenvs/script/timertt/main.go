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

	"github.com/hiroyaonoe/bcop-go/contrib/net/http/bcophttp"
	"github.com/hiroyaonoe/bcop-go/propagation"
	"github.com/hiroyaonoe/bcop-go/protocol/header"
)

var (
	errCount     = 0
	emu          sync.Mutex
	successCount = 0
	smu          sync.Mutex
)

func main() {
	url := flag.String("url", "", "URL")

	envID := flag.String("env-id", "", "BCoP env-id")
	bcop := flag.Bool("bcop", false, "use BCoP or not")
	reqPerSec := flag.Int("req-per-sec", 1000, "request per second")
	reqDuration := flag.Int("duration", 10, "duration second")
	clientNum := flag.Int("client-num", 16, "the number of client connections (equals the number of environments if no using proxy)")
	payload := flag.Int("payload", 1000, "payload byte")
	useProxy := flag.Bool("proxy", false, "use proxy or not")

	flag.Parse()

	reqTotal := *reqPerSec * *reqDuration

	client := http.DefaultClient
	ctx := context.Background()

	if *bcop {
		client = &http.Client{
			Transport: bcophttp.NewTransport(nil, propagation.EnvID{}),
		}

		h := header.NewV1()
		h.Set(propagation.EnvIDHeader, *envID)
		ctx = propagation.EnvID{}.Extract(ctx, propagation.NewBCoPCarrier(h))
	}

	client.Timeout = 1 * time.Second

	ctx, cancel := context.WithCancel(ctx)

	data := bytes.Repeat([]byte("0"), *payload)

	var wg sync.WaitGroup

	interval := time.Duration(*clientNum) * time.Second / time.Duration(*reqPerSec)
	fmt.Printf("interval: %d", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt)

	URLList := make([]string, *clientNum+1)
	for i := 1; i <= *clientNum; i++ {
		if *useProxy {
			URLList[i] = *url
		} else {
			URLList[i] = fmt.Sprintf("%s%03d", *url, i) // port番号を環境によって連続で変化させる
		}
	}

	i := 0
	// 終了はカウントで測るので不要
	// after := time.After(time.Duration(*reqDuration) * time.Second)
	for {
		select {
		case <-ticker.C:
			for j := 1; j <= *clientNum; j++ {
				if i >= reqTotal {
					goto WAIT
				}
				wg.Add(1)
				go func(count, envNum int) {
					defer wg.Done()
					fmt.Printf("% 8d:Start:% 6d;% 3d\n", successCount, count, envNum)

					req, err := http.NewRequestWithContext(ctx, http.MethodGet, URLList[envNum], bytes.NewReader(data))
					if err != nil {
						raiseErr(fmt.Sprintf("% 6d;% 3d: making request struct", count, envNum), err)
						return
					}
					if !*bcop {
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
				}(i, j)
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
