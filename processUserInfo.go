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
	Level           float64
	BlackHole       time.Duration
	Projects        map[string]Project
}

func processUserInfo(userData UserInfo) UserInfoParsed {

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

	userDataParsed.BlackHole = userData.CursusUsers[1].BlackHoledAt.Sub(time.Now())
	userDataParsed.Level = userData.CursusUsers[1].Level

	userDataParsed.Projects = make(map[string]Project)

	for _, projectRaw := range userData.ProjectsUsers {
		if projectRaw.CursusIds[0] != 9 {
			project.ProjectName = projectRaw.Project.Name
			project.ProjectStatus = projectRaw.Status
			userDataParsed.Projects[projectRaw.Project.Name] = project
		}
	}

	return userDataParsed
}
