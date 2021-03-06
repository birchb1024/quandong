package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

var Version = "none"

// environMap - return the os.Environ as a map
func environMap() map[string]string {
	result := map[string]string{}
	e := os.Environ()
	for _, kev := range e {
		kv := strings.Split(kev, "=")
		result[kv[0]] = kv[1]
	}
	return result
}

// findProgramInPath - Look for the the real program we are usurping,
// skip this program when it's found in the PATH,
// return the path as a string.
func findProgramInPath(name string, currentExecutable string) (string, error) {
	PATH, ok := os.LookupEnv("PATH")
	if !ok {
		log.Fatal("ERROR: quandong No PATH environment variable")
	}
	pathList := strings.Split(PATH, ":")
	for _, dir := range pathList {
		log.Printf("INFO: quandong scanning path segment '%s'", dir)
		files, err := ioutil.ReadDir(dir) // some PATH entries are bogus - so ignore error ones
		if err != nil {
			log.Printf("INFO: quandong ReadDir error on '%s': %s", dir, err)
			continue
		}

		for _, f := range files {
			target := filepath.Join(dir, f.Name())
			if f.Name() == name && !f.Mode().IsDir() {
				target, err = filepath.EvalSymlinks(target)
				if err != nil {
					log.Fatal(err)
				}

				if target == currentExecutable {
					continue
				}
				log.Printf("INFO: quandong found target '%s'", target)
				return target, nil
			}
		}
	}
	return "", fmt.Errorf("ERROR: quandong target executable '%s' not found in PATH", name)
}
func main() {

	if runtime.GOOS != "linux" {
		panic("quandong runs on Linux only")
	}
	currentExecutable, err := filepath.EvalSymlinks("/proc/self/exe")

	name := filepath.Base(os.Args[0])
	target, err := findProgramInPath(name, currentExecutable)
	if err != nil {
		log.Fatal(err)
	}

	tempPrefix := fmt.Sprintf("quandong-%s-*.json", name)
	logFile, err := ioutil.TempFile(".", tempPrefix)
	if err != nil {
		log.Fatalf("ERROR: quandong cannot create temporary file: %s. %s", tempPrefix, err)
	}
	defer func() { _ = logFile.Close() }()

	var toLog = map[string]interface{}{
		"args":    os.Args,
		"environ": environMap(),
		"target":  target,
		"quandong": map[string]string{
			"version":    Version,
			"executable": currentExecutable,
		}}
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
