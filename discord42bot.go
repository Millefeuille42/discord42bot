package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func writeUsers(api Api42, session *discordgo.Session, callNbr int) {

	var userList = os.Args
	userDataParsed := UserInfoParsed{}

	for i, user := range userList {
		if i != 0 {
			var oldTime = time.Now()
			api.UserData.getUserInfo(user, api.Token)
			fmt.Println(fmt.Sprintf("Request %06d:\n\tGot raw data from %s", i+((len(userList)-1)*callNbr), user))

			userDataParsed = processUserInfo(api.UserData)
			fmt.Println("\tProcessed raw data")

			checkUserFile(user, userDataParsed, session)
			time.Sleep(3000*time.Millisecond - time.Now().Sub(oldTime))
		}
	}
}

func main() {
	api := Api42{}

	fmt.Println("Started init")

	err := godotenv.Load()
	checkError(err)

	api.Token.getToken()
	fmt.Println("42 Token acquired")
	fmt.Println("Expires in:", api.Token.ExpiresIn)

	discordBot, err := discordgo.New("Bot " + os.Getenv("BOTTOKEN"))
	checkError(err)
	fmt.Println("Discord bot created")

	discordBot.AddHandler(messageHandler)

	err = discordBot.Open()
	checkError(err)
	fmt.Println("Discord Bot up and running")

	setupCloseHandler(discordBot)
	var callNbr = 0
	for {
		writeUsers(api, discordBot, callNbr)
		callNbr++
	}
}

func setupCloseHandler(session *discordgo.Session) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		_ = session.Close()
		os.Exit(0)
	}()
}
