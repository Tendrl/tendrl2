package gd2

import (
	"testing"
)

func TestGd2(t *testing.T) {
	peers, _ := Client.Peers()
	actual := len(peers)
	expected := 1
	if expected != actual {
		t.Errorf("Error occured while testing GD2 Peers: %d != %d", expected, actual)
	}
}
