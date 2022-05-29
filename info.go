package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gookit/color"
)

// Final user info struct
type UserInstance struct {
	Username     string
	Company      string
	Blog         string
	Location     string
	Email        string
	Bio          string
	Repositories int
	Gists        int
	Followers    int
	Following    int
	Created      string
}

// Final repo info struct
type RepoInstance struct {
	Name          string
	Repository    string
	Description   string
	Created       string
	Modified      string
	DefaultBranch string
	Stars         int
	Forks         int
	Language      string
	License       struct {
		Name string
	}
}

// Request function checks status code and evalualtes accordingly
func Request(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
	case 404:
		fmt.Fprintln(os.Stderr, "RequestError: repository/user not found")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "RequestError: returned status code => %v ", res.Status)
		os.Exit(0)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return body
}

// global variables
var (
	setcolor uint8
	text     = color.C256(231)
	wg       sync.WaitGroup
)

// Print user stats
func (user *User) Print(host HostCode, url string) {

	instance := new(UserInstance)
	rand.Seed(time.Now().UnixNano())
	setcolor = uint8(rand.Intn(230-20) + 20)
	var host_ string

	switch host {
	case github:
		instance = (*UserInstance)(&user.Github)
		host_ = "github.com"
	case gitlab:
		instance = (*UserInstance)(&user.Gitlab)
		host_ = "gitlab.com"
	}
	keycolor := color.C256(setcolor)

	keycolor.Printf("\n\t%v", host_)
	fmt.Printf("/")
	keycolor.Printf("%v\n\n", instance.Username)
	instance.Username = ""

	if instance.Bio != "" {
		keycolor.Printf("%v\n\n", instance.Bio)
	}
	instance.Bio = ""

	refVal := reflect.ValueOf(*instance)
	typeOfRefVal := refVal.Type()

	for i := 0; i < refVal.NumField(); i++ {
		key := typeOfRefVal.Field(i).Name
		value := refVal.Field(i).Interface()

		if value != "" {

			func() {
				fmt.Printf("=> ")
				keycolor.Printf(strings.ToLower(key))
				fmt.Printf(": ")
			}()

			if key == "Created" {
				text.Printf("%v\n", strings.Split(value.(string), "T")[0])
				continue
			}
			text.Printf("%v\n", value)

		}
	}
}

// Print repo stats
func (repo *Repository) Print(host HostCode, url string) {

	var commits map[string]int
	instance := new(RepoInstance)

	switch host {
	case github:
		instance = (*RepoInstance)(&repo.Github)
		setcolor = repo.getColorCode(strings.ToLower(repo.Github.Language))
	case gitlab:
		instance = (*RepoInstance)(&repo.Gitlab)
		rand.Seed(time.Now().UnixNano())
		setcolor = uint8(rand.Intn(230-20) + 20)
	}

	keycolor := color.C256(setcolor)

	printInfo := func() {
		header := func() {
			spl := strings.Split(instance.Name, "/")
			keycolor.Printf("\n\t%v", spl[0])
			fmt.Printf("/")
			keycolor.Printf("%v\n\n", spl[1])
			text.Printf("%v\n\n", instance.Description)

			instance.Name = ""
			instance.Description = ""
		}

		header()
		refVal := reflect.ValueOf(*instance)
		typeOfRefVal := refVal.Type()

		for i := 0; i < refVal.NumField(); i++ {

			key := typeOfRefVal.Field(i).Name
			value := refVal.Field(i).Interface()

			if value != "" {

				func() {
					fmt.Printf("=> ")
					keycolor.Printf(strings.ToLower(key))
					fmt.Printf(": ")
				}()

				if key == "License" {
					text.Printf("%v\n", value.(struct{ Name string }).Name)
					continue
				}
				if key == "Created" || key == "Modified" {
					text.Printf("%v\n", strings.Split(value.(string), "T")[0])
					continue
				}

				text.Printf("%v\n", value)
			}

		}
	}
	wg.Add(1)
	go func() {
		commits = repo.branchCommits(url)
		wg.Done()
	}()
	go printInfo()
	wg.Wait()
	if hostcode == github { // supports only github for now
		dl := func() {
			fmt.Printf("=> ")
			keycolor.Printf("commits")
			fmt.Printf(": {\n")
			for k, v := range commits {
				text.Printf("\t%v: ", k)
				keycolor.Println(v)
			}
			fmt.Println("   }")
		}
		defer dl()
	}
}
