package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var templetePath string
var dataPath string
var initPath string
var subject string

func init() {
	flag.StringVar(&templetePath, "tmpl", "./tmpl.html", "Path to the tmpl file to process")
	flag.StringVar(&dataPath, "data", "./data.txt", "Path to the data file to process")
	flag.StringVar(&initPath, "init", "./init.txt", "Path to the file with init data")
	flag.StringVar(&subject, "subject", "", "Subject of emails")
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()
	fmt.Println("Template path: " + templetePath)
	fmt.Println("Data path: " + dataPath)
	notifier := CreteNotifier(InitData{
		DataFilePath:       dataPath,
		TmplFilePath:       templetePath,
		SenderInfoFilePath: initPath,
		MessageSubject:     subject,
	})
	notifier.Init()
	notifier.ProcessTemplate()
	notifier.Send()
}
