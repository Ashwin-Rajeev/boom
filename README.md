# api-profiler
Golang package for testing the REST api

# Installing
> go get -u github.com/Ashwin-Rajeev/api-profiler

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

 API-Profiler running for 10s over the api:  https://www.google.com/ 
 5 Active Concurrent connections!

|     Statistics     |     value     |
| ================================== |
 + Total   Reqs          502 
 + Fastest Reqs          83.435901ms 
 + Slowest Reqs          262.846001ms 
 + Average Reqs          817.296542ms 
 + Error   Count         0 
― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―
|     Status Code    |     Count     |
| ================================== |
 + 1XX                   0 
 + 2XX                   502 
 + 3XX                   0 
 + 4XX                   0 
 + 5XX                   0 
 + Others                0 
― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―

