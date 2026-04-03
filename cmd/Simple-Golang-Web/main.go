package main

import (
	"Simple-Golang-Web/internal/auth"
	"Simple-Golang-Web/internal/handler"
	"Simple-Golang-Web/views/layouts"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"golang.org/x/oauth2"
)

// Main function
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	app := pocketbase.New()

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Static files
		e.Router.GET("/static/{path...}", func(c *core.RequestEvent) error {
			http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(c.Response, c.Request)
			return nil
		})

		// Home
		e.Router.GET("/", func(c *core.RequestEvent) error {
			var isAdmin bool
			var isLoggedIn bool

			data := []string{
				"Anderson Tirza Liman",
				"Angelo Benhanan Abinaya Fuun",
				"Jovanus Irwan Susanto",
				"Raymundo Rafaelito Maryos Von Woloblo",
				"Ilmu Komputer 2024",
				"Ilmu Komputer 2024",
				"Ilmu Komputer 2024",
				"Ilmu Komputer 2024",
				"2406355893",
				"2406495432",
				"2406434140",
				"240640462",
			}

			adminEmails := []string{
				"andersontirza@gmail.com",
				"a.b.abinaya.fuun@gmail.com",
				"jovanusirwan@gmail.com",
				"woloblocrafael@gmail.com",
			}

			isAdmin = false
			isLoggedIn = false
			email := getLoggedInEmail(c, app)
			if email != "" {
				isLoggedIn = true

				if slices.Contains(adminEmails, email) {
					isAdmin = true
				}
			}

			return handler.Render(c, http.StatusOK, layouts.BaseLayout(data, isAdmin, isLoggedIn))
		})

		// Redirect user to Google
		e.Router.GET("/auth/google/login", func(c *core.RequestEvent) error {
			url := auth.GoogleOAuthConfig().AuthCodeURL("", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "select_account"))
			http.Redirect(c.Response, c.Request, url, http.StatusTemporaryRedirect)
			return nil
		})

		// Handle Google's callback
		e.Router.GET("/auth/google/callback", func(c *core.RequestEvent) error {
			code := c.Request.FormValue("code")
			if code == "" {
				http.Error(c.Response, "Missing code", http.StatusBadRequest)
				return nil
			}

			token, err := auth.GoogleOAuthConfig().Exchange(context.Background(), code)
			if err != nil {
				http.Error(c.Response, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
				return nil
			}

			client := auth.GoogleOAuthConfig().Client(context.Background(), token)
			resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
			if err != nil {
				http.Error(c.Response, "Failed to get user info", http.StatusInternalServerError)
				return nil
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			var googleUser GoogleUser
			json.Unmarshal(body, &googleUser)

			record, err := app.FindAuthRecordByEmail("users", googleUser.Email)
			if err != nil {
				collection, err := app.FindCollectionByNameOrId("users")
				if err != nil {
					http.Error(c.Response, "Collection not found", http.StatusInternalServerError)
					return nil
				}
				record = core.NewRecord(collection)
				record.Set("email", googleUser.Email)
				record.SetPassword(generateState())
				record.Set("emailVisibility", true)
			}

			record.Set("name", googleUser.Name)
			if err := app.Save(record); err != nil {
				http.Error(c.Response, "Failed to save user: "+err.Error(), http.StatusInternalServerError)
				return nil
			}

			authToken, err := record.NewAuthToken()
			if err != nil {
				http.Error(c.Response, "Failed to create auth token", http.StatusInternalServerError)
				return nil
			}

			http.SetCookie(c.Response, &http.Cookie{
				Name:     "auth_token",
				Value:    authToken,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   60 * 60 * 24,
			})

			http.Redirect(c.Response, c.Request, "/", http.StatusTemporaryRedirect)
			return nil
		})

		// Logout user
		e.Router.GET("/auth/logout", func(c *core.RequestEvent) error {
			http.SetCookie(c.Response, &http.Cookie{
				Name:     "auth_token",
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				MaxAge:   -1,
			})
			http.Redirect(c.Response, c.Request, "/", http.StatusTemporaryRedirect)
			return nil
		})

		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// User model for Google OAuth
type GoogleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// Generates random state string for OAuth
func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Gets the logged-in user's email
func getLoggedInEmail(c *core.RequestEvent, app *pocketbase.PocketBase) string {
	cookie, err := c.Request.Cookie("auth_token")
	if err != nil {
		return ""
	}

	record, err := app.FindAuthRecordByToken(cookie.Value)
	if err != nil {
		return ""
	}

	return record.GetString("email")
}
