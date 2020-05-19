package main

import (
	"time"
)

type Project struct {
	ProjectName   string
	ProjectStatus string
}

type UserInfoParsed struct {
	Login           string
	Email           string
	Location        string
	CorrectionPoint int
	Wallet          int
	BlackHole       int
	Level           float64
	Projects        map[string]Project
}

func processUserInfo(userData UserInfo) (UserInfoParsed, error) {

	project := Project{}
	userDataParsed := UserInfoParsed{}

	userDataParsed.Login = userData.Login
	userDataParsed.Email = userData.Email
	userDataParsed.Wallet = userData.Wallet
	userDataParsed.CorrectionPoint = userData.CorrectionPoint

	userDataParsed.Location = userData.Location
	if userData.Location == "" {
		userDataParsed.Location = "null"
	}

	for _, cursus := range userData.CursusUsers {
		if cursus.CursusID == 21 {
			userDataParsed.BlackHole = int(cursus.BlackHoledAt.Sub(time.Now()).Hours() / 24)
			userDataParsed.Level = cursus.Level
		}
	}

	userDataParsed.Projects = make(map[string]Project)

	for _, projectRaw := range userData.ProjectsUsers {
		if projectRaw.CursusIds[0] != 9 {
			project.ProjectName = projectRaw.Project.Name
			project.ProjectStatus = projectRaw.Status
			userDataParsed.Projects[projectRaw.Project.Name] = project
		}
	}

	return userDataParsed, nil
}
