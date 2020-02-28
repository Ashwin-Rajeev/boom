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

	"github.com/cheggaaa/pb/v3"
	clr "github.com/fatih/color"
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
	flag.StringVar(&headerValues, "h", "", "header values separated with ';'")
	flag.StringVar(&requestBody, "b", "", "Request body file name (Relative path)")
	flag.IntVar(&requestDurationInSeconds, "d", 5, "Request duration")
	flag.IntVar(&requestTimeOut, "to", 2000, "Request time out in seconds")
	flag.BoolVar(&help, "help", false, "know more about the usage of boom")
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
		hv := strings.Split(headerValues, ";")
		for _, hd := range hv {
			header := strings.SplitN(hd, ":", 2)
			requestHeader[header[0]] = header[1]
		}
	}
	if help {
		fmt.Println("Usage: boom [<flags>] <url>")
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
		log.Fatalln("[Info] Requested url is invalid, Please check the input")
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
	fmt.Printf(" Boom running for %vs over the api: ", requestDurationInSeconds)
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", requestURL)
	clr.Unset()
	fmt.Printf(" %v Active Concurrent connections!\n", numberOfConcurrentConnections)

	// progressbar configuration
	bar := pb.Simple.Start(requestDurationInSeconds)
	go func(rd int, bar *pb.ProgressBar) {
		d1, _ := time.ParseDuration(fmt.Sprintf("%vs", rd))
		timeout := time.After(d1)
		for {
			select {
			case <-timeout:
				bar.Finish()
				return
			default:
				bar.Increment()
				time.Sleep(time.Second)
			}
		}
	}(requestDurationInSeconds, bar)

	for i := 0; i < numberOfConcurrentConnections; i++ {
		go config.request()
	}

	minions := 0
	statics := APIStatus{
		TotalDuration:  time.Minute,
		MinRequestTime: time.Minute,
		MaxRequestTime: time.Minute,
		StatusCodes:    &StatusCodes{},
	}
	for minions < numberOfConcurrentConnections {
		select {
		case <-sigChannel:
			config.stop()
			os.Exit(0)

		case s := <-staticsChan:
			statics.NumberOfRequests += s.NumberOfRequests
			statics.ErrorCount += s.ErrorCount
			statics.TotalDuration += s.TotalDuration
			statics.MaxRequestTime = s.MaxRequestTime
			statics.MinRequestTime = s.MinRequestTime
			statics.TotalResponseSize = s.TotalResponseSize
			statics.StatusCodes.OneXX += s.StatusCodes.OneXX
			statics.StatusCodes.TwoXX += s.StatusCodes.TwoXX
			statics.StatusCodes.ThreeXX += s.StatusCodes.ThreeXX
			statics.StatusCodes.FourXX += s.StatusCodes.FourXX
			statics.StatusCodes.FiveXX += s.StatusCodes.FiveXX
			statics.StatusCodes.Others += s.StatusCodes.Others
			minions++
		}
	}
	if statics.NumberOfRequests == 0 {
		fmt.Println("[Info] No request found")
		return
	}

	defer func() {
		printResult(statics)
	}()
}

// printResult out the result into console
func printResult(statics APIStatus) {
	if statics.NumberOfRequests == 0 {
		statics.NumberOfRequests = 1
	}
	statics.AvgReqTime = statics.TotalDuration / time.Duration(statics.NumberOfRequests)
	fmt.Printf("\n")
	fmt.Printf(`|     Statistics     |     value     |`)
	fmt.Printf("\n")
	fmt.Printf(`| ================================== |`)
	fmt.Printf("\n")
	fmt.Printf(` + Total   Reqs`)
	fmt.Printf("\t\t")
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.NumberOfRequests)
	clr.Unset()
	fmt.Printf(" + Fastest Reqs\t\t")
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.MinRequestTime)
	clr.Unset()
	fmt.Printf(" + Slowest Reqs\t\t")
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.MaxRequestTime)
	clr.Unset()
	fmt.Printf(" + Average Reqs\t\t")
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.AvgReqTime)
	clr.Unset()
	fmt.Printf(` + Error   Count`)
	fmt.Printf(`        `)
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.ErrorCount)
	clr.Unset()
	fmt.Printf(`― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―`)
	fmt.Printf("\n")
	fmt.Printf(`|     Status Code    |     Count     |`)
	fmt.Printf("\n")
	fmt.Printf(`| ================================== |`)
	fmt.Printf("\n")
	fmt.Printf(` + 1XX`)
	fmt.Printf(`                  `)
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.StatusCodes.OneXX)
	clr.Unset()
	fmt.Printf(" + 2XX")
	fmt.Printf(`                  `)
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.StatusCodes.TwoXX)
	clr.Unset()
	fmt.Printf(" + 3XX")
	fmt.Printf(`                  `)
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.StatusCodes.ThreeXX)
	clr.Unset()
	fmt.Printf(" + 4XX")
	fmt.Printf(`                  `)
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.StatusCodes.FourXX)
	clr.Unset()
	fmt.Printf(" + 5XX")
	fmt.Printf(`                  `)
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.StatusCodes.FiveXX)
	clr.Unset()
	fmt.Printf(" + Others")
	fmt.Printf(`               `)
	clr.Set(clr.FgGreen)
	fmt.Printf(" %v \n", statics.StatusCodes.Others)
	clr.Unset()
	fmt.Printf(`― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―`)
	fmt.Printf("\n")
}
