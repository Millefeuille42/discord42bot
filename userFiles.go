package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func checkUserFile(user string, userData UserInfoParsed) {
	_, err := os.Stat(fmt.Sprintf("./data/%s.json", user))
	if os.IsNotExist(err) {
		fmt.Println("\tData file not found")

		file, err := os.Create(fmt.Sprintf("./data/%s.json", user))
		checkError(err)
		defer file.Close()

		fmt.Println("\tData file created")
	}

	jsonData, err := json.MarshalIndent(userData, "", "\t")
	checkError(err)

	err = ioutil.WriteFile(fmt.Sprintf("./data/%s.json", user), jsonData, 0644)
	checkError(err)
	fmt.Println("\tData written to file")
}
