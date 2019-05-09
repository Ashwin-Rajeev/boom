package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
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
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)

	flag.IntVar(&numberOfConcurrentConnections, "-g", 5, "Number of concurrent connections")
	flag.StringVar(&requestMethod, "-m", "GET", "Request method")
	flag.StringVar(&headerValues, "-h", "", "header values seperated with ','")
	flag.StringVar(&requestBody, "-b", "", "Request body file name (Relative path)")
	flag.IntVar(&requestDurationInSeconds, "-d", 5, "Request duration")
	flag.IntVar(&requestTimeOut, "-to", 100, "Request time out in seconds")
	flag.BoolVar(&help, "-help", false, "know more about the usage of api-profiler")
}

func main() {
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

	config := newAPIConfig(
		numberOfConcurrentConnections,
		requestURL,
		requestMethod,
		requestHeader,
		requestDurationInSeconds,
		requestBody,
		requestTimeOut,
	)

	for i := 0; i < numberOfConcurrentConnections; i++ {
		go config.request()
	}
}
