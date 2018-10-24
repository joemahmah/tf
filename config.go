package main

import (
	"os/user"
	"os"
	"path/filepath"
	"encoding/json"
	"bufio"
	"fmt"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var Config TFConfig
var DB *sql.DB

const ConfigDirName = ".tf_data"

type TFConfig struct {
	DBPath string
}

func GetHomeDirPath() (string, error) {
	currentUser, err := user.Current()

	if err != nil {
		return "", err
	}

	return currentUser.HomeDir, nil
}

func LoadConfig() error {
	homeDir, err := GetHomeDirPath()

	if err != nil {
		return err
	}

	configDirPath := filepath.Join(homeDir, ConfigDirName)

	//Create folder if needed
	if _,err = os.Stat(configDirPath); os.IsNotExist(err){
		fmt.Println("No config folder detected...\nCreated config folder.")
		os.Mkdir(configDirPath, 0770)
	}

	configFilePath := filepath.Join(configDirPath, "config.json")

	//create config if needed
	if _,err = os.Stat(configFilePath); os.IsNotExist(err){
		fmt.Println("No config file detected...\nCreating config file.")
		err = CreateConfig(configFilePath)
		if err != nil {
			return err
		}
	}

	//Load config
	fmt.Println("Loading config file.")
	configFile, err := os.Open(configFilePath)

	if err != nil {
		fmt.Println("Error loading config file!")
		return err
	}

	defer configFile.Close()

	configBufferedReader := bufio.NewReader(configFile)

	json.NewDecoder(configBufferedReader).Decode(&Config)

	//Load DB
	dbFilePath := filepath.Join(configDirPath, Config.DBPath)

	fmt.Println("Loading database.")
	err = LoadDB(dbFilePath)

	return err

}

func CreateConfig(path string) error {

	baseConfig := TFConfig{DBPath: "db.sqlite3"}

	configFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0660)

	if err != nil {
		fmt.Println("Error creating config.")
		return err
	}

	defer configFile.Close()

	configBufferedWriter := bufio.NewWriter(configFile)
	defer configBufferedWriter.Flush()

	json.NewEncoder(configBufferedWriter).Encode(&baseConfig)

	return nil
}

func LoadDB(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)

	if err != nil {
		fmt.Println("Error loading database.")
		return err
	}

	statement, err := DB.Prepare("CREATE TABLE IF NOT EXISTS files (id INTEGER PRIMARY KEY, filepath TEXT NOT NULL)")
	if err != nil {
		fmt.Println("SQL error.")
		return err
	}
	statement.Exec()
	statement.Close()

	statement, err = DB.Prepare("CREATE TABLE IF NOT EXISTS tags (id INTEGER PRIMARY KEY, tag TEXT NOT NULL)")
	if err != nil {
		fmt.Println("SQL error.")
		return err
	}
	statement.Exec()
	statement.Close()

	statement, err = DB.Prepare("CREATE TABLE IF NOT EXISTS filetags (tag_id INTEGER NOT NULL, file_id INTEGER NOT NULL, PRIMARY KEY (tag_id, file_id))")
	if err != nil {
		fmt.Println("SQL error.")
		return err
	}
	statement.Exec()
	statement.Close()

	statement, err = DB.Prepare("CREATE TABLE IF NOT EXISTS tagchildren (tag_id INTEGER NOT NULL, child_id INTEGER NOT NULL, PRIMARY KEY (tag_id, child_id))")
	if err != nil {
		fmt.Println("SQL error.")
		return err
	}
	statement.Exec()
	statement.Close()

	return nil
}
