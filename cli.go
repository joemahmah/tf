package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

//Handles CLI input.
//Finds and processes subcommands.
func ProcessCli() {

	tagCmd := flag.NewFlagSet("tag", flag.ExitOnError)
	var tags = tagCmd.String("t", "", "Tag(s) to be added.")

	//untagCmd := flag.NewFlagSet("untag", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	/*lstagCmd := flag.NewFlagSet("lstag", flag.ExitOnError)

	lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	mvCmd := flag.NewFlagSet("mv", flag.ExitOnError)
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
			ProcessTagCmd(*tags, tagCmd.Args())
		case "list":
			listCmd.Parse(os.Args[2:])
			ProcessListCmd(listCmd.Args())
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

func ProcessListCmd(args []string) {
	for _, filepath := range args {
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
