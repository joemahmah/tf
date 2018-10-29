package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

//Handles CLI input.
//Finds and processes subcommands.
func ProcessCli() {

	tagCmd := flag.NewFlagSet("tag", flag.ExitOnError)
	var tagTags = tagCmd.String("t", "", "Tag(s) to be added.")
	var tagRecursive = tagCmd.Bool("r", false, "Recursively add tags to files and directories.")

	untagCmd := flag.NewFlagSet("untag", flag.ExitOnError)
	var untagTags = untagCmd.String("t", "", "Tag(s) to be removed.")
	var untagRecursive = untagCmd.Bool("r", false, "Recursively remove tags from files and directories.")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	lstagCmd := flag.NewFlagSet("lstag", flag.ExitOnError)

	lsCmd := flag.NewFlagSet("ls", flag.ExitOnError)
	mvCmd := flag.NewFlagSet("mv", flag.ExitOnError)
	cpCmd := flag.NewFlagSet("cp", flag.ExitOnError)
	rmCmd := flag.NewFlagSet("rm", flag.ExitOnError)

	/*
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
		ProcessTagCmd(*tagTags, *tagRecursive, tagCmd.Args())
	case "untag":
		untagCmd.Parse(os.Args[2:])
		ProcessUntagCmd(*untagTags, *untagRecursive, untagCmd.Args())
	case "list":
		listCmd.Parse(os.Args[2:])
		ProcessListCmd(listCmd.Args())
	case "lstag":
		lstagCmd.Parse(os.Args[2:])
		ProcessLstagCmd(lstagCmd.Args())
	case "ls":
		lsCmd.Parse(os.Args[2:])
		ProcessLsCmd(lsCmd.Args())
	case "mv":
		mvCmd.Parse(os.Args[2:])
		ProcessMvCmd(mvCmd.Args())
	case "cp":
		cpCmd.Parse(os.Args[2:])
		ProcessCpCmd(cpCmd.Args())
	case "rm":
		rmCmd.Parse(os.Args[2:])
		ProcessRmCmd(rmCmd.Args())

	}

}

//var tags is the -t flag
//var recursive is the -r flag
func ProcessTagCmd(tags string, recursive bool, args []string) {
	//Check for errors
	if tags == "" && len(args) <= 1 { //no flag, not enough args
		fmt.Println("Need at least one tag and one file.")
		return
	} else if tags != "" && len(args) < 1 { //flag, not enough args
		fmt.Println("No file entered.")
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
	for _, filePath := range args {
		//Check if directory (also checks if file exists)
		file, fileInfo, err := OpenFile(filePath)
		if err != nil {
			fmt.Println(filePath+": ", err)
			continue
		}
		defer file.Close()

		if fileInfo.IsDir() {
			if recursive {
				//Get the files in the directory
				subFiles, err := file.Readdir(-1)
				if err != nil {
					fmt.Println(filePath+": ", err)
					continue
				}
				subFilePaths := make([]string, 0)
				for _, subFileInfo := range subFiles {
					subFilePaths = append(subFilePaths, filepath.Join(filePath, subFileInfo.Name()))
				}

				//Send directory to be processed if directory is non-empty
				if len(subFilePaths) > 0 {
					ProcessTagCmd(tags, recursive, subFilePaths)
				}
			} else {
				fmt.Println(filePath + ": cannot tag a directory. Use -r flag to tag the files in the directory.")
			}

			//Since file is directory, it is handled
			//Continue to next file
			continue
		}

		//File is not a directory, so tag it
		AddFile(filePath)
		fid, err := GetFile(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}

		//add tags to file
		tidlist, err := GetOrCreateTags(taglist...)
		TagFile(fid, tidlist...)
	}

}

func ProcessUntagCmd(tags string, recursive bool, args []string) {
	//Check for errors
	if len(args) < 1 { //check if files passed
		fmt.Println("Need at least one file...")
		return
	}

	for _, filePath := range args {
		file, fileInfo, err := OpenFile(filePath)
		if err != nil {
			fmt.Println(filePath+": ", err)
			continue
		}
		defer file.Close()

		if fileInfo.IsDir() {
			if recursive {
				//Get the files in the directory
				subFiles, err := file.Readdir(-1)
				if err != nil {
					fmt.Println(filePath+": ", err)
					continue
				}
				subFilePaths := make([]string, 0)
				for _, subFileInfo := range subFiles {
					subFilePaths = append(subFilePaths, filepath.Join(filePath, subFileInfo.Name()))
				}

				//Send directory to be untagged
				if len(subFilePaths) > 0 {
					ProcessUntagCmd(tags, recursive, subFilePaths)
				}
			} else {
				fmt.Println(filePath + ": cannot untag a directory. Use -r flag to untag the files in the directory.")
			}

			//Directory is handled
			//Continue to next file
			continue
		}

		//Note: registers the file to ensure error returned
		//should only happen with actual errors
		AddFile(filePath)
		fid, err := GetFile(filePath)
		if err != nil {
			fmt.Println(filePath+": ", err)
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
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
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
		fmt.Println(filepath+": ", tags)
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
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Println(dir + ": directory does not exist")
			continue
		}

		//Check if directory
		file, fileInfo, err := OpenFile(dir)
		if err != nil {
			fmt.Println(dir+": ", err)
			continue
		}
		defer file.Close()

		if !fileInfo.IsDir() {
			fmt.Println(dir + ": not a directory")
			continue
		}

		//perform ls for all files
		fmt.Println("Directory " + dir + ":\n---------")
		filesInfo, err := file.Readdir(-1)
		filePaths := make([]string, 0)
		for _, subfileInfo := range filesInfo {
			filePaths = append(filePaths, filepath.Join(dir, subfileInfo.Name()))
		}
		//pass the file list over to list function
		ProcessListCmd(filePaths)
		fmt.Println("")
	}
}

func ProcessMvCmd(args []string) {
	//Check if at least two arguments
	if len(args) != 2 {
		fmt.Println("Need exactly one source and one destination.")
		return
	}

	//Get paths
	srcPath := args[0]
	destPath := args[1]

	//Rename
	err := os.Rename(srcPath, destPath)
	if err != nil {
		fmt.Println(srcPath+": ", err)
		return
	}

	//Get source file id
	srcFid, err := GetFile(srcPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Get destination file id
	AddFile(destPath)
	destFid, err := GetFile(destPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Copy tags
	err = CopyTagFile(srcFid, destFid)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Untag old file
	UntagAllFile(srcFid)
}

func ProcessRmCmd(args []string) {
	//Check for at least one arguments
	if len(args) < 1 {
		fmt.Println("Need at least one file.")
		return
	}

	//Get paths
	paths := args

	for _, path := range paths {
		//delete
		err := os.Remove(path)
		if err != nil {
			fmt.Println(path+": ", err)
			return
		}

		//Get file id
		fid, err := GetFile(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		//Untag file
		UntagAllFile(fid)
	}
}

func ProcessCpCmd(args []string) {

	//Check if at least two arguments
	//If not, error out
	if len(args) < 2 {
		fmt.Println("Need at least one source and one destination.")
		return
	}

	sourceFile, fileInfo, err := OpenFile(args[0]) //Open file
	if err != nil {
		fmt.Println(args[0] + ": cannot open file")
		return
	}
	defer sourceFile.Close()

	//Cannot copy directories
	if fileInfo.IsDir() {
		fmt.Println(args[0] + ": cannot copy directory")
		return
	}

	//Get source file id
	srcFid, err := GetFile(args[0])
	if err != nil {
		fmt.Println(args[0] + ": file not tagged.")
		return
	}

	//Get destination list
	destinations := args[1:]

	//Handle each copy
	for _, destinationPath := range destinations {
		//Open/create destination file
		destinationFile, err := os.OpenFile(destinationPath, os.O_RDWR|os.O_CREATE, fileInfo.Mode())
		if err != nil {
			fmt.Println("Error creating/opening file ", destinationPath, ": ", err)
			continue
		}
		defer destinationFile.Close()

		//Perform copy
		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			fmt.Println("Unable to copy file to ", destinationPath)
			continue
		}

		//Reset source file
		_, err = sourceFile.Seek(0, 0)
		if err != nil {
			fmt.Println("Seek failed on source file.")
			break //Break since seek failed; creating all other files will also fail.
		}

		//Copy tags
		AddFile(destinationPath)
		destFid, err := GetFile(destinationPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = CopyTagFile(srcFid, destFid)
		if err != nil {
			fmt.Println(err)
			return
		}

	}

}

//returns the file, fileinfo, error
func OpenFile(path string) (*os.File, os.FileInfo, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	fileInfo, err := file.Stat()
	return file, fileInfo, err
}
