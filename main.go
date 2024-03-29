package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	flag "github.com/spf13/pflag"
	"github.com/surface-security/scanner-go-entrypoint/scanner"
)

const DEFAULT_USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.142 Safari/537.36 github.com/surface-security/scanner-httpx"

func main() {
	s := scanner.Scanner{Name: "httpx"}
	options := s.BuildOptions()
	userAgent := flag.String("ua", DEFAULT_USER_AGENT, "Choose user-agent - use empty value for a random UA")
	scanner.ParseOptions(options)

	err := os.MkdirAll(options.Output, 0755)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// pass temporary file to binary instead of final path, as only finished files should be placed there
	file, err := os.CreateTemp("", s.Name)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer os.Remove(file.Name())

	args := []string{
		"-silent", "-no-fallback", "-pipeline", "-tech-detect",
		"-json", "-output", file.Name(),
		"-l", options.Input,
	}

	if *userAgent != "" {
		args = append(args, "-H", fmt.Sprintf("User-Agent: %s", *userAgent))
	}

	err = s.Exec(args...)
	if err != nil {
		log.Fatalf("Failed to run scanner: %v", err)
	}

	realOutputFile := path.Join(options.Output, "output.txt")
	outputFile, err := os.Create(realOutputFile)
	if err != nil {
		log.Fatalf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, file)
	if err != nil {
		log.Fatalf("Writing to output file failed: %v", err)
	}
}
