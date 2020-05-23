package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/openpgp/errors"
	"io/ioutil"
	"log"
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

func (token *OAuthToken) getToken() error {

	var url = "https://api.intra.42.fr/oauth/token"
	var uid = os.Getenv("APIUID")
	var secret = os.Getenv("APISECRET")
	var req = fmt.Sprintf("%s?grant_type=client_credentials&client_id=%s&client_secret=%s", url, uid, secret)

	response, err := http.Post(req, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte("")))
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		defer response.Body.Close()
		return err
	}

	defer response.Body.Close()

	err = json.Unmarshal(body, &token)
	return err
}

func getUserInfo(user string, token OAuthToken, userData UserInfo) (UserInfo, OAuthToken, error) {
	var url = fmt.Sprintf("https://api.intra.42.fr/v2/users/%s/", user)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		log.Print("http.NewRequest(\"GET\", url, bytes.NewBuffer([]byte(\"\")))")
		return userData, token, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print("http.DefaultClient.Do(req)")
		return userData, token, err
	}

	if res.Status == "429 Too Many Requests" {
		defer res.Body.Close()
		timeToSleep, _ := strconv.Atoi(res.Header["Retry-After"][0])
		time.Sleep(time.Duration(timeToSleep+2) * time.Second)
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Print("http.DefaultClient.Do(req)")
			return userData, token, err
		}
	}

	if res.Status != "200 OK" {
		defer res.Body.Close()
		err = token.getToken()
		log.Print("200 !OK")
		return userData, token, errors.UnsupportedError("API Error")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		defer res.Body.Close()
		log.Print("ioutil.ReadAll(res.Body)")
		return userData, token, err
	}

	defer res.Body.Close()

	err = json.Unmarshal(body, &userData)
	return userData, token, err
}
