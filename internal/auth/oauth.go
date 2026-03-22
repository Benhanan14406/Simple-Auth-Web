package auth

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleOAuthConfig() *oauth2.Config {
	redirectURL := os.Getenv("OAUTH_REDIRECT_URL")
	if redirectURL == "" {
		redirectURL = "http://localhost:8090/auth/google/callback"
	}

	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
