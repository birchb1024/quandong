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
				target := filepath.Join(dir, f.Name())
				target, err = filepath.EvalSymlinks(target)
				if err != nil {
					log.Fatal(err)
				}

				if target == currentExecutable {
					continue
				}
				return target, nil
			}
		}
	}
	return "", fmt.Errorf("ERROR: executable '%s' not found in PATH", name)
}
func main() {

	if runtime.GOOS != "linux" {
		panic("Program runs on Linux only")
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
		log.Fatalf("ERROR: Cannot create temporary file: %s. %s", tempPrefix, err)
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
