package main

import (
	"flag"
	"fmt"
    "net/url"
)

var verbose bool
var videoLinks []string
var maxConcurrent int

func init() {
    flag.BoolVar(&verbose, "verbose", false, "Be verbose")    
    flag.IntVar(&maxConcurrent, "goroutines", 8, "Maximum concurrent downloads")
}

func main() {
    parseArgs()

    fmt.Println("Links to be downloaded ", videoLinks)

    Download(&videoLinks)

}

func parseArgs() {
    flag.Parse()
    for _, arg := range flag.Args(){
        if !isValidURL(arg) {
            fmt.Println("Skipping ", arg)
            continue
        }
        videoLinks = append(videoLinks, arg)
    }
}

func isValidURL(str string) bool {
    _, err := url.ParseRequestURI(str)

    if err != nil {
        return false
    }

    u, err := url.Parse(str)

    if err != nil || u.Scheme != "https" {
        return false
    }

    return true
}
