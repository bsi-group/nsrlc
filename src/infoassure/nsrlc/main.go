package main

import (
	"bufio"
	"os"
	"log"
	"github.com/voxelbrain/goptions"
	"net/http"
	"net/url"
	"bytes"
	"strings"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

// ##### Variables ###########################################################

var (
	opt	*Options
)

// ##### Constants ###########################################################

const APP_TITLE string = "NSRL Client"
const APP_NAME string = "nsrlc"
const APP_VERSION string = "1.00"

const BATCH_SIZE int = 1000

// Formats for output
const (
	FORMAT_ALL 			= "a"
	FORMAT_IDENTIFIED 	= "i"
	FORMAT_UNIDENTIFIED = "u"
)

const DEFAULT_SERVER string = "127.0.0.1:8000"

// ##### Methods #############################################################

// Application entry point
func main() {

	fmt.Printf("\n%s (%s) %s\n\n", APP_TITLE, APP_NAME, APP_VERSION)

	// Set the default values for the command line options
	opt = new(Options)
	opt.Server = DEFAULT_SERVER
	opt.Format = FORMAT_ALL
	opt.BatchSize = BATCH_SIZE

	goptions.ParseAndFail(opt)

	// Lets make sure that the users input file actually exists
	if _, err := os.Stat(opt.InputFile); os.IsNotExist(err) {
		log.Fatal("Input file does not exist")
	}

	processInputFile()
}

// Import the import data file into the BTree
func processInputFile() {

	file, err := os.Open(opt.InputFile)
	if err != nil {
		log.Fatal("Error opening the input file: %v", err)
	}
	defer file.Close()

	fileOutput, err := os.Create(opt.OutputFile)
	if err != nil {
		log.Fatal("Error creating the output file: %v", err)
	}
	defer fileOutput.Close()

	// Output some CSV file headers
	switch (opt.Format) {
	case FORMAT_IDENTIFIED: // Found
		fileOutput.Write([]byte(fmt.Sprintf("%s\n", "Hash")))
	case FORMAT_UNIDENTIFIED: // Not Found
		fileOutput.Write([]byte(fmt.Sprintf("%s\n", "Hash")))
	case FORMAT_ALL: // All
		fileOutput.Write([]byte(fmt.Sprintf("%s,%s\n", "Hash", "Status")))
	}

	log.Print("Starting processing")

	client := &http.Client{}

	batchData := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}

		batchData = append(batchData, scanner.Text())
		if len(batchData) >= BATCH_SIZE {
			if sendBatch(client, batchData, fileOutput) == false {
				log.Print("The processing has terminated due to an error")
				os.Exit(-1)
			}

			// Clear the batch slice down ready for the next batch
			batchData = batchData[:0]
		}
	}

	// Send any remaining data
	if len(batchData) > 0 {
		if sendBatch(client, batchData, fileOutput) == false {
			log.Print("The processing has terminated due to an error")
			os.Exit(-1)
		}

		// Clear the batch slice down
		batchData = batchData[:0]
	}

	log.Print("Processing complete")
}

// Sends a HTTP POST request to the server. The POST data
// consists of the hashes separated by a hash character
func sendBatch(client *http.Client, data []string, fileOutput *os.File) bool {

	// Create the HTTP POST data delimited by the # character
	postData := strings.Join(data, "#")
	tempData := url.Values{}
	tempData.Set("hashes", postData)

	// Initialise the HTTP request and set the appropriate headers for a POST request
	r, _ := http.NewRequest("POST", "http://" + opt.Server + "/bulk", bytes.NewBufferString(tempData.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(tempData.Encode())))

	// Send the actual HTTP request
	resp, err := client.Do(r)
	if err != nil {
		log.Printf("Error performing HTTP lookup: %v", err)
		return false
	}

	// The only acceptable HTTP response code is 200
	if resp.StatusCode != 200 {
		log.Printf("Server returned an invalid status code: %v", resp.StatusCode)
		return false
	}

	// Get the HTTP response body
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading HTTP response body: %v", err)
		return false
	}

	// Turn the HTTP response text into actual struct data
	responseData := make([]JsonResult, 0)
	err = json.Unmarshal(contents, &responseData)
	if err != nil {
		log.Printf("Error unmarshalling JSON data: %v", err)
		return false
	}

	// Iterate through the slice of structs and output to the output file, using the desired output format
	for _, d := range responseData {
		if d.Exists == true {
			switch (opt.Format) {
			case FORMAT_IDENTIFIED: // Found
				fileOutput.Write([]byte(fmt.Sprintf("%s,%s\n", d.Hash, "FOUND")))
			case FORMAT_UNIDENTIFIED: // Not Found
			// Ignore
			case FORMAT_ALL: // All
				fileOutput.Write([]byte(fmt.Sprintf("%s,%s\n", d.Hash, "FOUND")))
			}
		} else {
			switch (opt.Format) {
			case FORMAT_IDENTIFIED: // Found
			// Ignore
			case FORMAT_UNIDENTIFIED: // Not Found
				fileOutput.Write([]byte(fmt.Sprintf("%s,%s\n", d.Hash, "NOT FOUND")))
			case FORMAT_ALL: // All
				fileOutput.Write([]byte(fmt.Sprintf("%s,%s\n", d.Hash, "NOT FOUND")))
			}
		}
	}

	return true
}