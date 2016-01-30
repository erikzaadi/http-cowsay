package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
)

const (
	cowsDefaultPath   = "/usr/share/cowsay/cows"
	cowsayDefaultPath = "/usr/games/cowsay"
)

var (
	port               = 3000
	app                = "app"
	cowTemplateName    = "default"
	cowEnvPath         = os.Getenv("COWPATH")
	cowsayOverridePath = ""
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func load(cowfile, thoughts string) string {
	cowPath := cowsDefaultPath
	if cowEnvPath != "" {
		cowPath = cowEnvPath
	}
	cowsayPath := cowsayDefaultPath
	if cowsayOverridePath != "" {
		cowsayPath = cowsayOverridePath
	}
	var filePath string
	var absolute, _ = regexp.MatchString("/", cowfile)
	if absolute == true {
		filePath = fmt.Sprintf("%s.cow", cowfile)
	} else {
		filePath = fmt.Sprintf("%s/%s.cow", cowPath, cowfile)
	}
	zeArgs := []string{"-f", filePath, thoughts}
	cmdOut, err := exec.Command(cowsayPath, zeArgs...).Output()

	checkError(err)

	return string(cmdOut)
}

func handler(w http.ResponseWriter, r *http.Request) {
	thoughts := fmt.Sprintf("Hello from %s", app)
	cowString := load(cowTemplateName, thoughts)
	fmt.Fprint(w, cowString)
}

func main() {
	flag.IntVar(&port, "p", 3000, "http port to serve")
	flag.StringVar(&app, "a", "app", "app name")
	flag.StringVar(&cowTemplateName, "c", "default", "cow name")
	flag.StringVar(&cowsayOverridePath, "e", "", "cowsay executable pat")
	flag.Parse()
	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
