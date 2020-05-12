package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
)

func compareData(fileData []byte, userData UserInfoParsed, session *discordgo.Session) {
	fileDataJson := UserInfoParsed{}

	err := json.Unmarshal(fileData, &fileDataJson)
	checkError(err)

	if fileDataJson.Location != userData.Location {
		switch {
		case fileDataJson.Location == "null":
			announceLocation("login", userData, fileDataJson, session)
		case userData.Location == "null":
			announceLocation("logout", userData, fileDataJson, session)
		default:
			announceLocation("newPos", userData, fileDataJson, session)
		}
	}

	for project, newProjectData := range userData.Projects {
		if oldProjectData, exists := fileDataJson.Projects[project]; !exists {
			announceProject("started", userData, project, session)
		} else if status := newProjectData.ProjectStatus; status != oldProjectData.ProjectStatus {
			if status == "finished" {
				announceProject(status, userData, project, session)
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
