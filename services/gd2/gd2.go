package gd2

import (
	"encoding/json"
	"github.com/ddliu/go-httpclient"
	"github.com/gluster/glusterd2/pkg/restclient"
)

func init() {
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "tendrl2-go",
		"Accept-Language":        "en-us",
	})
}

func Client(endpoint string) *restclient.Client {
	return restclient.New(endpoint, "", "", "", true)
}

func Version(endpoint string) map[string]interface{} {
	res, _ := httpclient.Get(endpoint + "/version")
	body, _ := res.ToString()
	var version map[string]interface{}
	if err := json.Unmarshal([]byte(body), &version); err != nil {
		panic(err)
	}
	return version
}
