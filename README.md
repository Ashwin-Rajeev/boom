# api-profiler
Golang package for testing the REST api

# Installing
go get github.com/Ashwin-Rajeev/api-profiler

# Usage

Usage: api-profiler <flags> <url>
        
        -b       Request body file name (Relative path) (Default value = )
        -d       Request duration (Default value = 5)
        -g       Number of concurrent connections (Default value = 5)
        -h       header values seperated with ',' (Default value = )
        -help    know more about the usage of api-profiler (Default value = false)
        -m       Request method (Default value = GET)
        -to      Request time out in seconds (Default value = 1000)

#  example

api-profiler -d 10 -g 5 https://www.google.com/

## output

✔ API-profiler running for 10s over the api: https://www.google.com/ ✔

☺ 5 goroutines running concurrently! Stay alert ☺

❤ [Info]❤ Total Memory Read:    137327 Bytes

❤ [Info]❤ Total Requests:       34

❤ [Info]❤ Fastest Request:      121.935ms

❤ [Info]❤ Slowest Request:      1m0s

❤ [Info]❤ Average Request Time: 10.992457138s

❤ [Info]❤ Number of Errors:     41
