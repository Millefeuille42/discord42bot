package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
)

func compareData(fileData []byte, newUserData UserInfoParsed, session *discordgo.Session) error {
	fileDataJson := UserInfoParsed{}

	err := json.Unmarshal(fileData, &fileDataJson)
	if err != nil {
		return err
	}

	if fileDataJson.Location != newUserData.Location {
		switch {
		case fileDataJson.Location == "null":
			announceLocation("login", newUserData, fileDataJson, session)
		case newUserData.Location == "null":
			announceLocation("logout", newUserData, fileDataJson, session)
		default:
			announceLocation("newPos", newUserData, fileDataJson, session)
		}
	}

	for project, oldProjectData := range fileDataJson.Projects {
		if _, exists := newUserData.Projects[project]; !exists {
			newUserData.Projects[project] = oldProjectData
		}
	}

	for project, newProjectData := range newUserData.Projects {
		if oldProjectData, exists := fileDataJson.Projects[project]; !exists {
			announceProject("started", newUserData, project, session, fileDataJson)
		} else if status := newProjectData.ProjectStatus; status != oldProjectData.ProjectStatus {
			if status == "finished" {
				announceProject(status, newUserData, project, session, fileDataJson)
			}
		}
	}
	return nil
}

func checkUserFile(user string, userData UserInfoParsed, session *discordgo.Session) error {
	var path = fmt.Sprintf("./data/%s.json", user)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("\tData file not found")

		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		fmt.Println("\tData file created")
	} else {
		fileData, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		err = compareData(fileData, userData, session)
		if err != nil {
			return err
		}
	}

	jsonData, err := json.MarshalIndent(userData, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}
	fmt.Println("\tData written to file")
	return nil
}
