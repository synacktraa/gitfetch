package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

//get color code according to repo language
func (repo *Repository) getColorCode(lang string) uint8 {

	code := map[string]uint8{}

	jsonFile, err := os.Open("code.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &code)

	if code[lang] == 0 {
		rand.Seed(time.Now().UnixNano())
		return uint8(256 + rand.Uint32())
	}
	return code[lang]
}

// Branches struct to store the branch name
type Branches []struct {
	Name string `json:"name"`
}

// get all available branches for specified repository
func (storage *Branches) getBranch(url string) []string {
	branchUrl := fmt.Sprintf("%s/branches", url)
	branchJson := Request(branchUrl)

	branches := func() []string {
		json.Unmarshal(branchJson, &storage)

		var branchVector []string
		for _, val := range *storage {
			branchVector = append(branchVector, val.Name)
		}
		return branchVector
	}()
	return branches
}

// get total commits of all the branches
func (repo *Repository) branchCommits(url string) map[string]int {

	br := new(Branches)
	branches := br.getBranch(url)
	commitMap := make(map[string]int, len(branches))

	for _, branch := range branches {

		CommitUrl := fmt.Sprintf("%s/commits?sha=%s&per_page=1&page=1", url, branch)

		res, err := http.Get(CommitUrl)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		defer res.Body.Close()

		func() {
			link := res.Header.Get("Link")
			re := regexp.MustCompile("page=([0-9]{1,})>; rel=\"last\"")

			cnt, err := strconv.Atoi(re.FindStringSubmatch(link)[1])
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			commitMap[branch] = cnt
		}()
	}
	return commitMap

}
