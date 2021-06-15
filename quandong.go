package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"
)

var Version = "none"

// findProgramInPath - Look for the the real program we are usurping,
// return the path to it as a string.
func findProgramInPath(name string) (string, error) {
	PATH, ok := os.LookupEnv("PATH")
	if !ok {
		log.Fatal("ERROR: No PATH environment variable")
	}
	pathList := strings.Split(PATH, ":")
	for _, dir := range pathList {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			if f.Name() == name && !f.Mode().IsDir() && f.Mode().IsRegular() && (f.Mode()&0111) != 0 {
				return fmt.Sprintf("%s/%s", dir, f.Name()), nil
			}
		}
	}
	return "", fmt.Errorf("ERROR: executable '%s' not found in PATH", name)
}
func main() {

	fmt.Printf("%#v\n", os.Args)
	a0 := strings.Split(os.Args[0], "/")
	name := a0[len(a0)-1]
	target, err := findProgramInPath(name)
	if err != nil {
		log.Fatal(err)
	}

	tempPrefix := fmt.Sprintf("quandong-%s-", name)
	logFile, err := ioutil.TempFile(".", tempPrefix)
	if err != nil {
		log.Fatalf("ERROR: Cannot create temporary file: %s. %s", tempPrefix, err)
	}
	defer func() { _ = logFile.Close() }()

	var toLog = map[string]interface{}{"Args": os.Args, "Environ": os.Environ(), "Target": target, "quandong": map[string]string{"Version": Version}}
	js, err := json.MarshalIndent(toLog, "", "  ")
	n, err := fmt.Fprintln(logFile, string(js))
	if n == 0 || err != nil {
		log.Fatal(err)
	}
	err = syscall.Exec(target, os.Args, os.Environ())
	if err != nil {
		log.Fatal(err)
	}
}
