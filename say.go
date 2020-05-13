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

func leaderboard(session *discordgo.Session, message *discordgo.MessageCreate) {
	var userList = os.Args
	var userPair = make([]levelNamePair, len(userList)-1)
	var leadMessage = fmt.Sprintf("@%s", message.Author.ID)
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
		leadMessage = fmt.Sprintf("%s\n%2d: %-9s%f", message, i+1, user.name, user.level)
	}
	_, err := session.ChannelMessageSend(message.ChannelID, leadMessage)
	checkError(err)
}
