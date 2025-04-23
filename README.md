<p align="center"><a href="https://github.com/rogerwelin/cassowary"><img src="https://i.imgur.com/V9BaioC.png" alt="cassowary"></a></p>
<p align="center">
  <a href="https://goreportcard.com/badge/github.com/rogerwelin/cassowary"><img src="https://goreportcard.com/badge/github.com/rogerwelin/cassowary?style=for-the-badge&logo=go" alt="Go Report Card"></a>
  <a href="https://github.com/rogerwelin/cassowary/actions/workflows/pullrequest.yaml"><img src="https://img.shields.io/github/actions/workflow/status/rogerwelin/cassowary/pullrequest.yaml?branch=master&style=for-the-badge&logo=github" alt="Build status"></a>
  <a href="https://github.com/avelino/awesome-go"><img src="https://awesome.re/mentioned-badge.svg" height="28" alt="Mentioned in Awesome Go"></a>
  <a href="https://github.com/rogerwelin/cassowary/blob/master/go.mod"><img src="https://img.shields.io/github/go-mod/go-version/rogerwelin/cassowary?style=for-the-badge&logo=go" alt="Go version"></a>
  <a href="https://github.com/rogerwelin/cassowary/releases"><img src="https://img.shields.io/github/v/release/rogerwelin/cassowary?style=for-the-badge&logo=github&color=orange" alt="Current Release"></a>
  <a href="https://godoc.org/github.com/rogerwelin/cassowary"><img src="https://godoc.org/github.com/rogerwelin/cassowary?status.svg" height="28" alt="godoc"></a>
  <a href="https://gocover.io/github.com/rogerwelin/cassowary/pkg/client"><img src="https://gocover.io/_badge/github.com/rogerwelin/cassowary/pkg/client" height="28" alt="Coverage"></a>
  <a href="https://github.com/rogerwelin/cassowary/blob/master/LICENSE"><img src="https://img.shields.io/badge/LICENSE-MIT-orange?style=for-the-badge" alt="License"></a>
</p>

English | [中文](README-ZH.md)


**Cassowary** is a modern HTTP/S, intuitive & cross-platform load testing tool built in Go for developers, testers and sysadmins. Cassowary draws inspiration from awesome projects like k6, ab & httpstat.

---

Toc
----

- [Features](#features)
- [Installation](#installation)
- [Running Cassowary](#running-cassowary)
  * [Regular Load Test](#regular-load-test)
  * [File Slurp Mode](#file-slurp-mode)
  * [Exporting Metrics to File](#exporting-metrics-to-file)
  * [Exporting Metrics to Prometheus](#exporting-metrics-to-prometheus)
  * [Exporting Metrics to Cloudwatch](#exporting-metrics-to-cloudwatch)
  * [Histogram](#histogram)
  * [Boxplot](#boxplot)
  * [Load Test with POST Data](#load-test-with-post-data)
  * [Specifying a Duration](#specifying-a-duration-for-the-load-test)
  * [Adding HTTP Headers](#adding-http-headers)
  * [Disabling HTTP keep-alive](#disabling-http-keep-alive)
  * [x509 Authentication](#x509-authentication)
  * [Distributed Load Testing](#distributed-load-testing)
- [Importing cassowary as a module](#importing-cassowary-as-a-modulelibrary)
- [Versioning](#versioning)
- [Contributing](#contributing)


Features  
--------

📌 &nbsp; **2 Load Testing modes**: one standard and one spread mode where URL Paths can be specified from a file (ideal if you want to hit several underlying microservices)  
📌 &nbsp;  **CI Friendly**: Well-suited to be part of a CI pipeline step  
📌 &nbsp; **Flexible metrics**: Cloudwatch metrics, Prometheus metrics (pushing metrics to Prometheus PushGateway), JSON file  
📌  &nbsp; **Configurable**: Able to pass in arbitrary HTTP headers, able to configure the HTTP client  
📌  &nbsp; **Supports GET, POST, PUT & PATCH** - POST, PUT and PATCH data can be defined in a file  
📌  &nbsp; **Cross Platform**: One single pre-built binary for Linux, Mac OSX and Windows  
📌  &nbsp; **Importable** - Besides the CLI tool cassowary can be imported as a module in your Go app  
📌 &nbsp; **Visualizations** - Cassowary can export request data as a histogram and boxplot in png format  

<img src="https://imgur.com/ac8F8eD.gif" />

Installation  
--------

Grab a pre-built binary from the [GitHub Releases page](https://github.com/rogerwelin/cassowary/releases). You can optionally put the **cassowary** binary in your `PATH` so you can run cassowary from any location. Alternative you can:

### Homebrew on Mac OSX  
You can install **cassowary** using the Homebrew package manager on Mac:

```bash
$ brew update && brew install cassowary
```

### Docker  

You can run **cassowary** directly by using its official Docker image:  

```bash
$ docker run rogerw/cassowary:v0.14.1 -u http://www.example.com -c 1 -n 10
```

For local development:

```bash
$ GOOS=linux go build -o dist/docker/cassowary cmd/cassowary/*.go
$ docker build -f dist/docker/Dockerfile -t test_cassowary dist/docker
$ docker run test_cassowary -u http://www.example.com -c 1 -n 10
```
Without argument when doing "docker run" it prints help usage

### ArchLinux/Manjaro

You can install the development version of Cassowary from the [AUR](https://aur.archlinux.org/packages/cassowary-git):

```bash
yay -S cassowary-git
```

Or build and install it manually:

```bash
git clone https://aur.archlinux.org/cassowary-git.git
cd cassowary-git
makepkg -si
```

### Nix/NixOS

Cassowary can be installed via the [Nix](https://nixos.org) package manager.
```
nix-env -iA cassowary
```

### CentOS/RHEL (RPM)

If you want to roll out your own RPM you can use the spec file [cassowary.spec](https://github.com/rogerwelin/cassowary/blob/master/dist/rpm/cassowary.spec) to build an RPM package



Running Cassowary  
--------

### Regular Load Test  
Example running **cassowary** against www.example.com with 100 requests spread out over 10 concurrent users:

```bash
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

### File Slurp Mode  
Example running **cassowary** in file slurp mode where all URL paths are specified from an external file (which can also be fetched from http if specified). By default cassowary will, without the -n flag specified, make one request per path specified in the file. However with the -n flag you can also specify how many request you want cassowary to generate against those URL paths. Example:

```bash
$ ./cassowary run -u http://localhost:8000 -c 1 -f urlpath.txt

# NOTE: from v0.10.0 and below file slurp mode had it's own command
# $ ./cassowary run-file -u http://localhost:8000 -c 1 -f urlpath.txt

Starting Load Test with 5 requests using 1 concurrent users

[ omitted ]


$ ./cassowary run -u http://localhost:8000 -c 10 -n 100 -f urlpath.txt

Starting Load Test with 100 requests using 10 concurrent users

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

### Exporting Metrics to File/CSV  
Example exporting **cassowary** json metrics to a file:

```bash
$ ./cassowary run --json-metrics --json-metrics-file=metrics.json -u http://localhost:8000 -c 125 -n 100000

Starting Load Test with 100000 requests using 125 concurrent users

 100% |████████████████████████████████████████| [0s:0s]            984.9862ms


 TCP Connect.....................: Avg/mean=-0.18ms     Median=0.00ms   p(95)=1ms
 Server Processing...............: Avg/mean=0.16ms      Median=0.00ms   p(95)=1ms
 Content Transfer................: Avg/mean=0.01ms      Median=0.00ms   p(95)=0ms

Summary:
 Total Req.......................: 100000
 Failed Req......................: 0
 DNS Lookup......................: 2.00ms
 Req/s...........................: 101524.27
```

> If `json-metrics-file` flag is missing then the default filename is `out.json`.

If you rather want the raw data of all the request (in CSV format) you can export it using the `--raw-output` flag

> outputted filename is `raw.csv`

### Exporting Metrics to Prometheus
Example exporting **cassowary** metrics to Prometheus by supplying an Prometheus PushGatway URL:

```bash
$ ./cassowary run -u http://localhost:8000 -c 125 -n 100000 -p http://pushgatway:9091

Starting Load Test with 100000 requests using 125 concurrent users

[ omitted for brevity ]

```

### Exporting Metrics to Cloudwatch
**Cassowary** can export metrics to AWS Cloudwatch just by adding the *--cloudwatch* flag without a value. Take note that you will need to tell Cassoway which AWS Region you want to use. The easiest way is using an environment variable as shown below:

```bash
$ export AWS_REGION=eu-north-1 && ./cassowary run -u http://localhost:8000 -c 125 -n 100000 --cloudwatch

Starting Load Test with 100000 requests using 125 concurrent users

[ omitted for brevity ]

```

### Histogram  
Supply the *--histogram* flag without a value to let cassowary compute and output a histogram in png format (hist.png in current dir) for the total duration of the requests in ms. Example: 

<img src="https://i.imgur.com/VLEsVOY.png" width="300" height="300" />

### Boxplot  
Supply the *--boxplot* flag without a value to tell cassowary to generate a boxplot visualization in png format (boxplot.png in current dir) for the total duration of the requests in ms

### Load Test with POST Data  
Example hitting a POST endpoint where POST json data is defined in a file:

```bash
$ ./cassowary run -u http://localhost:8000/add-user -c 10 -n 1000 --postfile user.json

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

### Load Test with PATCH Data  
Example hitting a PATCH endpoint where PATCH json data is defined in a file:

```bash
$ ./cassowary run -u http://localhost:8000/add-user -c 5 -n 200 --patchfile user.json

Starting Load Test with 200 requests using 5 concurrent users

[ omitted for brevity ]

```

### Specifying a Duration for the Load Test  
Example specifying a *duration* for your load test, in the command below we specify that we want send 100 requests over a duration of 30 seconds:

```bash
$ ./cassowary run -u http://localhost:8000 -n 100 -d 30

Starting Load Test with 100 requests using 1 concurrent users

[ omitted for brevity ]

```

### Adding HTTP Headers  
Example adding an HTTP header when running **cassowary**

```bash
$ ./cassowary run -u http://localhost:8000 -c 10 -n 1000 -H 'Host: www.example.com'

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

### Disabling HTTP keep-alive  
Example disabling http keep-alive (by default keep-alive are enabled):

```bash
$ ./cassowary run -u http://localhost:8000 -c 10 -n 1000 --disable-keep-alive

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

### Specifying ca certificate
Example specifying ca certificate

```bash
$ ./cassowary run -u http://localhost:8000 -c 10 -n 1000 --ca /path/to/ca.pem

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

### x509 Authentication
Example specifying client authentication for mTLS

```bash
$ ./cassowary run -u https://localhost:8443 -c 10 -n 1000 --cert /path/to/client.pem --key /path/to/client-key.pem --ca /path/to/ca.pem --renegotiation once

Starting Load Test with 1000 requests using 10 concurrent users, allows the server to renegotiate client certificates once. This is needed for servers that optionally request client certificates.

[ omitted for brevity ]

```

### Distributed Load Testing  
You can run cassowary in a distributed way if you need to scale up the load testing on more than one machine. Simplest way to do that is having access to a Kubernetes cluster. Use the batch type in Kubernetes and in the *spec.paralellism* key you can specify how many instances of cassowary you want to run simultaneously:

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: cassowary
spec:
  parallelism: 10
  template:
    spec:
      containers:
      - command: ["-u", "http://my-microservice.com:8000", "-c", "1", "-n", "10"]
        image: rogerw/cassowary:v0.14.1
        name: cassowary
      restartPolicy: Never
```

Just apply this yaml like so:

```bash
$ kubectl apply -f cassowary.yaml
```


Importing cassowary as a module/library
--------

Cassowary can be imported and used as a module in your Go app. Start by fetching the dependency by using go mod:

```bash
$ go mod init test && go get github.com/rogerwelin/cassowary/pkg/client
```

And below show a simple example on how to trigger a load test from your code and printing the results:

```go
package main

import (
        "encoding/json"
	"fmt"

	"github.com/rogerwelin/cassowary/pkg/client"
)

func main() {
	cass := &client.Cassowary{
		BaseURL:               "http://www.example.com",
		ConcurrencyLevel:      1,
		Requests:              10,
		DisableTerminalOutput: true,
	}
	metrics, err := cass.Coordinate()
	if err != nil {
		panic(err)
	}

        // print results
	fmt.Printf("%+v\n", metrics)

        // or print as json
	jsonMetrics, err := json.Marshal(metrics)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonMetrics))
}
```

More library examples [can be found here](docs/LIBRARY.md)


Versioning
--------

Cassowary follows semantic versioning. The public library (pkg/client) may break backwards compatibility until it hits a stable v1.0.0 release.

Contributing
--------

Contributions are welcome! To request a feature create a new issue with the label `feature-request`. Find a bug? Please add an issue with the label `bugs`. Pull requests are also welcomed but please add an issue on the requested feature first (unless it's a simple bug fix or readme change)
