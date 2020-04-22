package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type IpRange struct {
	Start uint64
	End   uint64
}

func main() {
	start := time.Now()

	file, err := os.Open(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	counter := 0

	ipRanges := []IpRange{}

	for scanner.Scan() {
		line := scanner.Text()
		res := strings.Split(line, "-")
		start, err := strconv.ParseUint(res[0], 10, 64)

		if err != nil {
			log.Fatal(err)
		}

		end, err := strconv.ParseUint(res[1], 10, 64)

		if err != nil {
			log.Fatal(err)
		}

		ipRanges = append(ipRanges, IpRange{start, end})

		counter++
	}

	ipRanges = mergeRanges(ipRanges)

	//fmt.Println(ipRanges)
	fmt.Println("Original filters: ", counter)
	fmt.Println("Resulting filters: ", len(ipRanges))

	elapsed := time.Since(start)
	fmt.Println("Time elapsed: ", elapsed)
}

func mergeRanges(ipRanges []IpRange) []IpRange {
	// sort by start range
	sort.Slice(ipRanges[:], func(i, j int) bool {
		return ipRanges[i].Start < ipRanges[j].Start
	})

	mergedIpRanges := []IpRange{}

	for i, ipRange := range ipRanges {
		if i == 0 {
			mergedIpRanges = append(mergedIpRanges, ipRange)
			continue
		}

		lastMergedIndex := len(mergedIpRanges) - 1
		lastMerged := mergedIpRanges[lastMergedIndex]

		if lastMerged.End < ipRange.Start {
			mergedIpRanges = append(mergedIpRanges, ipRange)
		} else {
			if mergedIpRanges[lastMergedIndex].End < ipRange.End {
				mergedIpRanges[lastMergedIndex].End = ipRange.End
			}
		}
	}

	return mergedIpRanges
}
