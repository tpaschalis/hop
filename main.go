package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type result int

const (
	passed result = iota
	failed
	noCompile
	noStart

	goTestFail   = "exit status 1"
	goTestNoComp = "exit status 2"

	usage = `Usage:
› Press a to run all tests.
› Press f to run failed tests.
› Press c to run all tests, forcing a re-run of cached results.
› Press p to run all tests satisfying a regex
› Press l to list all tests.
› Press d to list tests filtered by a regex.
› Press t to only run tests.
› Press b to only run benchmarks.
› Press q to quit.
› Press Enter to trigger a test run.
`

	summary = `Test Summary
› Ran %v tests %v %v.
› %v.
`
)

// var StdoutBuf, StderrBuf bytes.Buffer
var prev string

func main() {

	clearPlusUsage()

	for {
		fmt.Printf(prev)
		prev = ""
		scanner := bufio.NewScanner(os.Stdin)
		ok := scanner.Scan()
		if !ok {
			fmt.Println("\nExiting...")
			break
		}

		switch scanner.Text() {
		case "a":
			clear()
			res, stdout, stderr := goTest()
			handleOutput(res, stdout, stderr)

		case "c":
			clear()
			res, stdout, stderr := goTest("-count=1")
			handleOutput(res, stdout, stderr)

		case "p":
			clearPlusUsage()
			fmt.Printf("\npattern › ")

			scanner := bufio.NewScanner(os.Stdin)
			ok := scanner.Scan()
			if !ok {
				break
			}
			pattern := scanner.Text()

			clear()

			fmt.Printf("Running `go test` for pattern : %s\n", pattern)

			res, stdout, stderr := goTest("-run", pattern)
			handleOutput(res, stdout, stderr)

		case "l":
			clear()
			prev = ""
			res, stdout, stderr := goTest("-list", ".")
			handleOutput(res, stdout, stderr)
		case "d":
			clearPlusUsage()
			fmt.Printf("\npattern › ")

			scanner := bufio.NewScanner(os.Stdin)
			ok := scanner.Scan()
			if !ok {
				break
			}
			pattern := scanner.Text()

			clear()

			fmt.Printf("Running `go test -list` for pattern : %s\n", pattern)

			res, stdout, stderr := goTest("-list", pattern)
			handleOutput(res, stdout, stderr)
		case "t":
			clear()
			res, stdout, stderr := goTest("-run", "^Test")
			handleOutput(res, stdout, stderr)
		case "b":
			clear()
			res, stdout, stderr := goTest("-run", "^Benchmark")
			handleOutput(res, stdout, stderr)
		case "q":
			os.Exit(0)
		case "":
			clear()
		default:
			prev = "Selected option '" + scanner.Text() + "' is not supported.\n"
			clear()
		}

		// StdoutBuf.Reset()
		fmt.Printf(prev)
		fmt.Printf(usage)
		fmt.Printf(`› `)
		// fmt.Printf(StdoutBuf.String())

	}
}

func goTest(arguments ...string) (result, bytes.Buffer, bytes.Buffer) {
	var res result
	var stdoutBuf, stderrBuf bytes.Buffer

	args := []string{"test", "-v", "./..."}
	args = append(args, arguments...)
	cmd := exec.Command("go", args...)

	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	if runtime.GOOS == "windows" {
		cmd = exec.Command("tasklist")
	}

	err := cmd.Start()
	if err != nil {
		return noStart, stdoutBuf, stderrBuf
	}

	err = cmd.Wait()
	if err != nil {
		if err.Error() == goTestFail {
			res = failed
		}
		if err.Error() == goTestNoComp {
			res = noCompile
		}
	}

	cmd.Process.Kill()

	return res, stdoutBuf, stderrBuf
}

func handleOutput(res result, stdout, stderr bytes.Buffer) {
	//fmt.Println(stdout.String())
	//fmt.Println(stderr.String())

	// var events []event

	// scanner := bufio.NewScanner(&stdout)
	// for scanner.Scan() {
	// 	var ev event
	// 	line := scanner.Bytes()

	// 	err := json.Unmarshal(line, &ev)
	// 	if err != nil {
	// 		fmt.Println("could not unmarshal line : ", string(line), "with error :", err)
	// 		continue
	// 	}
	// 	//ev.Output = strings.Replace(ev.Output, "\n", `\n`, -1)
	// 	ev.Output = strings.Replace(ev.Output, "\n", "", -1)
	// 	events = append(events, ev)
	// }

	// for _, ev := range events {
	// 	if ev.Action == "skip" || strings.Contains(ev.Output, "\t(cached)") {
	// 		skip(ev.Output + "\n")
	// 	}
	// 	if strings.Contains(ev.Output, "[no test files]") {
	// 		ignore("Package " + ev.Package + " skipped -- no test files\n\n")
	// 	}

	// 	if ev.Action == "run" {
	// 		running("Running " + ev.Package + "/" + ev.Test + "\n")
	// 	}

	// 	if ev.Action == "output" && !strings.Contains(ev.Output, "=== RUN") && !strings.Contains(ev.Output, "--- PASS") && !strings.Contains(ev.Output, "--- FAIL") && ev.Output != "PASS" && ev.Output != "FAIL" && !strings.Contains(ev.Output, "\t(cached)") && !strings.Contains(ev.Output, "\t[no test files]") {
	// 		fmt.Println(ev.Output)
	// 	}
	// 	if ev.Action == "pass" {
	// 		if ev.Test == "" {
	// 			success("\nPackage " + ev.Package + " passed\n\n")
	// 		} else {
	// 			pass("   PASS " + ev.Package + "/" + ev.Test + "\n\n")
	// 		}
	// 	}
	// 	if ev.Action == "fail" {
	// 		if ev.Test == "" {
	// 			failure("\nPackage " + ev.Package + " failed\n\n")
	// 		} else {
	// 			fail("   FAIL " + ev.Package + "/" + ev.Test + "\n\n")
	// 		}
	// 	}
	// }

	//w := new(tabwriter.Writer)
	//w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	//for _, e := range events {
	//	fmt.Fprintln(w, e.Time, e.Action, e.Package, e.Test, e.Elapsed, e.Output)
	//
	//}
	//fmt.Fprintln(w)
	//w.Flush()

	switch res {
	case noCompile:
		red("Failed to run tests due to compiler errors\n")
		fmt.Println(stderr.String())
	case noStart:
		red("Failed to start the `go test` command\n")
		fmt.Println(stderr.String())
	case passed:
		success("\nPASS\n")
	case failed:
		failure("\nFAIL\n")
	}

}

type event struct {
	Time    *time.Time `json:",omitempty"`
	Action  string
	Package string   `json:",omitempty"`
	Test    string   `json:",omitempty"`
	Elapsed *float64 `json:",omitempty"`
	Output  string   `json:",omitempty"`
}
