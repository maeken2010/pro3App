package controllers

import (
	"encoding/json"
	"github.com/revel/revel"
	"github.com/mrjones/oauth"
	"pro3App/app/models"
)

type Application struct {
	*revel.Controller
}

var TWITTER = oauth.NewConsumer(
	"snWg3l6s2lJvZZ2y3IXxjxivv",
	"JSeK3pCLHtjPhnmEf0YpGHjMGpw1ZoFcbDhUrlNq8dOQcAhh04",
	oauth.ServiceProvider{
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	},
)


func (c Application) Index() revel.Result {
	return c.Render()
}

func (c Application) EnterDemo(user, demo string) revel.Result {
	users := getUser()

	if user != "" {
	 c.Validation.Required(user)
	 c.Validation.Required(demo)
	}else{
	resp, _ := TWITTER.Get(
	"https://api.twitter.com/1.1/account/verify_credentials.json",
	map[string]string{},
	users.AccessToken,
	)
	defer resp.Body.Close()
	account := struct {
		Name            string `json:"screen_name"`
		ProfileImageURL string `json:"profile_image_url"`
	}{}
	_ = json.NewDecoder(resp.Body).Decode(&account)
	// // 表示用情報をセット
	// setUserData(account.Name, account.ProfileImageURL)
	user = account.Name
	// demo = "websocket"
	}

	if c.Validation.HasErrors() {
		c.Flash.Error("Please choose a nick name and the demonstration type.")
		return c.Redirect(Application.Index)
	}

	return c.Redirect("/websocket/room?user=%s", user)

}

func (c Application) Authenticate(oauth_verifier string) revel.Result {
	user := getUser()
	if oauth_verifier != "" {
		// We got the verifier; now get the access token, store it and back to index
		accessToken, err := TWITTER.AuthorizeToken(user.RequestToken, oauth_verifier)
		if err == nil {
			user.AccessToken = accessToken
		} else {
			revel.ERROR.Println("Error connecting to twitter:", err)
		}
		return c.Redirect(Application.EnterDemo)
	}

	requestToken, url, err := TWITTER.GetRequestTokenAndUrl("http://localhost:9000/Application/Authenticate")
	if err == nil {
		// We received the unauthorized tokens in the OAuth object - store it before we proceed
		user.RequestToken = requestToken
		return c.Redirect(url)
	} else {
		revel.ERROR.Println("Error connecting to twitter:", err)
	}
	return c.Redirect(Application.EnterDemo)
}


func getUser() *models.User {
	return models.FindOrCreate("guest")
}

func init() {
	TWITTER.Debug(true)
}

