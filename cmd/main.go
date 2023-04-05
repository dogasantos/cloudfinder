package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"sync"
	"time"

	cloudfinder "github.com/dogasantos/cloudfinder/pkg/runner"
)


// This is a nice comment to make lint happy. hello lint, i'm here!
type Options struct {
	TargetListFile		string
	Version				bool
	Verbose				bool
}

var version = "0.1"

func parseOptions() *Options {
	options := &Options{}
	flag.StringVar(&options.TargetListFile, 		"l", "target.txt", "Target file (fqdn)")
	flag.BoolVar(&options.Version, 					"i", false, "Version info")
	flag.BoolVar(&options.Verbose, 					"v", false, "Verbose mode")
	flag.Parse()
	return options
}

func main() {

	options := parseOptions()
	if options.Version {
		fmt.Println(version)
	}

	if options.TargetListFile != "" {
		if options.Verbose == true {
			fmt.Printf("[+] cloudfinder v%s\n",version)
		}
		TargetFilestream, _ := ioutil.ReadFile(options.TargetListFile)
		targetContent := string(TargetFilestream)
		targets := strings.Split(targetContent, "\n") 
		
		if options.Verbose == true {
			fmt.Printf("  + Targets loaded: %d\n",len(targets))
		}

		wg := new(sync.WaitGroup)
		routinescounter := 0
		for _, target := range targets {
			target = strings.ReplaceAll(target, " ", "")
			if len(target) > 1 {
				wg.Add(1)
				go cloudfinder.start(target, options.Verbose, wg)
				if routinescounter == int(math.Round(float64(len(targets)) / 10))  {
					time.Sleep(5 * time.Second)
					routinescounter = 0
				} else {
					routinescounter = routinescounter+1
				}
			}
		}
		wg.Wait()
	}

}