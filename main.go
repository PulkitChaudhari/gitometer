package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"log"
	"os/User"
	"bufio"
)

func recursiveRepositoriesScan(path string) []string {

	entries, err := os.ReadDir(path)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	directories := []string{}

	for _,entry := range entries {

		if entry.IsDir() {

			if entry.Name() == ".git" {
				fmt.Println("Found git project at ", path)
				directories = append(directories,path)
			} else {
				var builder strings.Builder
				builder.WriteString(path)
				builder.WriteString("/")
				builder.WriteString(entry.Name())
				subPath := builder.String()
				subPathDirectories := recursiveRepositoriesScan(subPath)
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

func sliceContains(parsedDotFileSlice []string, directory string) bool {
	for _, dir := range parsedDotFileSlice {
		if dir == directory {
			return true
		}
	}
	return false
}

func joinDirectories(parsedDotFileSlice []string, gitDirectories []string) []string {

	for _,directory := range gitDirectories {
		if !sliceContains(parsedDotFileSlice,directory) {
			parsedDotFileSlice = append(parsedDotFileSlice,directory)
		}
	}
	return parsedDotFileSlice
}

func writeToDotFile(parsedDotFileSlice []string, dotFilePath string, gitDirectories []string) {
	file, err := os.OpenFile(dotFilePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file : ", err)
		return
	}
	fmt.Println("File opened successfully")

	parsedDotFileSlice = joinDirectories(parsedDotFileSlice,gitDirectories)

	for _,gitDirectory := range parsedDotFileSlice {
		_, err = file.WriteString(gitDirectory + "\n")
		if err != nil {
			fmt.Println("Error writing to file : ", err)
			return
		}
	}
	file.Close()
}

func parseDotFileToSlice(dotFilePath string) []string {
	file, err := os.Open(dotFilePath)
	if err != nil {
		fmt.Println("Error opening file : ", err)
	}
	scanner := bufio.NewScanner(file)
	parsedDotFile := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		parsedDotFile = append(parsedDotFile,line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file : ", err)
	}
	return parsedDotFile
}

func getDotFilePath() string {
	const dotFileName = ".gogitlocalstats"
	usrHomeDir := getUsrHomeDir()
	dotFilePath := usrHomeDir + "/" + dotFileName
	return dotFilePath
}

func addNewRepositoriesToDotFile(gitDirectories []string) {
	dotFilePath := getDotFilePath()
	createDotFileIfDoesNotExist(dotFilePath)
	parsedDotFileSlice := parseDotFileToSlice(dotFilePath)
	writeToDotFile(parsedDotFileSlice, dotFilePath, gitDirectories)
}

func scan(path string) {
	fmt.Println("Scanning ", path, " for git repositories")
	gitDirectories := recursiveRepositoriesScan(path)
	addNewRepositoriesToDotFile(gitDirectories)
	fmt.Println("Successfully added add repositories under ", path, "to knowledge base.")
}

func fillCommits(email string, path string, commitMap map[int]int) map[int]int {
	return make(map[int]int, 30)
}

func processRepositories(email string) map[int]int {
	dotFilePath := getDotFilePath()
	gitRepoDirectories := parseDotFileToSlice(dotFilePath)
	daysInMap := 30
	commitMap := make(map[int]int,daysInMap)
	for i := 0; i < daysInMap; i++ {
		commitMap[i] = 0
	}
	for _, path := range gitRepoDirectories {
		commitMap = fillCommits(email,path,commitMap)
	}
	return commitMap
}

func printCommitStats(commits []string) {
	fmt.Println("printCommitStats called")
}

func stats(email string) {
	commitMap := processRepositories(email)
	for i := 0; i < 30; i++ {
		fmt.Println(commitMap[i])
	}
}

func main() {
	var statsFlag = flag.String("stats","", "get git statistics")
	var addFlag = flag.String("add","", "add folder for Gitometer statistics")
	flag.Parse()
	if *addFlag != "" {
		scan(*addFlag)
		return
	}
	stats(*statsFlag)
}
