package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func main() {
	os.Exit(buildReadme())
}

func buildReadme() int {
	oldReadme, err := ioutil.ReadFile("README.md")

	if err != nil {
		fmt.Println(err)
		return 1
	}

	blogEntries, err := blogEntries()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	data := strings.Join(blogEntries, "\n")

	readme := replaceEntries(string(oldReadme), "blog", data)

	err = ioutil.WriteFile("README.md", []byte(readme), 0664)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

func blogEntries() ([]string, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://blog.mzumi.com/index.xml")
	if err != nil {
		return nil, err
	}

	entries := make([]string, 5)

	for i, item := range feed.Items[:5] {
		layout := "Mon, 02 Jan 2006 15:04:05 -0700"
		published, err := time.Parse(layout, item.Published)

		if err != nil {
			return nil, err
		}

		str := fmt.Sprintf("* [%s](%s) - %s", item.Title, item.Link, published.Format("2006-01-02"))
		fmt.Println(str)

		entries[i] = str
	}

	return entries, nil
}

func replaceEntries(readme string, comment string, str string) string {
	commentPrefix := fmt.Sprintf("<!-- %s starts -->", comment)
	commentSuffix := fmt.Sprintf("<!-- %s ends -->", comment)

	rep := regexp.MustCompile(fmt.Sprintf(`(?s)%s.+%s`, commentPrefix, commentSuffix))
	return rep.ReplaceAllString(readme, fmt.Sprintf("%s\n%s\n%s", commentPrefix, str, commentSuffix))
}
