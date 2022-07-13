package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/soumagar22/google-translate-project/cli"
)

var (
	sourceLang string
	targetLang string
	sourceText string
	wg         sync.WaitGroup
)

func init() {
	//init executes before main()
	//fmt.Println("Inside init function to check function order")
	flag.StringVar(&sourceLang, "s", "en", "Source Language[en]")
	flag.StringVar(&targetLang, "t", "fr", "Target Language[fr]")
	flag.StringVar(&sourceText, "st", "en", "Text to translate")
}

func main() {
	//main executes after init
	//fmt.Println("Inside main function to check function order")
	flag.Parse()

	if flag.NFlag() == 0 {
		fmt.Println("Options: ")
		flag.PrintDefaults()
		os.Exit(1)
	}

	strChan := make(chan string)

	reqBody := &cli.RequestBody{SourceLang: sourceLang, TargetLang: targetLang, SourceText: sourceText}

	wg.Add(1)
	go cli.RequestTranslate(reqBody, strChan, &wg)
	processedString := strings.ReplaceAll(<-strChan, "+", " ")
	fmt.Printf("%s \n", processedString)
	close(strChan)
	wg.Wait()
}
