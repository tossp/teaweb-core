package agents

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestSocketConnectivitySource_Execute(t *testing.T) {
	source := NewSocketConnectivitySource()
	source.Address = "127.0.0.1:27018"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}

func TestSocketConnectivitySource_Execute_2(t *testing.T) {
	source := NewSocketConnectivitySource()
	source.Address = "127.0.0.1:27017"
	source.Network = "tcp"
	value, err := source.Execute(nil)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(value, t)
}
