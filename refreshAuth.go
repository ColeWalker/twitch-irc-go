package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

//AccessToken stores a single access token from refreshAuth
type AccessToken struct {
	Token string `json:"access_token"`
}

//refresh authentication
func refreshAuth(refreshToken string, clientID string, secret string) string {
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s", refreshToken, clientID, secret)
	resp, err := http.Post(url, "application/json", nil)

	if err != nil {
		fmt.Println("Error refreshing auth token")
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Error in parsing refresh token response")
		log.Fatalln(err)
	}

	var NewToken AccessToken
	err = json.Unmarshal(body, &NewToken)

	if err != nil {
		log.Println("Error in unmarshalling refresh token response")
		log.Fatalln(err)
	}

	if NewToken.Token == "" {
		log.Println("Didn't receive refresh token from twitch")
		log.Fatalln(string(body))
	}
	return NewToken.Token
}
