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
		return err
	}

	defer response.Body.Close()

	checkError(json.Unmarshal(body, &token))
	return nil
}

func getUserInfo(user string, token OAuthToken, userData UserInfo) (UserInfo, OAuthToken, error) {
	var url = fmt.Sprintf("https://api.intra.42.fr/v2/users/%s/", user)

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		return userData, token, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return userData, token, err
	}

	if res.Status == "429 Too Many Requests" {
		timeToSleep, _ := strconv.Atoi(res.Header["Retry-After"][0])
		time.Sleep(time.Duration(timeToSleep+2) * time.Second)
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			return userData, token, err
		}
	}

	if res.Status != "200 OK" {
		err = token.getToken()
		if err != nil {
			return userData, token, err
		}
		fmt.Println("42 Token acquired")
		fmt.Println("Expires in:", token.ExpiresIn)
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			return userData, token, err
		}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return userData, token, err
	}

	defer res.Body.Close()

	checkError(json.Unmarshal(body, &userData))

	return userData, token, nil
}
