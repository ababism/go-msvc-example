package service

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func InitializeOAuthGoogle() {
	//GOOD
	OauthConfGl.ClientID = viper.GetString("google.clientID")
	logrus.Info(OauthConfGl.ClientID)
	OauthConfGl.ClientSecret = viper.GetString("google.clientSecret")
	logrus.Info(OauthConfGl.ClientSecret)

	initializeAuthService()
}

var (
	OauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
)
