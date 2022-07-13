package cli

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Jeffail/gabs"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const translateURL = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body *RequestBody, str chan string, wg *sync.WaitGroup) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", translateURL, nil)
	if err != nil {
		log.Fatalf("Something wrong happened %s", err)
		os.Exit(1)
	}

	query := req.URL.Query()
	//Adding google translate (gtx) as client
	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("q", body.SourceText)
	query.Add("dt", "t")

	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error occured while initiating request %v", err)
		os.Exit(1)
	}

	if resp != nil {
		if resp.StatusCode == http.StatusTooManyRequests {
			log.Fatalf("You have been rate limited")
			str <- "You have been rate limited, try again later"
			wg.Done()
		}

		parsedJson, err := gabs.ParseJSONBuffer(resp.Body)
		if err != nil {
			log.Fatalf("Error occured while parsing to json %s", err)
			os.Exit(1)
		}

		nestOne, err := parsedJson.ArrayElement(0)
		if err != nil {
			log.Fatalf("Error occured while retrieving first nestead object from api %s", err)
			os.Exit(1)
		}

		nestTwo, err := nestOne.ArrayElement(0)
		if err != nil {
			log.Fatalf("Error occured while retrieving second nestead object from api %s", err)
			os.Exit(1)
		}

		translatedString, err := nestTwo.ArrayElement(0)
		if err != nil {
			log.Fatalf("Error occured while retrieving third nestead object from api %s", err)
			os.Exit(1)
		}

		str <- translatedString.Data().(string)
		wg.Done()

	} else {
		log.Fatal("The response received from api was empty")
		os.Exit(1)
	}

	defer resp.Body.Close()
}
