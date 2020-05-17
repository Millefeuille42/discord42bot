package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
)

func compareData(fileData []byte, newUserData UserInfoParsed, session *discordgo.Session) {
	fileDataJson := UserInfoParsed{}

	err := json.Unmarshal(fileData, &fileDataJson)
	checkError(err)

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
}

func checkUserFile(user string, userData UserInfoParsed, session *discordgo.Session) {
	var path = fmt.Sprintf("./data/%s.json", user)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("\tData file not found")

		file, err := os.Create(path)
		checkError(err)
		defer file.Close()

		fmt.Println("\tData file created")
	} else {
		fileData, err := ioutil.ReadFile(path)
		checkError(err)
		compareData(fileData, userData, session)
	}

	jsonData, err := json.MarshalIndent(userData, "", "\t")
	checkError(err)

	err = ioutil.WriteFile(path, jsonData, 0644)
	checkError(err)
	fmt.Println("\tData written to file")
}
