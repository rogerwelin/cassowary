[![Go Report Card](https://goreportcard.com/badge/github.com/rogerwelin/cassowary)](https://goreportcard.com/report/github.com/rogerwelin/cassowary)
[![Build Status](https://travis-ci.org/rogerwelin/cassowary.svg?branch=master)](https://travis-ci.org/rogerwelin/cassowary)


**Cassowary** is a modern HTTP(S), intuitive & cross-platform load testing tool built in Go for developers, testers and sysadmins. Cassowary draws inspiration from awesome projects like k6, ab & httpstat.


Features  
--------

- **2 Load Testing modes**: one standard and one spread mode where URL Paths can be specified from a file (ideal if you want to hit several underlying microservices)
- **CI Friendly**: Well-suited to be part of a CI pipeline step
- **Flexible metrics**: Prometheus metrics (pushing metrics to Prometheus PushGateway), JSON file
- **Configurable**: Able to pass in arbitrary HTTP headers
- **Cross Platform**: One single pre-built binary for Windows, Mac OSX and Windows


