package main

import (
	"fmt"
	"os"
	"strings"
)

func displayVersion() {
	text.Println("gitfetch v1.0.0")
}

func main() {

	if len(os.Args) != 2 || os.Args[1] == "help" {

		displayVersion()
		text.Println("\nUsage: gitfetch host/username\n", "      gitfetch host/username/repository")
		text.Println("\nOptions:\n    help  displays the help section")
		text.Println("    version  displays the gitfetch version")
		text.Println("\n    hosts:-")
		text.Println("\tgithub -> gh, github, github.com")
		text.Println("\tgitlab -> gl, gitlab, gitlab.com")
		text.Println("\tgitea -> gt, gitea, gitea.com")
		os.Exit(0)

	} else if os.Args[1] == "version" {
		displayVersion()
		os.Exit(0)
	}

	var (
		host_ string
		arg   []string
	)
	url := func() string {

		host := new(Host)
		host.MapHosts()

		vector := strings.Split(os.Args[1], "/")
		switch len(vector) {
		case 1, 2:
			if host.Hosts[vector[0]] == 0 {
				host_ = "github"
				hostcode = github
				arg = vector
			} else {
				host_ = vector[0]
				arg = []string{vector[1]}
				hostcode = host.Hosts[host_]
			}
		case 3:
			if host.Hosts[vector[0]] == 0 {
				fmt.Fprintln(os.Stderr, "Error: unable to resolve host")
				os.Exit(1)
			} else {
				host_ = vector[0]
				arg = []string{vector[1], vector[2]}
				hostcode = host.Hosts[host_]
			}
		}
		return host.generateUrl(host_, arg)
	}()

	switch choice {
	case "user":
		user := new(User)
		user.Request(url)
		user.Print(hostcode, url)
	case "repo":
		repo := new(Repository)
		repo.Request(url)

		repo.Print(hostcode, url)
	}
}
