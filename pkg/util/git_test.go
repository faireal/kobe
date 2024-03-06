package util

import (
	"os"
	"testing"
)

func TestCloneRepositoryWithAuth(t *testing.T) {
	username := ""
	password := ""
	url := "https://gitee.com/faireal/ansible.git"
	dest := "D:/trusfort/kobe/test"
	err := os.Mkdir(dest, 0755)
	defer os.RemoveAll(dest)
	if err != nil {
		t.Fatal(err)
	}
	err = CloneRepositoryWithAuth(url, dest, username, password)
	if err != nil {
		t.Fatal(err)
	}
}
