package client

import (
	"fmt"
	"github.com/KubeOperator/kobe/api"
	"os"
	"testing"
)

var (
	host      = "10.10.10.88"
	port      = 8080
	inventory = &api.Inventory{
		Hosts: []*api.Host{
			{
				Ip:       "10.10.10.88",
				Name:     "test",
				Port:     22,
				User:     "root",
				Password: "Trusfort@20151010",
				Vars:     map[string]string{},
			},
		},
		Groups: []*api.Group{
			{
				Name:     "master",
				Children: []string{},
				Vars:     map[string]string{},
				Hosts:    []string{"test"},
			},
		},
	}
)

func TestKobeClient_RunAdhoc(t *testing.T) {
	client := NewKobeClient(host, port)
	adhoc, err := client.RunAdhoc("all", "ping", "", inventory)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(adhoc)
}

func TestKobeClient_RunPlaybook(t *testing.T) {
	client := NewKobeClient(host, port)
	result, err := client.RunPlaybook("test", "test.yml", "abc", inventory)
	if err != nil {
		t.Fatal(err)
		return
	}
	err = client.WatchRun(result.Id, os.Stdout)
	if err != nil {
		t.Fatal(err)
		return
	}
}
