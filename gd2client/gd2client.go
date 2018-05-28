package gd2client

import (
	"github.com/ddliu/go-httpclient"
)

func BuildHello() string {
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "tendrl2-go",
		"Accept-Language":        "en-us",
	})

	res, err := httpclient.Get("http://google.com/search", map[string]string{
		"q": "news",
	})

	println(res.StatusCode, err)
	return "Hello, world."
}

func BuildHi() string {
	return "Hi, world."
}
