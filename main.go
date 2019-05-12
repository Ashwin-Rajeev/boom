package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"
)

var (
	numberOfConcurrentConnections int
	requestURL                    string
	requestMethod                 string
	requestHeader                 map[string]string
	headerValues                  string
	requestDurationInSeconds      int
	requestBody                   string
	requestTimeOut                int
	help                          bool
)

func init() {
	flag.IntVar(&numberOfConcurrentConnections, "g", 5, "Number of concurrent connections")
	flag.StringVar(&requestMethod, "m", "GET", "Request method")
	flag.StringVar(&headerValues, "h", "", "header values seperated with ','")
	flag.StringVar(&requestBody, "b", "", "Request body file name (Relative path)")
	flag.IntVar(&requestDurationInSeconds, "d", 5, "Request duration")
	flag.IntVar(&requestTimeOut, "to", 1000, "Request time out in seconds")
	flag.BoolVar(&help, "help", false, "know more about the usage of api-profiler")
}

func main() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	flag.Parse()
	if !flag.Parsed() {
		log.Fatalln("[Info] Command line flags parsing failed, Please check the input")
	}
	requestHeader = make(map[string]string)
	if headerValues != "" {
		hv := strings.Split(headerValues, ",")
		for _, hd := range hv {
			header := strings.SplitN(hd, ":", 2)
			requestHeader[header[0]] = header[1]
		}
	}
	if help {
		fmt.Println("Usage: api-profiler <flags> <url>")
		flag.VisitAll(func(flag *flag.Flag) {
			fmt.Println("\t-"+flag.Name, "\t", flag.Usage, "(Default value = "+flag.DefValue+")")
		})
		return
	}

	if len(requestBody) > 0 {
		data, err := ioutil.ReadFile(requestBody)
		if err != nil {
			fmt.Println(fmt.Errorf("[Info] Could not read file %s: %s", requestBody, err.Error()))
			os.Exit(1)
		}
		requestBody = string(data)
	}
	requestURL = flag.Arg(0)

	if len(requestURL) == 0 {
		log.Fatalln("[Info] request url is invalid, Please check the input")
	}
	staticsChan := make(chan *APIStatus, numberOfConcurrentConnections)
	config := newAPIConfig(
		numberOfConcurrentConnections,
		requestURL,
		requestMethod,
		requestHeader,
		requestDurationInSeconds,
		requestBody,
		requestTimeOut,
		staticsChan,
	)
	fmt.Printf("✔ API-profiler running for %vs over the api: %v ✔\n", requestDurationInSeconds, requestURL)
	fmt.Printf("\t☺ %v goroutines running concurrently! Stay alert ☺\n\n", numberOfConcurrentConnections)
	for i := 0; i < numberOfConcurrentConnections; i++ {
		go config.request()
	}

	minions := 0
	statics := APIStatus{
		TotalDuration:  time.Minute,
		MinRequestTime: time.Minute,
		MaxRequestTime: time.Minute,
	}
	for minions < numberOfConcurrentConnections {
		select {
		case <-sigChannel:
			config.stop()
			fmt.Println("[Info] Api-profiler shutting down...")
			os.Exit(0)

		case s := <-staticsChan:
			statics.NumberOfRequests += s.NumberOfRequests
			statics.ErrorCount += s.ErrorCount
			statics.TotalDuration += s.TotalDuration
			statics.MaxRequestTime = s.MaxRequestTime
			statics.MinRequestTime = s.MinRequestTime
			statics.TotalResponseSize = s.TotalResponseSize
			minions++
		}
	}
	if statics.NumberOfRequests == 0 {
		fmt.Println("[Info] No request found")
		return
	}
	avgReqTime := statics.TotalDuration / time.Duration(statics.NumberOfRequests)

	fmt.Printf("❤ [Info]❤ Total Memory Read:\t%v Bytes\n", statics.TotalResponseSize)
	fmt.Printf("❤ [Info]❤ Total Requests:\t%v\n", statics.NumberOfRequests)
	fmt.Printf("❤ [Info]❤ Fastest Request:\t%v\n", statics.MinRequestTime)
	fmt.Printf("❤ [Info]❤ Slowest Request:\t%v\n", statics.MaxRequestTime)
	fmt.Printf("❤ [Info]❤ Average Request Time:\t%v\n", avgReqTime)
	fmt.Printf("❤ [Info]❤ Number of Errors:\t%v\n", statics.ErrorCount)
}
