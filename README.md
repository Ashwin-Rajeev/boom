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

api-profiler -d 10 -g 20 https://www.google.com/

## output

<pre> API-Profiler running for 10s over the api: <font color="#4E9A06"> https://www.google.com/ </font>
 20 Active Concurrent connections!

|     Statistics     |     value     |
| ================================== |
 + Total   Reqs		<font color="#4E9A06"> 529 </font>
 + Fastest Reqs		<font color="#4E9A06"> 267.343248ms </font>
 + Slowest Reqs		<font color="#4E9A06"> 870.364788ms </font>
 + Average Reqs		<font color="#4E9A06"> 2.766071181s </font>
 + Error   Count        <font color="#4E9A06"> 0 </font>
― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―
|     Status Code    |     Count     |
| ================================== |
 + 1XX                  <font color="#4E9A06"> 0 </font>
 + 2XX                  <font color="#4E9A06"> 529 </font>
 + 3XX                  <font color="#4E9A06"> 0 </font>
 + 4XX                  <font color="#4E9A06"> 0 </font>
 + 5XX                  <font color="#4E9A06"> 0 </font>
 + Others               <font color="#4E9A06"> 0 </font>
― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―
</pre>

