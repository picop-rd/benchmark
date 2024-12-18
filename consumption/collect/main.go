package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	name := flag.String("name", "", "project name")
	timestamp := flag.String("timestamp", "", "timestamp to exec command (RFC3339)")
	interval := flag.Int("interval", 10, "interval to exec command (second)")
	duration := flag.Int("duration", 300, "duration to exec command (second)")
	flag.Parse()
	cmd := flag.Args()

	fmt.Printf("name: %s\n", *name)
	fmt.Printf("timestamp: %s\n", *timestamp)
	fmt.Printf("interval: %d\n", *interval)
	fmt.Printf("duration: %d\n", *duration)
	fmt.Printf("cmd: %s\n", strings.Join(cmd, " "))

	if len(*name) == 0 || *interval == 0 || *duration == 0 || len(cmd) < 2 || len(*timestamp) == 0 {
		flag.Usage()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit
		cancel()
	}()

	dirname := filepath.Join(*name, *timestamp)
	if _, err := os.Stat(dirname); err == nil {
		fmt.Printf("duplicated timestamp: %s\n", dirname)
		return
	}
	err := os.MkdirAll(dirname, 0755)
	if err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		return
	}
	log.Printf("created: %s\n", dirname)

	collect(ctx, dirname, cmd, *interval, *duration)
	log.Println("completed")
}

func collect(ctx context.Context, dirname string, cmd []string, interval, duration int) {
	var wg sync.WaitGroup
	after := time.After(time.Duration(duration) * time.Second)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for cnt := 0; ; cnt++{
		select {
		case <-ctx.Done():
			return
		case <-after:
			log.Println("waiting")
			wg.Wait()
			return
		case <-ticker.C:
			wg.Add(1)
			go func(cnt int) {
				log.Printf("exec command: %d\n", cnt)
				err := execCmd(dirname, cmd)
				if err != nil {
					fmt.Printf("exec command %d error: %v\n", cnt, err)
				}
				wg.Done()
				log.Printf("finished: %d\n", cnt)
			}(cnt)
		}
	}
}

func execCmd(dirname string, cmd []string) error {
	now := time.Now().Unix()
	c := exec.Command(cmd[0], cmd[1:]...)
	bytes, err := c.Output()
	c.Stderr = os.Stderr
	if err != nil {
		return err
	}
	file := filepath.Join(dirname, fmt.Sprintf("%d", now))
	err = os.WriteFile(file, bytes, 0644)
	if err != nil {
		return err
	}
	log.Printf("saved file: %s\n", file)
	return nil
}
