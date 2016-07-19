package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Line struct {
	date  time.Time
	amt   float64
	desc  string
	found bool
}

const dateFmt = "01/02/2006"

func parseCsv(file string, dateCol int, debitCol int, creditCol int, descCol int, start time.Time) []Line {
	f, _ := os.Open(file)
	reader := csv.NewReader(bufio.NewReader(f))
	lines := make([]Line, 0, 4096)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if len(record) < dateCol || len(record) < creditCol || len(record) < descCol {
			fmt.Printf("Record too short, Skipping.  %v\n", record)
			continue
		}

		date, err := time.Parse(dateFmt, record[dateCol])
		if err != nil {
			fmt.Printf("Date parse error, Skipping %s (%v)\n", record[descCol], err)
			continue
		}

		if !start.IsZero() && date.Before(start) {
			continue
		}

		amtStr := record[debitCol]
		if amtStr == "" {
			amtStr = record[creditCol]
		}

		amtStr = strings.Replace(amtStr, ",", "", -1)

		amt, err := strconv.ParseFloat(amtStr, 64)
		if err != nil {
			fmt.Printf("Amount parse error!? %v (%v) (%v) \n", record, err, amtStr)
			continue
		}

		lines = append(lines, Line{date: date, amt: math.Abs(amt), desc: record[descCol]})
	}

	return lines
}

func compareLines(a1 []Line, a2 []Line) ([]Line, []Line) {
	for i1 := range a1 {
		if a1[i1].found {
			continue
		}

		for i2 := range a2 {
			_ = "breakpoint"
			if a1[i1].amt == a2[i2].amt && a1[i1].date.Equal(a2[i2].date) && !a2[i2].found {
				a1[i1].found = true
				a2[i2].found = true

				break
			}
		}
	}

	return a1, a2
}

func main() {
	// first parameter is the verity export
	a1 := parseCsv(os.Args[1], 1, 2, 2, 7, time.Time{})

	// second format is from quickbooks
	a2 := parseCsv(os.Args[2], 0, 4, 5, 3, a1[0].date)

	compareLines(a1, a2)

	fmt.Println("In first but not in second: ")
	for _, v := range a1 {
		if v.found {
			continue
		}
		fmt.Printf("%s\t%0.2f\t%s\t%b\n", v.date.Format(dateFmt), v.amt, v.desc, v.found)
	}

	fmt.Println("In second but not in first: ")
	for _, v := range a2 {
		if v.found {
			continue
		}
		fmt.Printf("%s\t%v\t%s\n", v.date.Format(dateFmt), v.amt, v.desc)
	}
}
