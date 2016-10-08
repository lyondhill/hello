package main

import "github.com/shirou/gopsutil/disk"
import "fmt"

func main() {
	parts, _ := disk.DiskPartitions(false)
	fmt.Printf("%v\n\n", parts)

	for _, part := range parts {
		usage, _ := disk.DiskUsage(part.Mountpoint)
		fmt.Printf("usage for %s: %v\n\n", part.Device, usage)
	}
}