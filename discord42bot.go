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

func writeUsers(api Api42, session *discordgo.Session, callNbr int) Api42 {

	var userList = os.Args

	for i, user := range userList[1:] {
		userData := UserInfo{}
		userDataParsed := UserInfoParsed{}
		var err error

		userData, api.Token, err = getUserInfo(user, api.Token, userData)
		if err != nil {
			logError(err)
			continue
		}
		fmt.Println(fmt.Sprintf("Request %06d:\n\tGot raw data from %s", i+((len(userList)-1)*callNbr), user))
		userDataParsed, err = processUserInfo(userData)
		if err != nil {
			logError(err)
			continue
		}
		fmt.Println("\tProcessed raw data")
		err = checkUserFile(user, userDataParsed, session)
		if err != nil {
			logError(err)
			continue
		}
		userDataToDB(user)
		time.Sleep(3000 * time.Millisecond)
	}

	return api
}

func main() {
	api := Api42{}

	fmt.Println("Started init")

	err := godotenv.Load("dev.env")
	checkError(err)

	err = api.Token.getToken()
	checkError(err)
	fmt.Println("42 Token acquired")
	fmt.Println("Expires in:", api.Token.ExpiresIn)

	discordBot, err := discordgo.New("Bot " + os.Getenv("BOTTOKEN"))
	checkError(err)
	fmt.Println("Discord bot created")

	discordBot.AddHandler(messageHandler)

	err = discordBot.Open()
	checkError(err)
	fmt.Println("Discord Bot up and running")

	startApi()

	setupCloseHandler(discordBot)
	var callNbr = 0
	for {
		api = writeUsers(api, discordBot, callNbr)
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
