package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	analyzer "github.com/FujitsuLaboratories/ChaincodeAnalyzer/analyze"
)

// Reference: github.com/golang/lint

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tccanalyzer [files]\n")
	fmt.Fprintf(os.Stderr, "\tccanalyzer [directory]\n")
	flag.PrintDefaults()
}

// main func
func main() {
	flag.Usage = usage
	flag.Parse()

	// logging
	f, err := os.OpenFile("ccanalyzer.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	logger := log.New(f, "", log.Lshortfile|log.LstdFlags)

	var rundir, runfile int
	var args []string

	if flag.NArg() == 0 {
		usage()
		os.Exit(2)
	} else {
		for _, arg := range flag.Args() {
			if isDir(arg) {
				rundir = 1
				args = append(args, arg)
			} else if isFileExists(arg) && !isTest(arg) {
				runfile = 1
				args = append(args, arg)
			}
		}
	}

	if rundir+runfile != 1 {
		usage()
		os.Exit(2)
	}

	switch {
	case rundir == 1:
		fmt.Println("Target Dir: ", args)
		for _, dir := range args {
			analyzeDir(logger, dir)
		}
	case runfile == 1:
		fmt.Println("Target Files: ", args)
		analyzeFiles(logger, args...)
	}
}

func isDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return err == nil && fi.IsDir()
}

func isFileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func isTest(filename string) bool {
	return strings.HasSuffix(filename, "_test.go")
}

func isGoFile(filename string) bool {
	return strings.HasSuffix(filename, ".go")
}

func analyzeDir(logger *log.Logger, dirname string) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	var paths []string
	for _, file := range files {
		if isGoFile(file.Name()) && !isTest(file.Name()) {
			paths = append(paths, filepath.Join(dirname, file.Name()))
		}
	}
	if len(paths) > 0 {
		analyzeFiles(logger, paths...)
	} else {
		fmt.Printf("The directory [%s] does not include go file\n", dirname)
		os.Exit(1)
	}
}

func analyzeFiles(logger *log.Logger, filenames ...string) {
	files := make(map[string][]byte)
	for _, filename := range filenames {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		files[filename] = src
	}

	a := new(analyzer.Analyzer)
	ps, err := a.AnalyzeFiles(logger, files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	for _, p := range ps {
		if p.Validity {
			fmt.Println("## Category ", p.Category)
			fmt.Println("## Function ", p.Function)
			fmt.Println("## VarName ", p.VarName)
			fmt.Println("## Position ", p.Position)
			fmt.Println(p.LineText)
			fmt.Println("## Affected Position ", p.AffectedPosition)
			fmt.Println(p.AffectedLineText)
		}
	}
}
