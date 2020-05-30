package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"io/ioutil"
	"os"
)

type User struct {
	gorm.Model
	Login           string
	Email           string
	Location        string
	CorrectionPoint int
	Wallet          int
	BlackHole       int
	Level           float64
}

type OverTimeData struct {
	gorm.Model
	Login           string
	BlackHole       int
	CorrectionPoint int
	Level           float64
}

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
			announceProject("started", project, newUserData, fileDataJson, session)
		} else if status := newProjectData.ProjectStatus; status != oldProjectData.ProjectStatus {
			if status == "finished" && newUserData.Level != fileDataJson.Level {
				announceProject(status, project, newUserData, fileDataJson, session)
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

func userDataToDB(user string) {
	var queryUser User
	userData := UserInfoParsed{}

	fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/%s.json", user))
	if err != nil {
		return
	}
	err = json.Unmarshal(fileData, &userData)
	if err != nil {
		return
	}

	db, err := gorm.Open("postgres", fmt.Sprintf("connect_timeout=10 host=%s user=%s dbname=segbot password=%s sslmode=disable",
		os.Getenv("DBHOST"),
		os.Getenv("DBUSER"),
		os.Getenv("DBPASSWORD")))
	checkError(err)

	db.AutoMigrate(&User{})
	db.Find(&queryUser, "login = ?", user)

	exists := queryUser.Login

	queryUser.Login = userData.Login
	queryUser.Email = userData.Email
	queryUser.Location = userData.Location
	queryUser.CorrectionPoint = userData.CorrectionPoint
	queryUser.Wallet = userData.Wallet
	queryUser.BlackHole = userData.BlackHole
	queryUser.Level = userData.Level

	if exists == "" {
		db.Create(&queryUser)
	} else {
		db.Model(&queryUser).Where("login = ?", user).Save(&queryUser)
	}

	db.AutoMigrate(&OverTimeData{})
	db.Create(&OverTimeData{
		Login:           queryUser.Login,
		BlackHole:       queryUser.BlackHole,
		CorrectionPoint: queryUser.CorrectionPoint,
		Level:           queryUser.Level,
	})

	defer db.Close()
}
