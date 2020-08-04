package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type event struct {
	Time    *time.Time `json:",omitempty"`
	Action  string
	Package string   `json:",omitempty"`
	Test    string   `json:",omitempty"`
	Elapsed *float64 `json:",omitempty"`
	Output  string   `json:",omitempty"`
}

func clearPlusUsage() {
	fmt.Printf("\033[2J")
	fmt.Printf("\033[1;1H")
	fmt.Printf(usage)
	fmt.Printf(`› `)
}
func clear() {
	fmt.Printf("\033[2J")
	fmt.Printf("\033[1;1H")
	fmt.Printf(`› `)
}

type colourizedWriter bytes.Buffer

type monochromeWriter bytes.Buffer

func (c *colourizedWriter) Write(p []byte) (n int, err error) {
	var events []event
	r := bytes.NewReader(p)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var ev event
		line := scanner.Bytes()

		err := json.Unmarshal(line, &ev)
		if err != nil {
			continue
		}
		//ev.Output = strings.Replace(ev.Output, "\n", `\n`, -1)
		ev.Output = strings.Replace(ev.Output, "\n", "", -1)
		events = append(events, ev)
	}

	for _, ev := range events {
		if ev.Action == "skip" || strings.Contains(ev.Output, "\t(cached)") {
			skip(ev.Output + "\n")
		}
		if strings.Contains(ev.Output, "[no test files]") {
			ignore("Package " + ev.Package + " skipped -- no test files")
		}

		if ev.Action == "run" {
			running("Running " + ev.Package + "/" + ev.Test + "\n")
		}

		if ev.Action == "output" && !strings.Contains(ev.Output, "=== RUN") && !strings.Contains(ev.Output, "--- PASS") && !strings.Contains(ev.Output, "--- FAIL") && ev.Output != "PASS" && ev.Output != "FAIL" && !strings.Contains(ev.Output, "\t(cached)") && !strings.Contains(ev.Output, "\t[no test files]") {
			fmt.Println(ev.Output)
		}
		if ev.Action == "pass" {
			if ev.Test == "" {
				success("\nPackage " + ev.Package + " passed\n\n")
			} else {
				pass("   PASS " + ev.Package + "/" + ev.Test + "\n\n")
			}
		}
		if ev.Action == "fail" {
			if ev.Test == "" {
				failure("\nPackage " + ev.Package + " failed\n\n")
			} else {
				fail("   FAIL " + ev.Package + "/" + ev.Test + "\n\n")
			}
		}
	}

	return len(p), nil
}

func (c *monochromeWriter) Write(p []byte) (n int, err error) {
	fmt.Printf(string(p))
	return len(p), nil
}
