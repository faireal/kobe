package client

import (
	"fmt"
	"github.com/faireal/kobe/api"
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
				Password: "123456",
				Vars: map[string]string{
					"aaa": "3a",
				},
			},
		},
		Groups: []*api.Group{
			{
				Name:     "master",
				Children: []string{},
				Vars: map[string]string{
					"aaaa": "4a",
				},
				Hosts: []string{"test"},
			},
		},
	}
	inventory2 = &api.Inventory{
		Hosts: []*api.Host{
			{
				Ip:       "10.10.10.88",
				Name:     "10.10.10.88",
				Port:     22,
				User:     "root",
				Password: "123456",
				Vars: map[string]string{
					"aaa":  "3a",
					"aaaa": "4a",
				},
			},
			{
				Ip:       "10.10.10.108",
				Name:     "10.10.10.108",
				Port:     22,
				User:     "root",
				Password: "123456",
				Vars: map[string]string{
					"aaa":  "3a",
					"aaaa": "4a",
				},
			},
		},
	}
	projectName = "ko"
	url         = "https://gitee.com/faireal/ansible.git"
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
	result, err := client.RunPlaybook(projectName, "02-test.yaml", "", inventory2)
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
	project, err := client.CreateProject(projectName, url)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(project)
}

func TestKobeClient_CreateProjectWithAuth(t *testing.T) {
	client := NewKobeClient(host, port)
	project, err := client.CreateProjectWithAuth(projectName, url, "faireal", "faireal032.")
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

func TestKobeClient_DeleteProject(t *testing.T) {
	client := NewKobeClient(host, port)
	err := client.DeleteProject(projectName)
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestKobeClient_UpdateProject(t *testing.T) {
	client := NewKobeClient(host, port)
	err := client.DeleteProject(projectName)
	if err != nil {
		t.Fatal(err)
		return
	}
	project, err := client.CreateProjectWithAuth(projectName, url, "faireal", "faireal032.")
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(project)
}

func TestKobeClient_GetInvertory(t *testing.T) {
	client := NewKobeClient(host, port)
	id := "fbdc4d3a-d37e-42ec-b1ac-236faab8a63c"
	i, err := client.GetInvertory(id)
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Println(i)
}
