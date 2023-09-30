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
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/yusufpapurcu/wmi"
	"math"
	"os"
	"os/exec"
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
	//                          WWWWWW 100000000
	percent, err := cpu.Percent(200000000, false)
	if err != nil {
		fmt.Printf("Error getting CPU usage: %v\n", err)
		return
	}

	for i, pct := range percent {
		fmt.Printf("CPU%d: %.2f%%\n", i, pct)
	}
}

func roundToDecimalPlaces(num float64, decimalPlaces int) float64 {
	// Calculate the factor for rounding
	roundingFactor := math.Pow(10, float64(decimalPlaces))

	// Round the number to the specified decimal places
	rounded := math.Round(num*roundingFactor) / roundingFactor

	return rounded
}

func parseBootTime(bootTimeStr string) (time.Time, error) {

	parts := strings.Split(bootTimeStr, ".")

	if len(parts) != 2 {
		fmt.Println("Invalid input format")
		return time.Time{}, nil
	}

	t, err := time.Parse("20060102150405", parts[0])
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}, err
	}

	return t, nil
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

			fmt.Printf("===CPU===\n")
			fmt.Printf("CPU Name: %s\n", processor.Name)
			fmt.Printf("CPU Manufacturer: %s\n", processor.Manufacturer)
			fmt.Printf("Base Clock Speed: %.2f GHz\n", roundToDecimalPlaces(float64(processor.CurrentClockSpeed)/1000.0, 2))
			fmt.Printf("Number of Cores: %d\n", processor.NumberOfCores)
			fmt.Printf("Number of Logical Processors: %d\n", processor.NumberOfLogicalProcessors)
			fmt.Printf("L2 Cache Size: %d KB\n", processor.L2CacheSize)
			fmt.Printf("L3 Cache Size: %d KB\n", processor.L3CacheSize)
			bootTime, err := parseBootTime(operatingSystem.LastBootUpTime)
			if err == nil {
				uptime := time.Since(bootTime)
				parts := strings.Split(uptime.String(), ".")
				fmt.Printf("Uptime: %s\n", parts[0]+"s")
			} else {
				fmt.Printf("Error parsing boot time: %v\n", err)
			}
			getCpuUsage()

			// Memory
			fmt.Printf("\n===Memory===\n")
			memoryInfo, err := mem.VirtualMemory()
			if err != nil {
				fmt.Printf("Error getting memory info: %v\n", err)
				return
			}

			fmt.Printf("Total Memory: %.2f GB\n", float64(memoryInfo.Total)/float64(1024*1024*1024))
			fmt.Printf("Used Memory: %.2f GB\n", float64(memoryInfo.Used)/float64(1024*1024*1024))
			fmt.Printf("Free Memory: %.2f GB\n", float64(memoryInfo.Free)/float64(1024*1024*1024))
			fmt.Printf("Memory Usage: %.2f%%\n", memoryInfo.UsedPercent)

			// Disk
			fmt.Printf("\n===Disk===\n")
			partitions, err := disk.Partitions(false)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			for _, partition := range partitions {
				fmt.Printf("Device: %s\n", partition.Device)
				fmt.Printf("Mount point: %s\n", partition.Mountpoint)
				fmt.Printf("File system type: %s\n", partition.Fstype)

				// Get disk usage statistics
				usage, err := disk.Usage(partition.Mountpoint)
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}

				fmt.Printf("Total capacity: %.2f GB\n", float64(usage.Total)/(1024*1024*1024))
				fmt.Printf("Used space: %.2f GB\n", float64(usage.Used)/(1024*1024*1024))
				fmt.Printf("Free space: %.2f GB\n", float64(usage.Free)/(1024*1024*1024))
				fmt.Printf("Usage percentage: %.2f%%\n", usage.UsedPercent)
			}

		} else {
			fmt.Println("No CPU information found.")
		}
	}
}
