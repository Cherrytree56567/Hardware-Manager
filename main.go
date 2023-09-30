/*
******************************
*      Hardware Manager      *
******************************

	Hardware Manager © 2023 by Cherrytree56567 is licensed under Attribution-NonCommercial-ShareAlike 4.0 International
*/
package main

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/yusufpapurcu/wmi"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Win32Processor struct {
	Name                      string
	Manufacturer              string
	MaxClockSpeed             uint32
	CurrentClockSpeed         uint32
	NumberOfCores             uint32
	NumberOfLogicalProcessors uint32
	L2CacheSize               uint32
	L3CacheSize               uint32
}

type Win32OperatingSystem struct {
	LastBootUpTime string
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls") // For Windows
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getCpuUsage() {
	_, err := cpu.Info()
	if err != nil {
		fmt.Printf("Error getting CPU info: %v\n", err)
		return
	}

	percent, err := cpu.Percent(1000000000, false)
	if err != nil {
		fmt.Printf("Error getting CPU usage: %v\n", err)
		return
	}

	for i, pct := range percent {
		fmt.Printf("CPU%d: %.2f%%\n", i, pct)
	}
}

func parseBootTime(bootTimeStr string) (time.Time, error) {

	parts := strings.Split(bootTimeStr, ".")

	if len(parts) != 2 {
		fmt.Println("Invalid input format")
		return time.Time{}, nil
	}

	// Parse the numeric part as a float64.
	numericPart, err := strconv.ParseFloat(parts[0]+"."+parts[1], 64)
	if err != nil {
		fmt.Println("Error parsing numeric part:", err)
		return time.Time{}, nil
	}

	// Calculate hours, minutes, and seconds.
	hours := int(numericPart / 3600)
	numericPart -= float64(hours) * 3600
	minutes := int(numericPart / 60)
	numericPart -= float64(minutes) * 60
	seconds := int(numericPart)

	fmt.Printf("%d hours, %d minutes, %d seconds\n", hours, minutes, seconds)

	return time.Time{}, nil
}

func main() {
	clearScreen()

	for {
		var processors []Win32Processor
		if err := wmi.Query("SELECT * FROM Win32_Processor", &processors); err != nil {
			fmt.Printf("Error querying Win32_Processor: %v\n", err)
			return
		}

		var operatingSystems []Win32OperatingSystem
		if err := wmi.Query("SELECT LastBootUpTime FROM Win32_OperatingSystem", &operatingSystems); err != nil {
			fmt.Printf("Error querying Win32_OperatingSystem: %v\n", err)
			return
		}

		if len(processors) > 0 {
			operatingSystem := operatingSystems[0]
			processor := processors[0]

			clearScreen()

			fmt.Printf("CPU Name: %s\n", processor.Name)
			fmt.Printf("CPU Manufacturer: %s\n", processor.Manufacturer)
			fmt.Printf("Max Clock Speed: %d MHz\n", processor.MaxClockSpeed)
			fmt.Printf("Current Clock Speed: %d MHz\n", processor.CurrentClockSpeed)
			fmt.Printf("Number of Cores: %d\n", processor.NumberOfCores)
			fmt.Printf("Number of Logical Processors: %d\n", processor.NumberOfLogicalProcessors)
			fmt.Printf("L2 Cache Size: %d KB\n", processor.L2CacheSize)
			fmt.Printf("L3 Cache Size: %d KB\n", processor.L3CacheSize)
			bootTime, err := parseBootTime(operatingSystem.LastBootUpTime)
			if err == nil {
				//uptime := time.Since(bootTime)
				fmt.Printf("Uptime: %f:%f:%f\n", time.Since(bootTime).Hours(), time.Since(bootTime).Minutes(), time.Since(bootTime).Seconds())
			} else {
				fmt.Printf("Error parsing boot time: %v\n", err)
			}
			getCpuUsage()
		} else {
			fmt.Println("No CPU information found.")
		}
	}
}
