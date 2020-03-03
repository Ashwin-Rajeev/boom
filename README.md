# Boom [![Build Status](https://travis-ci.com/Ashwin-Rajeev/boom.svg?branch=master)](https://travis-ci.com/Ashwin-Rajeev/boom) [![GoDoc](https://godoc.org/github.com/Ashwin-Rajeev/boom?status.svg)](https://godoc.org/github.com/Ashwin-Rajeev/boom) [![Go Report Card](https://goreportcard.com/badge/github.com/Ashwin-Rajeev/boom)](https://goreportcard.com/report/github.com/Ashwin-Rajeev/boom)
Boom is a HTTP(S) benchmarking tool written in Golang. Because of its lightning fast performance you can test any real world applications using this wonderfull benchmarking tool.

### Coded With Language
* [Golang](https://golang.org/)


<!-- GETTING STARTED -->
## Getting Started

### Installing
To get the latest and greatest run:

```
go get -u github.com/Ashwin-Rajeev/boom
```

## Usage
```
Usage:  boom [<flags>] <url>

        -help    know more about the usage of boom (Default value = false)
        -d       Request duration (Default value = 5)
        -g       Number of concurrent connections (Default value = 5)
        -b       Request body file name (Relative path) (Default value = )
        -h       header values separated with ';' (Default value = )
        -m       Request method (Default value = GET)
        -mul     Use multiple request mode (Default value = false)
        -to      Request time out in seconds (Default value = 1000)

```

##  Example
 
### Single request
<pre>
> boom -d 10 -g 20 https://www.google.com/
  Boom running for 10s over the api: <font color="#4E9A06"> https://www.google.com/ </font>
  20 Active Concurrent connections!
10 / 10 [--------------------------------------------------------------] 100.00%

|     Statistics     |     value     |
| ================================== |
 + Total   Reqs		<font color="#4E9A06"> 1424 </font>
 + Fastest Reqs		<font color="#4E9A06"> 128.59929ms </font>
 + Slowest Reqs		<font color="#4E9A06"> 273.003861ms </font>
 + Average Reqs		<font color="#4E9A06"> 1.026262071s </font>
 + Error   Count        <font color="#4E9A06"> 0 </font>
― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―
|     Status Code    |     Count     |
| ================================== |
 + 1XX                  <font color="#4E9A06"> 0 </font>
 + 2XX                  <font color="#4E9A06"> 1424 </font>
 + 3XX                  <font color="#4E9A06"> 0 </font>
 + 4XX                  <font color="#4E9A06"> 0 </font>
 + 5XX                  <font color="#4E9A06"> 0 </font>
 + Others               <font color="#4E9A06"> 0 </font>
― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―
</pre>

### Single POST request with body
<pre>
> boom -d 10 -g 20 -m POST -b body.json https://www.sample.com/
  Boom running for 10s over the api: <font color="#4E9A06"> https://www.google.com/ </font>
  20 Active Concurrent connections!
10 / 10 [--------------------------------------------------------------] 100.00%

|     Statistics     |     value     |
| ================================== |
 + Total   Reqs		<font color="#4E9A06"> 1424 </font>
 + Fastest Reqs		<font color="#4E9A06"> 128.59929ms </font>
 + Slowest Reqs		<font color="#4E9A06"> 273.003861ms </font>
 + Average Reqs		<font color="#4E9A06"> 1.026262071s </font>
 + Error   Count        <font color="#4E9A06"> 0 </font>
― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―
|     Status Code    |     Count     |
| ================================== |
 + 1XX                  <font color="#4E9A06"> 0 </font>
 + 2XX                  <font color="#4E9A06"> 1600 </font>
 + 3XX                  <font color="#4E9A06"> 0 </font>
 + 4XX                  <font color="#4E9A06"> 0 </font>
 + 5XX                  <font color="#4E9A06"> 0 </font>
 + Others               <font color="#4E9A06"> 0 </font>
― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―
</pre>

### Multiple requests using multi request mode
<pre>
> boom  --mul -d 10 -g 20 -b body.json
  Boom running for 10s over the api: <font color="#4E9A06"> </font>
  20 Active Concurrent connections!
10 / 10 [--------------------------------------------------------------] 100.00%

|     Statistics     |     value     |
| ================================== |
 + Total   Reqs		<font color="#4E9A06"> 1424 </font>
 + Fastest Reqs		<font color="#4E9A06"> 128.59929ms </font>
 + Slowest Reqs		<font color="#4E9A06"> 273.003861ms </font>
 + Average Reqs		<font color="#4E9A06"> 1.026262071s </font>
 + Error   Count        <font color="#4E9A06"> 0 </font>
― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―
|     Status Code    |     Count     |
| ================================== |
 + 1XX                  <font color="#4E9A06"> 0 </font>
 + 2XX                  <font color="#4E9A06"> 1600 </font>
 + 3XX                  <font color="#4E9A06"> 0 </font>
 + 4XX                  <font color="#4E9A06"> 0 </font>
 + 5XX                  <font color="#4E9A06"> 0 </font>
 + Others               <font color="#4E9A06"> 0 </font>
― ― ― ― ― ― ― ― ― ― ―― ― ― ― ― ― ― ― ―
</pre>

sample_request.json
```

[
    {
        "url": "http://www.sample.com",
        "method": "GET",
        "body": {
            "a": "1",
            "b": "2"
        },
        "header": {
            "key": "value"
        }
    },
    {
        "url": "http://www.sample.com",
        "method": "POST",
        "body": {
            "a": 1,
            "b": 2
        },
        "header": {
            "key1": "value1",
            "key2": "value2"
        }
    }
]
```
## Prerequisites

Golang should be installed on your computer, boom is compatible with go1.10 and above


<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/Ashwin-Rajeev/boom/issues) for a list of proposed features (and known issues).



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.



<!-- CONTACT -->
## Contact

Email   : [ashwinrajeev916@gmail.com](ashwinrajeev916@gmail.com)   
GitHub  : [https://github.com/Ashwin-Rajeev](https://github.com/Ashwin-Rajeev)   
Website : [https://ashwinrajeev.com](https://ashwinrajeev.com)
