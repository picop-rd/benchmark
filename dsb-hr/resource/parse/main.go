package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const podPrefixLength = len("-6fd889d5dd-fnkpv")

var (
	orgNames = []string{"prox", "ingress", "mongodb", "memcached"}
)

type PodUsage struct {
	Data map[int]int
}

func main() {
	name := flag.String("name", "", "project name")
	timestamp := flag.String("timestamp", "", "timestamp to exec command (RFC3339)")
	input := flag.String("input", "", "input dir path")
	output := flag.String("output", "", "output dir path")
	flag.Parse()

	if len(*input) == 0 || len(*output) == 0 || len(*name) == 0 || len(*timestamp) == 0 {
		flag.Usage()
		return
	}

	fmt.Printf("name: %s\n", *name)
	fmt.Printf("timestamp: %s\n", *timestamp)
	fmt.Printf("input dir path: %s\n", *input)
	fmt.Printf("output dir path: %s\n", *output)

	// Directory where the .txt files are stored
	inputDir := filepath.Join(*input, *name, *timestamp)
	outputDir := filepath.Join(*output, *name, *timestamp)

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
	// Export to Raw CSV
	rawDir := filepath.Join(outputDir, "raw")
	err = os.MkdirAll(rawDir, 0755)
	if err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		return
	}
	log.Printf("created: %s\n", rawDir)
	exportToRawCSV(cpuUsage, filepath.Join(rawDir, "cpu.csv"))
	log.Printf("exported: %s\n", filepath.Join(rawDir, "cpu.csv"))
	exportToRawCSV(memoryUsage, filepath.Join(rawDir, "memory.csv"))
	log.Printf("exported: %s\n", filepath.Join(rawDir, "memory.csv"))

	// Export to Org CSV
	orgDir := filepath.Join(outputDir, "org")
	err = os.MkdirAll(orgDir, 0755)
	if err != nil {
		fmt.Printf("mkdir error: %v\n", err)
		return
	}
	log.Printf("created: %s\n", orgDir)
	exportToOrgTxt(cpuUsage, filepath.Join(orgDir, "cpu.txt"))
	log.Printf("exported: %s\n", filepath.Join(orgDir, "cpu.txt"))
	exportToOrgTxt(memoryUsage, filepath.Join(orgDir, "memory.txt"))
	log.Printf("exported: %s\n", filepath.Join(orgDir, "memory.txt"))
	exportToOrgCsv(cpuUsage, filepath.Join(orgDir, "cpu.csv"))
	log.Printf("exported: %s\n", filepath.Join(orgDir, "cpu.csv"))
	exportToOrgCsv(memoryUsage, filepath.Join(orgDir, "memory.csv"))
	log.Printf("exported: %s\n", filepath.Join(orgDir, "memory.csv"))
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

func exportToRawCSV(data map[string]*PodUsage, filename string) {
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

func exportToOrgTxt(data map[string]*PodUsage, filename string) {
	res := 0.0
	for _, usage := range data {
		ints := make([]int, 0, len(usage.Data))
		for _, value := range usage.Data {
			ints = append(ints, value)
		}
		avg := calculateAverage(ints)
		res += avg
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating Txt file:", err)
		return
	}
	defer file.Close()

	file.Write([]byte(fmt.Sprintf("%f", res)))
}

func exportToOrgCsv(data map[string]*PodUsage, filename string) {
	orgData := make(map[string]float64, len(orgNames))
	for _, name := range orgNames {
		orgData[name] = 0.0
	}
	for name, usage := range data {
		ints := make([]int, 0, len(usage.Data))
		for _, value := range usage.Data {
			ints = append(ints, value)
		}
		yet := true
		for _, orgName := range orgNames {
			if strings.Contains(name, orgName) {
				orgData[orgName] += calculateAverage(ints)
				yet = false
				break
			}
		}
		if yet {
			orgData["app"] += calculateAverage(ints)
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating Csv file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for name, value := range orgData {
		writer.Write([]string{name, fmt.Sprintf("%f", value)})
	}
}

func calculateAverage(values []int) float64 {
	sort.Ints(values)
	trimmedValues := values[len(values)/20 : 19*len(values)/20] // 上位下位5%を除外
	var sum int
	for _, v := range trimmedValues {
		sum += v
	}
	return float64(sum) / float64(len(trimmedValues))
}
