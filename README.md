# httpd 

![](https://github.com/harley9293/httpd/workflows/Build/badge.svg)
[![codecov](https://codecov.io/gh/harley9293/httpd/graph/badge.svg?token=5Vfc2RLAC6)](https://codecov.io/gh/harley9293/httpd)

A simple and easy-to-use http server

## Installation

```shell
go get -u github.com/harley9293/httpd
```

## Usage

```go
import "github.com/harley9293/httpd"

func main() {
	service := httpd.NewService(&httpd.Config{})
	service.AddMiddleWare(func1, func2...)
	service.AddHandler(httpd.GET, "/", funcTest, MiddleWare1, MiddleWare2...)
	service.LinstenAndServe(":8080")
}
```


