package util

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
)

func CloneRepository(url, dest string) error {
	options := git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}
	_, err := git.PlainClone(dest, false, &options)
	return err
}

func CloneRepositoryWithAuth(url, dest string, username, password string) error {
	options := git.CloneOptions{
		URL:      url,
		Auth:     &http.BasicAuth{Username: username, Password: password},
		Progress: os.Stdout,
	}
	_, err := git.PlainClone(dest, false, &options)
	return err
}
