package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Api42 struct {
	Token    OAuthToken
	UserData UserInfo
}

type OAuthToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type UserInfo struct {
	Email           string `json:"email"`
	Login           string `json:"login"`
	Name            string `json:"displayname"`
	Location        string `json:"location"`
	CorrectionPoint int    `json:"correction_point"`
	Wallet          int    `json:"wallet"`
	CursusUsers     []struct {
		CursusID     int       `json:"cursus_id"`
		Level        float64   `json:"level"`
		BlackHoledAt time.Time `json:"blackholed_at"`
	} `json:"cursus_users"`
	ProjectsUsers []struct {
		Status    string `json:"status"`
		CursusIds []int  `json:"cursus_ids"`
		Project   struct {
			Name string `json:"name"`
		} `json:"project"`
	} `json:"projects_users"`
}

func (token *OAuthToken) getToken() {

	var url = "https://api.intra.42.fr/oauth/token"
	var uid = os.Getenv("APIUID")
	var secret = os.Getenv("APISECRET")
	var req = fmt.Sprintf("%s?grant_type=client_credentials&client_id=%s&client_secret=%s", url, uid, secret)

	response, err := http.Post(req, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte("")))
	checkError(err)

	body, err := ioutil.ReadAll(response.Body)
	checkError(err)

	defer response.Body.Close()

	checkError(json.Unmarshal(body, &token))
}

func getUserInfo(user string, token OAuthToken, userData UserInfo) UserInfo {
	var url = fmt.Sprintf("https://api.intra.42.fr/v2/users/%s/", user)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	checkError(err)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	checkError(err)

	if res.Status == "429 Too Many Requests" {
		timeToSleep, _ := strconv.Atoi(res.Header["Retry-After"][0])
		time.Sleep(time.Duration(timeToSleep+2) * time.Second)
		res, err = http.DefaultClient.Do(req)
		checkError(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	checkError(err)

	defer res.Body.Close()

	checkError(json.Unmarshal(body, &userData))

	return userData
}
