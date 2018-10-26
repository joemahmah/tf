package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"path/filepath"
)

//Handles CLI input.
//Finds and processes subcommands.
func ProcessCli() {

	tagCmd := flag.NewFlagSet("tag", flag.ExitOnError)
	var tagTags = tagCmd.String("t", "", "Tag(s) to be added.")

	untagCmd := flag.NewFlagSet("untag", flag.ExitOnError)
	var untagTags = untagCmd.String("t", "", "Tag(s) to be removed.")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	lstagCmd := flag.NewFlagSet("lstag", flag.ExitOnError)

	lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	/*mvCmd := flag.NewFlagSet("mv", flag.ExitOnError)
	cpCmd := flag.NewFlagSet("cp", flag.ExitOnError)

	mvtagCmd := flag.NewFlagSet("mvtag", flag.ExitOnError)
	cptagCmd := flag.NewFlagSet("cptag", flag.ExitOnError)
	rmtagCmd := flag.NewFlagSet("rmtag", flag.ExitOnError)
	mergetagCmd := flag.NewFlagSet("mergetag", flag.ExitOnError)
*/
	if len(os.Args) == 1 {
		fmt.Println("usage: tf <command>")
		return
	}

	switch os.Args[1] {
		case "tag":
			tagCmd.Parse(os.Args[2:])
			ProcessTagCmd(*tagTags, tagCmd.Args())
		case "untag":
			untagCmd.Parse(os.Args[2:])
			ProcessUntagCmd(*untagTags, untagCmd.Args())
		case "list":
			listCmd.Parse(os.Args[2:])
			ProcessListCmd(listCmd.Args())
		case "lstag":
			lstagCmd.Parse(os.Args[2:])
			ProcessLstagCmd(lstagCmd.Args())
		case "ls":
			lsCmd.Parse(os.Args[2:])
			ProcessLsCmd(lsCmd.Args())
	}

}

//var tags is the -t flag
func ProcessTagCmd(tags string, args []string) {
	//Check for errors
	if tags == "" && len(args) <= 1 { //no flag, not enough args
		fmt.Println("Need at least one tag an one file...")
		return
	} else if tags != "" && len(args) < 1 { //flag, not enough args
		fmt.Println("Need at least one tag an one file...")
		return
	}

	//if tag flag not used, assume first arg is flag
	if tags == "" {
		tags = args[0]
		args = args[1:]
	}

	//Slit tags into actual tags
	taglist := strings.Fields(tags)

	//handle each file
	for _, filepath := range args {
		//Check if file exists; if not ignore
		if _, err := os.Stat(filepath); os.IsNotExist(err){
			fmt.Println(filepath + " does not exist.")
			continue
		}

		AddFile(filepath)
		fid, err := GetFile(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}

		//add tags to file
		tidlist, err := GetOrCreateTags(taglist...)
		TagFile(fid, tidlist...)
	}

}

func ProcessUntagCmd(tags string, args []string) {
	//Check for errors
	if len(args) < 1 { //check if files passed
		fmt.Println("Need at least one file...")
		return
	}

	for _, filepath := range args {
		//Check if file exists; if not ignore
		if _, err := os.Stat(filepath); os.IsNotExist(err){
			fmt.Println(filepath + " does not exist.")
			continue
		}

		AddFile(filepath)
		fid, err := GetFile(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}

		if tags == "" { //if no tags entered, remove all tags
			UntagAllFile(fid)
		} else { //remove specified tags
			taglist := strings.Fields(tags)
			tidlist, _ := GetTags(taglist...)
			UntagFile(fid, tidlist...)
		}

	}
}

func ProcessListCmd(args []string) {
	for _, filepath := range args {
		//Check if file exists; if not ignore
		if _, err := os.Stat(filepath); os.IsNotExist(err){
			fmt.Println(filepath + ": file does not exist")
			continue
		}

		fid, err := GetFile(filepath)
		if err != nil {
			fmt.Println(filepath + ": file not tagged")
			continue
		}

		//Get tags
		tags, err := GetTagsForFile(fid)
		fmt.Println(filepath + ": ", tags)
	}
}

func ProcessLstagCmd(args []string) {
	for _, tag := range args {
		tid, err := GetTag(tag)
		if err != nil {
			fmt.Println(tag + ": tag does not exist")
			continue
		}

		//Get files
		files, err := GetFilesForTag(tid)
		fmt.Println(tag + ": ")
		for _, file := range files {
			fmt.Println("\t" + file)
		}
	}
}

func ProcessLsCmd(args []string) {
	//If no directory is given, use current one
	if len(args) <= 0 {
		args = append(args, ".")
	}

	for _, dir := range args {
		//Check if exists
		if _, err := os.Stat(dir); os.IsNotExist(err){
			fmt.Println(dir + ": directory does not exist")
			continue
		}

		//Checking if directory
		file, err := os.Open(dir)
		defer file.Close()
		if err != nil {
			fmt.Println(dir + ": unable to open")
			continue
		}

		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Println(dir + ": error reading")
			continue
		}

		//actual check
		if !fileInfo.IsDir() {
			fmt.Println(dir + ": not a directory")
			continue
		}


		//perform ls for all files
		fmt.Println("Directory " + dir + ":\n---------")
		filesInfo, err := file.Readdir(-1)
		filePaths := make([]string, 0)
		for _, subfileInfo := range filesInfo {
			filePaths = append(filePaths, filepath.Join(dir,subfileInfo.Name()))
		}
		//pass the file list over to list function
		ProcessListCmd(filePaths)
		fmt.Println("")
	}
}
