package gd2

import (
	"encoding/json"
	"fmt"
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

type versionDetail struct {
	GlusterdVersion string  `json:"glusterd-version"`
	ApiVersion      float64 `json:"api-version"`
}

func Version(endpoint string) map[string]interface{} {
	res, _ := httpclient.Get(endpoint + "/version")
	body, _ := res.ToString()
	fmt.Println(body)
	//var version versionDetail
	var version map[string]interface{}
	if err := json.Unmarshal([]byte(body), &version); err != nil {
		panic(err)
	}
	fmt.Println(version)
	return version
}
