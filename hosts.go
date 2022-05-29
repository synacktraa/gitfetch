package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Hostcode
type HostCode int

const (
	github HostCode = iota + 1
	gitlab
	gitea
)

// Host struct which will mapped accordingly
type Host struct {
	Hosts map[string]HostCode
}

// MapHosts maps hosts to their respective hostcode
func (host *Host) MapHosts() {
	host.Hosts = map[string]HostCode{
		"github.com": github,
		"gitlab.com": gitlab,
		"gitea.com":  gitea,
		"github":     github,
		"gitlab":     gitlab,
		"gitea":      gitea,
		"gh":         github,
		"gl":         gitlab,
		"gt":         gitea,
	}
}

var (
	hostcode HostCode
	choice   string
)

// Returns url according to host and host info
func (host *Host) generateURL(site string, arg []string) string {

	var url string
	if len(arg) == 1 {
		switch host.Hosts[site] {

		case github:
			url = fmt.Sprintf("https://api.github.com/users/%s", arg[0])
		case gitlab, gitea:
			// url = "gitlab.com/api/users"
			fmt.Fprintln(os.Stderr, "SupportError: gitfetch doesn't support gitlab and codeberg for now")
			os.Exit(0)
		}
		choice = "user"
	} else {
		switch host.Hosts[site] {
		case github:
			url = fmt.Sprintf("https://api.github.com/repos/%s/%s", arg[0], arg[1])
		case gitlab:
			url = fmt.Sprintf("https://gitlab.com/api/v4/projects/%s%s%s?license=1", arg[0], "%2F", arg[1])
		case gitea:
			fmt.Fprintln(os.Stderr, "SupportError: gitfetch doesn't support codeberg for now")
			os.Exit(0)
		}
		choice = "repo"
	}
	return url
}

// User struct to store the user json values
type User struct {
	Github struct {
		Username     string `json:"login"`
		Company      string
		Blog         string
		Location     string
		Email        string
		Bio          string
		Repositories int `json:"public_repos"`
		Gists        int `json:"public_gists"`
		Followers    int
		Following    int
		Created      string `json:"created_at"`
	}
	Gitlab struct {
		Username     string `json:"login"`
		Company      string
		Blog         string
		Location     string
		Email        string
		Bio          string
		Repositories int `json:"public_repos"`
		Gists        int `json:"public_gists"`
		Followers    int
		Following    int
		Created      string `json:"created_at"`
	}
	Gitea struct{}
}

// Repository struct to store the repository json values
type Repository struct {
	Github struct {
		Name          string `json:"full_name"`
		Repository    string `json:"html_url"`
		Description   string `json:"description"`
		Created       string `json:"created_at"`
		Modified      string `json:"updated_at"`
		DefaultBranch string `json:"default_branch"`
		Stars         int    `json:"stargazers_count"`
		Forks         int    `json:"forks_count"`
		Language      string `json:"language"`
		License       struct {
			Name string `json:"name"`
		} `json:"license"`
	}
	Gitlab struct {
		Name          string `json:"path_with_namespace"`
		Repository    string `json:"web_url"`
		Description   string `json:"description"`
		Created       string `json:"created_at"`
		Modified      string `json:"last_activity_at"`
		DefaultBranch string `json:"default_branch"`
		Stars         int    `json:"star_count"`
		Forks         int    `json:"forks_count"`
		Language      string `json:"language"`
		License       struct {
			Name string `json:"name"`
		} `json:"license"`
	}
	Gitea struct{}
}

// Request user info
func (user *User) Request(url string) {
	jsonDump := Request(url)
	switch hostcode {
	case github:
		json.Unmarshal(jsonDump, &user.Github)
	case gitlab:
		json.Unmarshal(jsonDump, &user.Gitlab)
	case gitea:
		json.Unmarshal(jsonDump, &user.Gitea)
	}
}

// Request repository info
func (repo *Repository) Request(url string) {
	jsonDump := Request(url)
	switch hostcode {
	case github:
		json.Unmarshal(jsonDump, &repo.Github)
	case gitlab:
		json.Unmarshal(jsonDump, &repo.Gitlab)
	case gitea:
		json.Unmarshal(jsonDump, &repo.Gitea)
	}
}
