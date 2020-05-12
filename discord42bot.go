package main

import (
	"fmt"
)

func initUsers(api Api42) {

	api.UserData.getUserInfo("mlabouri", api.Token)
	fmt.Println("Got info from mlabouri")
	fmt.Println(api.UserData)
}

func main() {

	api := Api42{}

	fmt.Println("Started Bot")

	api.Token.getToken()
	fmt.Println("42 Token acquired")
	fmt.Println("Expires in:", api.Token.ExpiresIn)

}
