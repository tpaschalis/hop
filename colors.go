package main

import "github.com/fatih/color"

func red(format string, a ...interface{}) {
	red := color.New(color.FgHiRed).Add(color.Bold)
	red.Printf(format, a...)
	color.Unset()
}

func success(format string, a ...interface{}) {
	boldBlack := color.New(color.FgBlack).Add(color.Bold)
	greenBackgroundBoldBlack := boldBlack.Add(color.BgHiGreen)
	greenBackgroundBoldBlack.Printf(format, a...)
	color.Unset()
}

func failure(format string, a ...interface{}) {
	boldBlack := color.New(color.FgBlack).Add(color.Bold)
	redBackgroundBoldBlack := boldBlack.Add(color.BgHiRed)
	redBackgroundBoldBlack.Printf(format, a...)
	color.Unset()
}

func ignore(format string, a ...interface{}) {
	boldWhite := color.New(color.FgHiWhite).Add(color.Bold)
	grayBackgroundBoldWhite := boldWhite.Add(color.BgHiBlack)
	grayBackgroundBoldWhite.Printf(format, a...)
	color.Unset()
}

func skip(format string, a ...interface{}) {
	italicsGrey := color.New(color.FgHiBlack).Add(color.Faint).Add(color.Italic)
	italicsGrey.Printf(format, a...)
	color.Unset()
}

func pass(format string, a ...interface{}) {
	faintGreen := color.New(color.FgHiGreen).Add(color.Faint)
	faintGreen.Printf(format, a...)
	color.Unset()
}

func fail(format string, a ...interface{}) {
	faintRed := color.New(color.FgHiRed).Add(color.Faint)
	faintRed.Printf(format, a...)
	color.Unset()
}
func running(format string, a ...interface{}) {
	underline := color.New().Add(color.Underline)
	underline.Printf(format, a...)
	color.Unset()
}
