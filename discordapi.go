package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func announceLocation(param string, newData UserInfoParsed, oldData UserInfoParsed, session *discordgo.Session) {
	switch param {
	case "login":
		message := fmt.Sprintf(phrasePicker("conf/login.txt"), newData.Login, newData.Location)
		fmt.Println(fmt.Sprintf("\t\tSending login for %s, on %s", newData.Login, newData.Location))
		_, err := session.ChannelMessageSend("277524661208612865", message)
		checkError(err)
	case "logout":
		message := fmt.Sprintf(phrasePicker("conf/logout.txt"), newData.Login)
		fmt.Println(fmt.Sprintf("\t\tSending logout for %s", newData.Login))
		_, err := session.ChannelMessageSend("277524661208612865", message)
		checkError(err)
	case "newPos":
		message := fmt.Sprintf(phrasePicker("conf/newPos.txt"), newData.Login, oldData.Location, newData.Location)
		fmt.Println(fmt.Sprintf("\t\tSending newPos for %s, from %s to %s", newData.Login, oldData.Location, newData.Location))
		_, err := session.ChannelMessageSend("277524661208612865", message)
		checkError(err)
	}
}

func announceProject(param string, newData UserInfoParsed, project string, session *discordgo.Session, oldData UserInfoParsed) {
	switch param {
	case "finished":
		message := fmt.Sprintf(phrasePicker("conf/finished.txt"), newData.Login, project, oldData.Level, newData.Level)
		fmt.Println(fmt.Sprintf("\t\tSending finished for %s, on %s", newData.Login, project))
		_, err := session.ChannelMessageSend("277524661208612865", message)
		checkError(err)
	case "started":
		message := fmt.Sprintf(phrasePicker("conf/started.txt"), newData.Login, project, phrasePicker("conf/proverb.txt"))
		fmt.Println(fmt.Sprintf("\t\tSending started for %s, on %s", newData.Login, project))
		_, err := session.ChannelMessageSend("277524661208612865", message)
		checkError(err)
	}
}

func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {

	botID, err := session.User("@me")
	checkError(err)

	if botID.ID == message.Author.ID {
		return
	}
	_, err = session.ChannelMessageSend(message.ChannelID, message.Content)
	checkError(err)
}
