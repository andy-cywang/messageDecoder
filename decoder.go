package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var scanner *bufio.Scanner

// Build a map store the keys and values based on the sequence
func getMapping(header string) map[string]string {
	mapping := make(map[string]string)
	//digits: number of digits for the current key
	//index: decimal representation of the key
	digits, index := 1, 0

	for _, v := range header {

		// the count of keys in each group who has same length of keys is (2 to the power of n - 1)
		// when that count has been hit, reset index to 0 and increase digits by 1 to move to next group
		if index == 2<<(digits-1) - 1 {
			index = 0
			digits++
		}
		// get the binary representation of the current index
		key := strconv.FormatInt(int64(index), 2)

		// add 0 at the front to fill a vacancy of digits
		prefixLen := digits - len(key)
		for j := 0; j < prefixLen; j++ {
			key = "0" + key
		}

		// put key and value in the map
		mapping[key] = string(v)
		index++
	}

	return mapping
}

func decode(out io.Writer, header string) {
	// get the mapping
	mapping := getMapping(header)

	// line: current scanned line
	// segment: binary segment length representation
	// key: key going to be used to map character in the mapping
	var line, segment, key string
	// segLen: the decimal length for the current segment
	var segLen int

	scanner.Scan()

	// if line equals "000" means the current message is terminated here
	for line = scanner.Text(); line != "000"; {

		// split out the current segment info from current line
		segment, line = line[:3], line[3:]
		// get segment length from the binary representation
		fmt.Sscanf(segment, "%b", &segLen)

		for {
			// split out the current key from current line
			key, line = line[:segLen], line[segLen:]

			// if the current line is shorter than the segLen, concat next line to current line.
			if len(line) < segLen {
				scanner.Scan()
				line += scanner.Text()
			}

			// if segment terminator encountered, jump to next segment
			if strings.Count(key, "1") == segLen {
				break
			}

			// output mapped character to the file
			fmt.Fprintf(out, "%s", mapping[key])
		}
	}
	fmt.Fprintln(out)
}

func main() {
	in, _ := os.Open("input.txt")
	out, _ := os.Create("output.txt")
	defer in.Close()
	defer out.Close()

	scanner = bufio.NewScanner(in)

	var line string
	for scanner.Scan() {
		if line = scanner.Text(); line == "" {
			break
		}
		decode(out, line)
	}
}
