package main

import (
	"path/filepath"
)

func AddTag(tag string) error {
	tid, _ := GetTag(tag)

	//If tag already exists, don't touch anything
	if tid != -1{
		return nil
	}

	statement, err := DB.Prepare("INSERT INTO tags (tag) VALUES ($1)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(tag)
	statement.Close()

	return err
}

func RemoveTag(tag string) error {
	return nil
}

func AddFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	fid, _ := GetFile(absPath)

	//If file already exists, don't touch anything
	if fid != -1 {
		return nil
	}

	statement, err := DB.Prepare("INSERT INTO files (filepath) VALUES ($1)")
	if err != nil {
		return err
	}

	_, err = statement.Exec(absPath)
	statement.Close()

	return err

	return nil
}

func TagFile(file int, tags ...int) error {
	statement, err := DB.Prepare("INSERT INTO filetags (file_id, tag_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	//handle each tag
	for _, tid := range tags {
		_, err = statement.Exec(file, tid)
		if err != nil {
			return err
		}
	}

	return nil
}

func UntagFile(file int, tags ...int) error {
	return nil
}

func GetTagsForFile(file int) ([]string, error) {
	tags := make([]string, 0)

	statement, err := DB.Prepare("SELECT tags.tag FROM tags JOIN filetags ON tags.id=filetags.tag_id WHERE filetags.file_id=$1")

	rows, err := statement.Query(file)
	if err != nil {
		return tags, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag string
		if err = rows.Scan(&tag); err != nil {
			return tags, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func GetFilesForTag(tag int) ([]string, error) {
	files := make([]string, 0)

	statement, err := DB.Prepare("SELECT files.filepath FROM files JOIN filetags on files.id=filetags.file_id WHERE filetags.tag_id=$1")

	rows, err := statement.Query(tag)
	if err != nil {
		return files, err
	}
	defer rows.Close()

	for rows.Next() {
		var file string
		if err = rows.Scan(&file); err != nil {
			return files, err
		}
		files = append(files, file)
	}

	return files, nil
}

func GetFile(path string) (int, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return -1, err
	}

	statement, err := DB.Prepare("SELECT id FROM files WHERE filepath = $1")
	if err != nil {
		return -1, err
	}

	row := statement.QueryRow(absPath)
	statement.Close()

	var id int = -1
	err = row.Scan(&id)

	return id, err

}

func GetTag(tag string) (int, error) {
	statement, err := DB.Prepare("SELECT id FROM tags WHERE tag = $1")
	if err != nil {
		return -1, err
	}

	row := statement.QueryRow(tag)
	statement.Close()

	var id int = -1
	err = row.Scan(&id)

	return id, err

}

func GetOrCreateTags(tags ...string) ([]int, error) {
	tids := make([]int, 0)

	for _, tag := range tags {
		tid, err := GetTag(tag)
		if tid == -1 { //If tag does not exist, add it
			AddTag(tag)
			tid, _ = GetTag(tag)
		} else { //If tag exists, check for errors
			if err != nil {
				return tids, err
			}
		}

		tids = append(tids, tid)
	}

	return tids, nil
}
