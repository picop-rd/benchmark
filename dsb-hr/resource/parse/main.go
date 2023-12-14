package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const podPrefixLength = len("-6fd889d5dd-fnkpv")

type PodUsage struct {
	Data map[int]int
}

func main() {
	input := flag.String("input", "", "input file path")
	output := flag.String("output", "", "output file path")
	flag.Parse()

	if len(*input) == 0 || len(*output) == 0 {
		flag.Usage()
		return
	}

	fmt.Printf("input file path: %s\n", *input)
	fmt.Printf("output file path: %s\n", *output)

	// Directory where the .txt files are stored
	inputDir := *input
	outputDir := *output

	// Maps to store pod usages
	cpuUsage := make(map[string]*PodUsage)
	memoryUsage := make(map[string]*PodUsage)

	// Process each file in the directory
	filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			unixTime, err := strconv.Atoi(info.Name())
			if err != nil {
				return err
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			log.Printf("processing: %s\n", path)
			processFile(file, unixTime, cpuUsage, memoryUsage)
		}
		return nil
	})

	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		return
	}
	log.Printf("created: %s\n", outputDir)
	// Export to CSV
	exportToCSV(cpuUsage, filepath.Join(outputDir, "cpu.csv"))
	log.Printf("exported: %s\n", filepath.Join(outputDir, "cpu.csv"))
	exportToCSV(memoryUsage, filepath.Join(outputDir, "memory.csv"))
	log.Printf("exported: %s\n", filepath.Join(outputDir, "memory.csv"))
}

func processFile(file *os.File, unixTime int, cpu, memory map[string]*PodUsage) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "CPU(cores)") {
			continue // Skip header line
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Printf("Invalid line: %s\n", line)
			continue // Invalid line
		}
		podName, err := converPodName(parts[0])
		if err != nil {
			fmt.Printf("Invalid pod name: %s error: %s\n", parts[0], err)
			continue // Invalid pod name
		}
		cpuValue, err := convertCPUValue(parts[1])
		if err != nil {
			fmt.Printf("Invalid CPU value: %s error: %s\n", parts[1], err)
			continue // Invalid CPU value
		}
		memoryValue, err := convertMemoryValue(parts[2])
		if err != nil {
			fmt.Printf("Invalid memory value: %s error: %s\n", parts[2], err)
			continue // Invalid memory value
		}

		if _, exists := cpu[podName]; !exists {
			cpu[podName] = &PodUsage{Data: make(map[int]int)}
		}
		cpu[podName].Data[unixTime] = cpuValue

		if _, exists := memory[podName]; !exists {
			memory[podName] = &PodUsage{Data: make(map[int]int)}
		}
		memory[podName].Data[unixTime] = memoryValue
	}
}

func converPodName(raw string) (string, error) {
	if len(raw) < podPrefixLength {
		return "", fmt.Errorf("invalid pod name length: %s", raw)
	}
	return raw[:len(raw)-podPrefixLength], nil
}

func convertCPUValue(raw string) (int, error) {
	switch raw[len(raw)-1:] {
	case "m":
		return strconv.Atoi(raw[:len(raw)-1])
	default:
		return 0, fmt.Errorf("invalid CPU value: %s", raw)
	}
}

func convertMemoryValue(raw string) (int, error) {
	switch raw[len(raw)-2:] {
	case "Mi":
		return strconv.Atoi(raw[:len(raw)-2])
	default:
		return 0, fmt.Errorf("invalid memory value: %s", raw)
	}
}

func exportToCSV(data map[string]*PodUsage, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	length := 0
	isFirst := true
	for podName, usage := range data {
		if isFirst {
			isFirst = false
			length = len(usage.Data) + 1
			// Write headers
			headers := make([]string, 0, length)
			headers = append(headers, "name")
			// Add Unix timestamps as headers
			for unixTime := range usage.Data {
				headers = append(headers, strconv.Itoa(unixTime))
			}
			writer.Write(headers)
		}
		row := make([]string, 0, length)
		row = append(row, podName)
		for _, value := range usage.Data {
			row = append(row, strconv.Itoa(value))
		}
		writer.Write(row)
	}
}
