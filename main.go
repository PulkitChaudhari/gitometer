package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"log"
	"os/User"
)

func stats(email *string) {
	fmt.Println("Statistics for ", *email, " are as follows : ")
}

func recursiveRepositoriesScan(path *string) []string {

	entries, err := os.ReadDir(*path)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	directories := []string{}

	for _,entry := range entries {

		if entry.IsDir() {

			if entry.Name() == ".git" {
				fmt.Println("Found git project at ", *path)
				directories = append(directories,*path)
			} else {
				var builder strings.Builder
				builder.WriteString(*path)
				builder.WriteString("/")
				builder.WriteString(entry.Name())
				subPath := builder.String()
				subPathDirectories := recursiveRepositoriesScan(&subPath)
				directories = append(directories,subPathDirectories...)
			}
		}
	}

	return directories
}

func getUsrHomeDir() string {
	usr, err := user.Current()
    if err != nil {
        log.Fatal(err)
    }
    return usr.HomeDir
}

func createDotFileIfDoesNotExist(dotFilePath string) {
	file, err := os.OpenFile(dotFilePath,os.O_RDWR|os.O_EXCL|os.O_CREATE, 0644)

	if err != nil {
		if os.IsExist(err) {
			fmt.Println("File already exists, Skipping file creation")
			return
		}
		fmt.Println("Error creating file")
	}
	file.Close()
}

func writeToDotFile(dotFilePath string, gitDirectories []string) {
	file, err := os.OpenFile(dotFilePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file : ", err)
		return
	}
	fmt.Println("File opened successfully")

	_, err = file.WriteString("Hello World!")
	if err != nil {
		fmt.Println("Error writing to file : ", err)
		return
	}
	file.Close()
}

func addNewRepositoriesToDotFile(gitDirectories []string) {
	const dotFileName = ".gogitlocalstats"
	usrHomeDir := getUsrHomeDir()
	dotFilePath := usrHomeDir + "/" + dotFileName
	createDotFileIfDoesNotExist(dotFilePath)
	writeToDotFile(dotFilePath, gitDirectories)
}

func scan(path *string) {
	fmt.Println("Scanning ", *path, " for git repositories")
	gitDirectories := recursiveRepositoriesScan(path)
	addNewRepositoriesToDotFile(gitDirectories)
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
