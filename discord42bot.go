package main

import (
	"fmt"
	"os"
	"time"
)

func writeUsers(api Api42) {

	var userList = os.Args
	userDataParsed := UserInfoParsed{}

	for i, user := range userList {
		if i != 0 {
			var oldTime = time.Now()
			api.UserData.getUserInfo(user, api.Token)
			fmt.Println(fmt.Sprintf("Request %03d:\n\tGot raw data from %s", i, user))

			userDataParsed = processUserInfo(api.UserData)
			fmt.Println("\tProcessed raw data")

			checkUserFile(user, userDataParsed)
			time.Sleep(500*time.Millisecond - time.Now().Sub(oldTime))
		}
	}
}

func main() {

	api := Api42{}

	fmt.Println("Started Bot")
	api.Token.getToken()
	fmt.Println("42 Token acquired")
	fmt.Println("Expires in:", api.Token.ExpiresIn)

	writeUsers(api)
}
