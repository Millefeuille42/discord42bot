package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func setVarsToMessage(phrase string, newData UserInfoParsed, oldData UserInfoParsed, project string) string {
	phrase = strings.ReplaceAll(phrase, "#{userName}", newData.Login)
	phrase = strings.ReplaceAll(phrase, "#{project}", project)
	phrase = strings.ReplaceAll(phrase, "#{proverb}", phrasePicker("conf/proverbs.txt"))
	phrase = strings.ReplaceAll(phrase, "#{oldLocation}", oldData.Location)
	phrase = strings.ReplaceAll(phrase, "#{newLocation}", newData.Location)
	phrase = strings.ReplaceAll(phrase, "#{oldLevel}", fmt.Sprintf("%.2f", oldData.Level))
	phrase = strings.ReplaceAll(phrase, "#{newLevel}", fmt.Sprintf("%.2f", newData.Level))

	return phrase
}

func announceLocation(param string, newData, oldData UserInfoParsed, session *discordgo.Session) {
	switch param {
	case "login":
		message := setVarsToMessage(phrasePicker("conf/login.txt"), newData, oldData, "")
		fmt.Println(fmt.Sprintf("\t\tSending login for %s, on %s", newData.Login, newData.Location))
		_, err := session.ChannelMessageSend("710820070284066822", message)
		checkError(err)
	case "logout":
		message := setVarsToMessage(phrasePicker("conf/logout.txt"), newData, oldData, "")
		fmt.Println(fmt.Sprintf("\t\tSending logout for %s", newData.Login))
		_, err := session.ChannelMessageSend("710820070284066822", message)
		checkError(err)
	case "newPos":
		message := setVarsToMessage(phrasePicker("conf/newPos.txt"), newData, oldData, "")
		fmt.Println(fmt.Sprintf("\t\tSending newPos for %s, from %s to %s", newData.Login, oldData.Location, newData.Location))
		_, err := session.ChannelMessageSend("710820070284066822", message)
		checkError(err)
	}
}

func announceProject(param, project string, newData, oldData UserInfoParsed, session *discordgo.Session) {
	switch param {
	case "finished":
		message := setVarsToMessage(phrasePicker("conf/finished.txt"), newData, oldData, project)
		fmt.Println(fmt.Sprintf("\t\tSending finished for %s, on %s", newData.Login, project))
		_, err := session.ChannelMessageSend("710820070284066822", message)
		checkError(err)
	case "started":
		message := setVarsToMessage(phrasePicker("conf/started.txt"), newData, oldData, project)
		fmt.Println(fmt.Sprintf("\t\tSending started for %s, on %s", newData.Login, project))
		_, err := session.ChannelMessageSend("710820070284066822", message)
		checkError(err)
	}
}

func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {

	botID, err := session.User("@me")
	checkError(err)

	if botID.ID == message.Author.ID {
		return
	}

	if message.Content == "!leaderboard" {
		leaderboard(session, message)
	}

	if strings.HasPrefix(message.Content, "!roadmap") {
		arg := strings.Split(message.Content, "-")
		if len(arg) > 1 {
			roadmap(session, message, arg[1])
		} else {
			roadmap(session, message, "in_progress")
		}
	}

	if strings.HasPrefix(message.Content, "!template") {
		arg := strings.Split(message.Content, "-")
		if len(arg) > 1 {
			template(session, message, arg[1])
		} else {
			template(session, message, "bin")
		}
	}

	if strings.HasPrefix(message.Content, "!user") {
		arg := strings.Split(message.Content, " ")
		if len(arg) > 1 {
			sendUser(session, message, arg[1])
		}
	}
}
