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

[English](README.md) | 中文

**Cassowary** 是一个现代化的 HTTP/S 负载测试工具，采用 Go 语言开发，设计直观且跨平台，专为开发者、测试人员和系统管理员打造。Cassowary 受到 k6、ab 和 httpstat 等优秀项目的启发。

---

目录
----

- [功能特性](#功能特性)
- [安装](#安装)
- [运行 Cassowary](#运行-cassowary)
  * [常规负载测试](#常规负载测试)
  * [文件读取模式](#文件读取模式)
  * [导出指标到文件](#导出指标到文件)
  * [导出指标到 Prometheus](#导出指标到-prometheus)
  * [导出指标到 Cloudwatch](#导出指标到-cloudwatch)
  * [直方图](#直方图)
  * [箱线图](#箱线图)
  * [POST 数据负载测试](#post-数据负载测试)
  * [指定测试持续时间](#指定测试持续时间)
  * [添加 HTTP 头](#添加-http-头)
  * [禁用 HTTP Keep-Alive](#禁用-http-keep-alive)
  * [x509 认证](#x509-认证)
  * [分布式负载测试](#分布式负载测试)
- [将 Cassowary 导入为模块](#将-cassowary-导入为模块)
- [版本控制](#版本控制)
- [贡献](#贡献)

功能特性  
--------

📌 **两种负载测试模式**：标准模式和扩展模式，扩展模式支持从文件中指定 URL 路径（适合测试多个底层微服务）  
📌 **CI 友好**：非常适合集成到 CI 流水线中  
📌 **灵活的指标输出**：支持 Cloudwatch 指标、Prometheus 指标（推送至 Prometheus PushGateway）以及 JSON 文件  
📌 **高度可配置**：允许传入任意 HTTP 头，可自定义 HTTP 客户端配置  
📌 **支持多种 HTTP 方法**：支持 GET、POST、PUT 和 PATCH，POST、PUT 和 PATCH 数据可通过文件定义  
📌 **跨平台**：提供适用于 Linux、Mac OSX 和 Windows 的单一预编译二进制文件  
📌 **可导入**：除了命令行工具外，Cassowary 还可作为模块导入到 Go 应用中  
📌 **可视化支持**：Cassowary 可将请求数据导出为直方图和箱线图（PNG 格式）  

<img src="https://imgur.com/ac8F8eD.gif" />

安装  
--------

从 [GitHub Releases 页面](https://github.com/rogerwelin/cassowary/releases) 下载预编译的二进制文件。你可以选择将 **cassowary** 二进制文件放入你的 `PATH` 中，以便在任意位置运行。或者，你也可以：

### Homebrew（Mac OSX）  
在 Mac 上使用 Homebrew 包管理器安装 **cassowary**：

```bash
$ brew update && brew install cassowary
```

### Docker  

通过官方 Docker 镜像直接运行 **cassowary**：

```bash
$ docker run rogerw/cassowary:v0.14.1 -u http://www.example.com -c 1 -n 10
```

本地开发：

```bash
$ GOOS=linux go build -o dist/docker/cassowary cmd/cassowary/*.go
$ docker build -f dist/docker/Dockerfile -t test_cassowary dist/docker
$ docker run test_cassowary -u http://www.example.com -c 1 -n 10
```

若运行 `docker run` 时不带参数，将打印帮助信息。

### ArchLinux/Manjaro

从 [AUR](https://aur.archlinux.org/packages/cassowary-git) 安装 Cassowary 的开发版本：

```bash
yay -S cassowary-git
```

或手动构建和安装：

```bash
git clone https://aur.archlinux.org/cassowary-git.git
cd cassowary-git
makepkg -si
```

### Nix/NixOS

通过 [Nix](https://nixos.org) 包管理器安装 Cassowary：

```
nix-env -iA cassowary
```

### CentOS/RHEL (RPM)

若需自行构建 RPM 包，可使用 [cassowary.spec](https://github.com/rogerwelin/cassowary/blob/master/dist/rpm/cassowary.spec) 文件。

运行 Cassowary  
--------

### 常规负载测试  
示例：对 www.example.com 运行 **cassowary**，使用 10 个并发用户执行 100 个请求：

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

### 文件读取模式  
示例：在文件读取模式下运行 **cassowary**，所有 URL 路径从外部文件指定（也可通过 HTTP 获取）。默认情况下，若未指定 `-n` 标志，Cassowary 将为文件中每个路径发起一次请求。使用 `-n` 标志可指定对这些 URL 路径的总请求数。示例：

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

### 导出指标到文件  
示例：将 **cassowary** 的 JSON 指标导出到文件：

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

> 如果未指定 `json-metrics-file` 标志，默认文件名将为 `out.json`。

若需所有请求的原始数据（CSV 格式），可使用 `--raw-output` 标志导出：

> 输出文件名将为 `raw.csv`。

### 导出指标到 Prometheus  
示例：通过指定 Prometheus PushGateway 的 URL，将 **cassowary** 指标导出到 Prometheus：

```bash
$ ./cassowary run -u http://localhost:8000 -c 125 -n 100000 -p http://pushgatway:9091

Starting Load Test with 100000 requests using 125 concurrent users

[ omitted for brevity ]

```

### 导出指标到 Cloudwatch  
**Cassowary** 可通过添加不带值的 `--cloudwatch` 标志将指标导出到 AWS Cloudwatch。请注意，你需要指定使用的 AWS 区域，最简单的方法是通过环境变量：

```bash
$ export AWS_REGION=eu-north-1 && ./cassowary run -u http://localhost:8000 -c 125 -n 100000 --cloudwatch

Starting Load Test with 100000 requests using 125 concurrent users

[ omitted for brevity ]

```

### 直方图  
通过添加不带值的 `--histogram` 标志，Cassowary 将生成请求总持续时间的直方图（PNG 格式，保存为当前目录下的 `hist.png`）。示例：

<img src="https://i.imgur.com/VLEsVOY.png" width="300" height="300" />

### 箱线图  
通过添加不带值的 `--boxplot` 标志，Cassowary 将生成请求总持续时间的箱线图（PNG 格式，保存为当前目录下的 `boxplot.png`）。

### POST 数据负载测试  
示例：对 POST 端点发起请求，POST 的 JSON 数据从文件中定义：

```bash
$ ./cassowary run -u http://localhost:8000/add-user -c 10 -n 1000 --postfile user.json

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

### PATCH 数据负载测试  
示例：对 PATCH 端点发起请求，PATCH 的 JSON 数据从文件中定义：

```bash
$ ./cassowary run -u http://localhost:8000/add-user -c 5 -n 200 --patchfile user.json

Starting Load Test with 200 requests using 5 concurrent users

[ omitted for brevity ]

```

### 指定测试持续时间  
示例：为负载测试指定持续时间，以下命令将在 30 秒内发送 100 个请求：

```bash
$ ./cassowary run -u http://localhost:8000 -n 100 -d 30

Starting Load Test with 100 requests using 1 concurrent users

[ omitted for brevity ]

```

### 添加 HTTP 头  
示例：在运行 **cassowary** 时添加 HTTP 头：

```bash
$ ./cassowary run -u http://localhost:8000 -c 10 -n 1000 -H 'Host: www.example.com'

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

### 禁用 HTTP Keep-Alive  
示例：禁用 HTTP Keep-Alive（默认启用）：

```bash
$ ./cassowary run -u http://localhost:8000 -c 10 -n 1000 --disable-keep-alive

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

### 指定 CA 证书  
示例：指定 CA 证书：

```bash
$ ./cassowary run -u http://localhost:8000 -c 10 -n 1000 --ca /path/to/ca.pem

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

### x509 认证  
示例：为 mTLS 指定客户端认证：

```bash
$ ./cassowary run -u https://localhost:8443 -c 10 -n 1000 --cert /path/to/client.pem --key /path/to/client-key.pem --ca /path/to/ca.pem

Starting Load Test with 1000 requests using 10 concurrent users

[ omitted for brevity ]

```

### 分布式负载测试  
若需在多台机器上扩展负载测试，可通过分布式方式运行 Cassowary。最简单的方法是使用 Kubernetes 集群。使用 batch 类型，并通过 `spec.parallelism` 键指定同时运行的 Cassowary 实例数量：

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

应用该 YAML 文件：

```bash
$ kubectl apply -f cassowary.yaml
```

将 Cassowary 导入为模块  
--------

Cassowary 可作为模块导入到你的 Go 应用中。首先通过 go mod 获取依赖：

```bash
$ go mod init test && go get github.com/rogerwelin/cassowary/pkg/client
```

以下示例展示如何从代码中触发负载测试并打印结果：

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

更多模块示例 [可在此查看](docs/LIBRARY.md)。

版本控制  
--------

Cassowary 遵循语义化版本控制。公共库（pkg/client）在达到稳定 v1.0.0 版本之前，可能会破坏向后兼容性。

贡献  
--------

欢迎贡献！如需请求新功能，请创建带有 `feature-request` 标签的问题。发现 bug？请创建带有 `bugs` 标签的问题。欢迎提交 Pull Request，但请先为请求的功能创建问题（除非是简单的 bug 修复或 README 修改）。