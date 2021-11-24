
## Using Cassowary as an Library

Cassowary can be imported and used as a module in your Go app. Start by fetching the dependency by using go mod:

```bash
$ go mod init test && go get github.com/rogerwelin/cassowary/pkg/client
```

**Example 1: Simple Load Test of an URL**  

```go
package main

import (
        "encoding/json"
	"fmt"

	cassowary "github.com/rogerwelin/cassowary/pkg/client"
)

func main() {
	cass := &cassowary.Cassowary{
		BaseURL:               "http://www.example.com",
		Groups: []QueryGroup{
			{
				Name:             "default",
				ConcurrencyLevel:      1,
				Requests:              10,
			},
		},
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

**Example 2: Load Test an URL across multiple URL paths**  

The following code will make 30 requests across the 3 URL paths declared in the URLPaths field:

```go
package main

import (
        "encoding/json"
	"fmt"

	cassowary "github.com/rogerwelin/cassowary/pkg/client"
)

func main() {
	cass := &cassowary.Cassowary{
		BaseURL:               "http://www.example.com",
		Groups: []QueryGroup{
			{
				ConcurrencyLevel:      2,
				Requests:              30,
				FileMode:	       true,
				URLPaths:	       []string{"/accounts", "/orders", "/customers"},
			},
		},
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

**Example 3: Custom TLS config**

```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"

	cassowary "github.com/rogerwelin/cassowary/pkg/client"
)

func main() {
	pemCerts, err := ioutil.ReadFile("testdata/ca.pem")
	if err != nil {
		panic("Invalid ca.pem path")
	}

	ca := x509.NewCertPool()
	if !ca.AppendCertsFromPEM(pemCerts) {
		panic("Failed to read CA from PEM")
	}

	cert, err := tls.LoadX509KeyPair("testdata/client.pem", "testdata/client-key.pem")
	if err != nil {
		panic("Invalid client.pem/client-key.pem path")
	}

	clientTLSConfig := &tls.Config{
		RootCAs:      ca,
		Certificates: []tls.Certificate{cert},
	}

	cass := &cassowary.Cassowary{
		BaseURL:               "http://www.example.com",
		Groups: []QueryGroup{
			{
				ConcurrencyLevel:      1,
				Requests:              10,
			},
		},
		TLSConfig:             clientTLSConfig,
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

**Example 4: Load Test an URL across multiple URL paths with iterator**  

The following code will make 30 requests across the 3 URL paths declared in the URLPaths field:

```go
package main

import (
        "encoding/json"
	"fmt"

	cassowary "github.com/rogerwelin/cassowary/pkg/client"
)

type URLIterator struct {
	pos  uint64
	data []string
	v    Validator
}

func (it *URLIterator) Next() *Query {
	for {
		pos := atomic.AddUint64(&it.pos, 1)
		if pos > uint64(len(it.data)) {
			if !atomic.CompareAndSwapUint64(&it.pos, pos, 0) {
				// retry
				continue
			}
			pos = 0
		} else {
			pos--
		}
		//return &Query{Method: "GET", URL: it.data[pos]}
		return &Query{Method: "POST", URL: it.data[pos], DataType: "application/json", Data: []byte("{ \"test\": \"POST\" }"), Validator: it.v}
	}
}

func NewURLIterator(data []string) *URLIterator {
	if len(data) == 0 {
		return nil
	}
	return &URLIterator{data: data, pos: 0}
}

func main() {
	it := cassowary.NewURLIterator([]string{"/test1", "/test2", "/test3"})

	cass := &cassowary.Cassowary{
		BaseURL:               "http://www.example.com",
		Groups: []QueryGroup{
			{
				Name:             "default",
				ConcurrencyLevel:      2,
				Requests:              30,
				FileMode:	       	   true,
				URLIterator:           it,
			},
		},
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