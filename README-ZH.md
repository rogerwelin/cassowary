<p align="center"><a href="https://github.com/rogerwelin/cassowary"><img src="cass-logo.png" alt="cassowary"></a></p>
<p align="center">
  <a href="https://goreportcard.com/badge/github.com/rogerwelin/cassowary"><img src="https://goreportcard.com/badge/github.com/rogerwelin/cassowary" alt="Go Report Card"></a>
  <a href="https://travis-ci.org/rogerwelin/cassowary"><img src="https://travis-ci.org/rogerwelin/cassowary.svg?branch=master" alt="Build status"></a>
  <a href="https://github.com/rogerwelin/cassowary/blob/master/LICENSE"><img src="https://img.shields.io/github/license/rogerwelin/cassowary" alt="License"></a>
  <a href="https://github.com/rogerwelin/cassowary/blob/master/go.mod"><img src="https://img.shields.io/github/go-mod/go-version/rogerwelin/cassowary" alt="Go version"></a>
  <a href="https://github.com/rogerwelin/cassowary/releases"><img src="https://img.shields.io/github/v/release/rogerwelin/cassowary.svg" alt="Current Release"></a>
  <a href="https://godoc.org/github.com/rogerwelin/cassowary"><img src="https://godoc.org/github.com/rogerwelin/cassowary?status.svg" alt="godoc"></a>
  <a href="https://gocover.io/github.com/rogerwelin/cassowary/pkg/client"><img src="https://gocover.io/_badge/github.com/rogerwelin/cassowary/pkg/client" alt="Coverage"></a>
</p>

**Cassowary** 是一个最新的跨平台HTTP/S负载测试工。我使用了Go编辑Cassowary，希望Developer, tester 以及sysadmins都可以便捷的进行负载测试。Cassowary受到了很多经典的开源项目的启发，比如k6, ab和httestat。

---


目录
----
- [功能](#功能)
- [安装](#安装)
- [使用](#使用)
- [反馈](#反馈)


功能
--------
- 两种测试模式:标准和自定义。在自定义模式下可以选择URL路径
- CI友好
- 灵活的算法: 可以向Prometheus PushGateway直接发送算法，也可以以JSON文件的形式发送算法
- 灵活调节:可以自由选择使用哪种HTTP头字段
- 跨平台: 一个二进制文件可同时支持Linux，Mac OSX和Windows

<img src="https://i.imgur.com/geJykYH.gif" />


安装
--------

从GitHub Releases page下载二进制文件。可以选择把Cassowary二进制文件放在PATH里，这样在任何页面写都可以运行Cassowary。


### Nix/NixOS

Cassowary可以安装在Nix OS上。


使用
--------

示例:10个用户同时向www.example.com 发送100个访问

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

示例:访问外部文件指定的URL路径(外部文件也可以是http路径的)

```bash
$ ./cassowary run-file -u http://localhost:8000 -c 10 -f urlpath.txt

Starting Load Test with 3925 requests using 10 concurrent users

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

示例:导出Cassowary的Json算法

```bash
$ ./cassowary run --json-metrics --json-metrics-file=metrics.json -u http://localhost:8000 -c 125 -n 100000

Starting Load Test with 100000 requests using 125 concurrent users

[ omitted for brevity ]

```

> 如果没有指定Json算法的导出文件名，系统会使用默认文件名out.json.


示例:指定一个Prometheus Pushgateway URL，把Cassowary的Json算法导出到Prometheus

```bash
$ ./cassowary run -u http://localhost:8000 -c 125 -n 100000 -p http://pushgatway:9091

Starting Load Test with 100000 requests using 125 concurrent users

[ omitted for brevity ]

```


示例:添加HTTP头字段

```bash
$ ./cassowary run -u http://localhost:8000 -c 10 -n 1000 -H 'Host: www.example.com'

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```


示例:关闭http-keep-alive(http-keep-alive在默认下是激活的)

```bash
$ ./cassowary run -u http://localhost:8000 -c 10 -n 1000 --disable-keep-alive

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

示例:指定自定义 ca 证书
```bash
$ ./cassowary run -u http://localhost:8000 -c 10 -n 1000 --ca /path/to/ca.pem

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

示例:指定客户端证书信息
```bash
$ ./cassowary run -u http://localhost:8000 -c 10 -n 1000 --cert /path/to/client.pem --key /path/to/client-key.pem

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```


**以模块或者library导入Cassowary**  

Cassowary可以以模块的形式倒入/使用在你的Go程序。我们从使用go mod下载依赖关系开始

```bash
$ go mod init test && go get github.com/rogerwelin/cassowary/pkg/client
```

以下是一个简单示例:如何激活一个load test并且显示结果

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


反馈
--------

非常欢迎收到各种意见和建议！如果你觉得有某项功能可以更加完善，可以在Issue下发帖，最好能使用feature-request标签。如果你找到一个bug,一定要在issue下发帖告诉我，使用bugs标签。也非常欢迎Pull requests，一样的在issue下发帖，使用feature-request标签就好。

