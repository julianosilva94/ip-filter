package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type IpRange struct {
	Start uint64
	End uint64
}

func main() {
	file, err := os.Open(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	ipRanges := []IpRange{}

	scanner := bufio.NewScanner(file)
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

		found := false
		for index , ipRange := range ipRanges {
			if isIpRangesInConflict(IpRange{start, end}, ipRange) {
				ipRanges[index] = mergeIpRanges(IpRange{start, end}, ipRange)
				found = true
				break
			}
		}

		if !found {
			ipRanges = append(ipRanges, IpRange{start, end})
		}
	}

	sort.Slice(ipRanges[:], func(i, j int) bool {
		return ipRanges[i].Start < ipRanges[j].Start
	})

	for {
		size := len(ipRanges)
		lastI := 0

		if size > 1 {
			found := false

			for i := lastI; i < size; i++ {
				prev := ipRanges[i];

				for j := i + 1; j < size; j++ {
					actual := ipRanges[j]

					if isIpRangesInConflict(prev, actual) {
						found = true
						ipRanges[i] = mergeIpRanges(prev, actual)
						lastI = i+1
					}

					if found {
						ipRanges = append(ipRanges[:j], ipRanges[j+1:]...)
						break
					}
				}

				if found {
					break
				}
			}

			if !found {
				break
			}

		} else {
			break
		}
	}

	fmt.Println(ipRanges)
}

func isIpRangesInConflict(a IpRange, b IpRange) bool {
	// a overlap b in left side
	if a.Start < b.Start && a.End > b.Start && a.End < b.End {
		return true
	}

	if b.Start < a.Start && b.End > a.Start && b.End < a.End {
		return true
	}

	// a overlap b in right side
	if a.Start > b.Start && a.Start < b.End && a.End > b.End {
		return true
	}

	if b.Start > a.Start && b.Start < a.End && b.End > a.End {
		return true
	}

	// a is inside b
	if a.Start > b.Start && a.End < b.End {
		return true
	}

	if b.Start > a.Start && b.End < a.End {
		return true
	}

	return false
}

func mergeIpRanges(a IpRange, b IpRange) IpRange {
	ipRange := IpRange{}

	if a.Start <= b.Start {
		ipRange.Start = a.Start
	} else {
		ipRange.Start = b.Start
	}

	if a.End >= b.End {
		ipRange.End = a.End
	} else {
		ipRange.End = b.End
	}

	return ipRange
}
