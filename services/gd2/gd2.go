package gd2

import (
	"github.com/gluster/glusterd2/pkg/restclient"
	"os"
)

var (
	Client *restclient.Client
)

func init() {
	Client = restclient.New(os.Getenv("GD2_ENDPOINT"), "", "", "", true)
}
