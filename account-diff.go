package main

import (
	"sort"
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

func parseCsv(file string, dateCol int, debitCol int, creditCol int, descCol int, start time.Time, dateFmt string) []Line {
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

  if (len(lines) == 0) {
    fmt.Printf("No lines found in %v!  Do not save file as UTF-8!\n", file)
    os.Exit(1)
  }

  sort.Slice(lines, func(i, j int) bool {
		return lines[i].date.Before(lines[j].date)
	})

	return lines
}

func compareLines(a1 []Line, a2 []Line) ([]Line, []Line) {
	for i1 := range a1 {
		if a1[i1].found {
			continue
		}

		for i2 := range a2 {
			if a2[i2].found {
				continue
			}

			days := a1[i1].date.Sub(a2[i2].date).Hours()

			if a1[i1].amt == a2[i2].amt && days >= -7*24 && days <= 7*24 {
				a1[i1].found = true
				a2[i2].found = true

				break
			}
		}
	}

	return a1, a2
}

func main() {
	// first parameter is the bank export
	var a1 []Line

	switch os.Args[1] {
	case "verity":
		fmt.Println("********* Parsing Verity ************")
		a1 = parseCsv(os.Args[2], 1, 2, 2, 7, time.Time{}, "1/2/2006")
	case "amex":
		fmt.Println("********* Parsing AmEx ************")
		a1 = parseCsv(os.Args[2], 0, 7, 7, 2, time.Time{}, "1/2/06")
	case "capone":
		fmt.Println("********* Parsing CapOne ************")
		a1 = parseCsv(os.Args[2], 0, 5, 6, 3, time.Time{}, "2006-01-02")
	default:
		panic("first arg must be 'verity', 'capone' or 'amex'")
	}

	// second format is from quickbooks
	fmt.Println("********* Parsing Quickbooks ************")

	// be sure and verify credit, debit and desc cols are 4, 5, 3  -- sometimes theyre 5, 6, 4!
	a2 := parseCsv(os.Args[3], 0, 5, 6, 4, a1[0].date, "01/02/2006")
	// fmt.Printf("a1: %v\n a2: %v\n", a1, a2)

	compareLines(a1, a2)

	fmt.Println("Missing from quickbooks: ")
	for _, v := range a1 {
		if v.found {
			continue
		}
		fmt.Printf("%s\t%0.2f\t%s\t%v\n", v.date.Format("01/02/2006"), v.amt, v.desc, v.found)
	}

	fmt.Println("Extra in quickbooks: ")
	for _, v := range a2 {
		if v.found {
			continue
		}
		fmt.Printf("%s\t%v\t%s\n", v.date.Format("01/02/2006"), v.amt, v.desc)
	}
}
