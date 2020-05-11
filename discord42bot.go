package main

import (
	api42 "discord42bot/api42"
	"fmt"
)

func main(){
	discord42Bot()
}

func discord42Bot() {
	fmt.Println("Started Bot")

	token := api42.OAuthToken{
		AccessToken: "test",
		TokenType:   "tes",
		ExpiresIn:   15,
	}

	fmt.Println("42 Token acquired")
	fmt.Println("Expires in:", token.ExpiresIn)


}