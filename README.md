<p align="center"><a href="https://github.com/rogerwelin/cassowary"><img src="cass-logo.png" alt="cassowary"></a></p>
<p align="center">
  <a href="https://goreportcard.com/badge/github.com/rogerwelin/cassowary"><img src="https://goreportcard.com/badge/github.com/rogerwelin/cassowary" alt="Go Report Card"></a>
  <a href="https://travis-ci.org/rogerwelin/cassowary"><img src="https://travis-ci.org/rogerwelin/cassowary.svg?branch=master" alt="Build status"></a>
  <a href="https://github.com/rogerwelin/cassowary/blob/master/LICENSE"><img src="https://img.shields.io/github/license/rogerwelin/cassowary" alt="License"></a>
  <a href="https://github.com/rogerwelin/cassowary/blob/master/go.mod"><img src="https://img.shields.io/github/go-mod/go-version/rogerwelin/cassowary" alt="Go version"></a>
  <a href="https://github.com/rogerwelin/cassowary/releases"><img src="https://img.shields.io/github/v/release/rogerwelin/cassowary.svg" alt="Current Release"></a>
</p>


**Cassowary** is a modern HTTP/S, intuitive & cross-platform load testing tool built in Go for developers, testers and sysadmins. Cassowary draws inspiration from awesome projects like k6, ab & httpstat.

---

Features  
--------

- **2 Load Testing modes**: one standard and one spread mode where URL Paths can be specified from a file (ideal if you want to hit several underlying microservices)
- **CI Friendly**: Well-suited to be part of a CI pipeline step
- **Flexible metrics**: Prometheus metrics (pushing metrics to Prometheus PushGateway), JSON file
- **Configurable**: Able to pass in arbitrary HTTP headers
- **Cross Platform**: One single pre-built binary for Linux, Mac OSX and Windows

<img src="https://i.imgur.com/geJykYH.gif" />


Installation  
--------

Grab a pre-built binary from the [GitHub Releases page](https://github.com/rogerwelin/cassowary/releases). You can optionally put the **cassowary** binary in your `PATH` so you can run cassowary from any location


Running Cassowary  
--------

Example running **cassowary** against www.example.com with 100 requests spread out over 10 concurrent users:

```
$ ./cassowary run -u http://www.example.com -c 10 -n 100

Starting Load Test with 100 requests using 10 concurrent users

 100% |████████████████████████████████████████| [1s:0s]            1.256773616s


 TCP Connect.....................: Avg/mean=101.90ms 	Median=102.00ms	p(95)=105ms
 Server Processing...............: Avg/mean=100.18ms 	Median=100.50ms	p(95)=103ms
 Content Transfer................: Avg/mean=0.01ms 	Median=0.00ms	p(95)=0ms

Summary:
 Total Req.......................: 100
 Failed Req......................: 0
 DNS Lookup......................: 115.00ms
 Req/s...........................: 79.57
```

Example running **cassowary** in file slurp mode where all URL paths are specified from an external file (which can also be fetched from http if specified):

```
$ ./cassowary run-file -u http://localhost:8000 -c 10 -f urlpath.txt

Starting Load Test with 0 requests using 10 concurrent users

 100% |████████████████████████████████████████| [0s:0s]            599.467161ms


 TCP Connect.....................: Avg/mean=1.80ms 	Median=2.00ms	p(95)=3ms
 Server Processing...............: Avg/mean=0.90ms 	Median=0.00ms	p(95)=3ms
 Content Transfer................: Avg/mean=0.00ms 	Median=0.00ms	p(95)=0ms

Summary:
 Total Req.......................: 3925
 Failed Req......................: 0
 DNS Lookup......................: 2.00ms
 Req/s...........................: 6547.48
```

Project Status & Contribute  
--------

Feature-wise I'm quite happy with `cassowary` and hence new features will only be added on a feature request and bug basis. To request a feature create a new issue with the label `feature-request`. For bugs use the label `bug`.
