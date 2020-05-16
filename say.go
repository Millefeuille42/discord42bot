package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"sort"
)

type levelNamePair struct {
	name  string
	level float64
}

func roadmap(session *discordgo.Session, message *discordgo.MessageCreate, status string) {
	roadMessage := ""
	userList := os.Args
	userDataParsed := UserInfoParsed{}
	projectList := make(map[string]string, 0)

	if !Find([]string{"finished", "in_progress"}, status) {
		return
	}

	for i, user := range userList {
		if i == 0 {
			continue
		}
		fileData, err := ioutil.ReadFile(fmt.Sprintf("data/%s.json", user))
		checkError(err)
		err = json.Unmarshal(fileData, &userDataParsed)
		checkError(err)

		for _, project := range userDataParsed.Projects {
			if project.ProjectStatus == status {
				if _, ok := projectList[project.ProjectName]; !ok {
					projectList[project.ProjectName] = "\n\t| " + user
				} else {
					projectList[project.ProjectName] = fmt.Sprintf("%s\n\t| %s", projectList[project.ProjectName], user)
				}
			}
		}
	}
	for projectName, projectUsers := range projectList {
		roadMessage = fmt.Sprintf("%s\n\n%s%10s", roadMessage, projectName, projectUsers)
	}
	roadMessage = fmt.Sprintf("<@%s>, Roadmap for '%s'```%s```", message.Author.ID, status, roadMessage)
	_, err := session.ChannelMessageSend(message.ChannelID, roadMessage)
	checkError(err)
}

func leaderboard(session *discordgo.Session, message *discordgo.MessageCreate) {
	var leadMessage = ""
	userList := os.Args
	userPair := make([]levelNamePair, len(userList)-1)
	userDataParsed := UserInfoParsed{}

	for i, user := range userList {
		if i != 0 {
			fileData, err := ioutil.ReadFile(fmt.Sprintf("data/%s.json", user))
			checkError(err)
			err = json.Unmarshal(fileData, &userDataParsed)
			checkError(err)
			userPair = append(userPair, levelNamePair{userDataParsed.Login, userDataParsed.Level})
		}
	}
	sort.Slice(userPair, func(i, j int) bool {
		return userPair[i].level > userPair[j].level
	})

	for i, user := range userPair {
		if i > (len(userList) - 2) {
			break
		}
		leadMessage = fmt.Sprintf("%s\n%2d: %-15s%.2f", leadMessage, i+1, user.name, user.level)
	}
	leadMessage = fmt.Sprintf("<@%s>```%s```", message.Author.ID, leadMessage)
	_, err := session.ChannelMessageSend(message.ChannelID, leadMessage)
	checkError(err)
}
