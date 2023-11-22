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

func TestKobeClient_CancelTask(t *testing.T) {
	client := NewKobeClient(host, port)
	result, err := client.RunAdhoc("master", "shell", "sleep 10", inventory)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)
	go func() {
		err = client.CancelTask(result.Id)
		if err != nil {
			t.Fatal(err)
			return
		}
	}()
	err = client.WatchRun(result.Id, os.Stdout)
	if err != nil {
		t.Log(err)
	}
	result, err = client.GetResult(result.Id)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(result)
}

func TestKobeClient_RunPlaybook(t *testing.T) {
	client := NewKobeClient(host, port)
	result, err := client.RunPlaybook("ko", "01-base.yaml", "", inventory)
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
func TestKobeClient_GetResult(t *testing.T) {
	client := NewKobeClient(host, port)
	result, err := client.GetResult("cc57c051-092c-4db5-bb38-afd189a319d6")
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(result)
}

func TestKobeClient_CreateProject(t *testing.T) {
	client := NewKobeClient(host, port)
	project, err := client.CreateProject("ko", "https://gitee.com/faireal/ansible.git")
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(project)
}

func TestKobeClient_ListProject(t *testing.T) {
	client := NewKobeClient(host, port)
	project, err := client.ListProject()
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(project)
}
