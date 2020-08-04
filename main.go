package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
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

var (
	customPrinter io.Writer
	colourPrinter colourizedWriter
	bwPrinter     monochromeWriter
	baseArgs      []string
	colorized     bool
	prev          string
)

func main() {

	flag.BoolVar(&colorized, "color", false, "enable colorful output")
	flag.Parse()

	if colorized {
		customPrinter = &colourPrinter
		baseArgs = []string{"test", "-v", "./...", "-json"}
	} else {
		customPrinter = &bwPrinter
		baseArgs = []string{"test", "-v", "./..."}
	}

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
			res, stdout, stderr := goTest("-bench=.", "-run", "^Benchmark")
			handleOutput(res, stdout, stderr)
		case "q":
			os.Exit(0)
		case "":
			clear()
		default:
			prev = "Selected option '" + scanner.Text() + "' is not supported.\n"
			clear()
		}

		fmt.Println("")
		fmt.Printf(prev)
		fmt.Printf(usage)
		fmt.Printf(`› `)

	}
}

func goTest(arguments ...string) (result, bytes.Buffer, bytes.Buffer) {
	var res result
	var stdoutBuf, stderrBuf bytes.Buffer

	args := baseArgs
	args = append(args, arguments...)

	cmd := exec.Command("go", args...)

	cmd.Stdout = io.MultiWriter(customPrinter, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(customPrinter, &stdoutBuf)

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
