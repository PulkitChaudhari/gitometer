package main

import (
	"flag"
	"fmt"
	"os"
	// "strings"
)

func stats(email *string) {
	fmt.Println("Statistics for ", *email, " are as follows : ")
}

func recursiveRepositoriesScan(path *string) {
	entries, err := os.ReadDir(*path)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _,entry := range entries {

		fmt.Println(entry, entry.IsDir())

		if !entry.IsDir() && entry.Name() == ".git" {
			fmt.Println("Found git project at ", *path)
		}

		if entry.IsDir() {
			fmt.Println("Found subdirectory : ", entry.Name())
			// var builder strings.Builder
			// builder.WriteString(*path)
			// builder.WriteString("/")
			// builder.WriteString(entry.Name())
			// subPath := builder.String()
			// recursiveRepositoriesScan(&subPath)
		}
	}
}

func getDotFilePath() {
	fmt.Println("Getting knowledge base file")
}

func addNewRepositoriesToFile() {
	fmt.Println("Adding newly scanned repositories")
}

func scan(path *string) {
	fmt.Println("Scanning ", *path, " for git repositories")
	recursiveRepositoriesScan(path)
	getDotFilePath()
	addNewRepositoriesToFile()
	fmt.Println("Successfully added add repositories under ", *path, "to knowledge base.")
}

func main() {
	var statsFlag = flag.String("stats","", "get git statistics")
	var addFlag = flag.String("add","", "add folder for Gitometer statistics")
	flag.Parse()
	if *addFlag != "" {
		scan(addFlag)
		return
	}
	stats(statsFlag)
}
