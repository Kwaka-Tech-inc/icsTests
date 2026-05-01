package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

func main() {
	var fileNames = []string{"link1.ics", "link2.ics", "link3.ics"}

	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		loc = time.Now().Location()
	}
	currentDate := time.Now().In(loc).Format("20060102")

	// Regex to match the date part of DTSTAMP, DTSTART, and DTEND
	re := regexp.MustCompile(`(DTSTAMP|DTSTART|DTEND)([^:]*:)(\d{8})(T\d{6}.*)`)

	for _, fName := range fileNames {
		err := updateICSFile(fName, currentDate, re)
		if err != nil {
			fmt.Printf("Error updating %s: %v\n", fName, err)
		} else {
			fmt.Printf("Successfully updated %s\n", fName)
		}
	}
}

func updateICSFile(filename, currentDate string, re *regexp.Regexp) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			line = re.ReplaceAllString(line, "${1}${2}"+currentDate+"${4}")
		}
		lines = append(lines, line)
	}
	file.Close() // Close the file before reopening to write

	if err := scanner.Err(); err != nil {
		return err
	}

	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, line := range lines {
		if _, err := outFile.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}
