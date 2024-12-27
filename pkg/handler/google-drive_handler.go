package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"web-service/config"
	googledrive "web-service/pkg/google-drive"
	"web-service/pkg/utils"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

func getAuthURL() string {
	return googledrive.OauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func handleGoogleDriveAuth(w http.ResponseWriter, r *http.Request) {
	url := getAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleDriveCallback(w http.ResponseWriter, r *http.Request) utils.Response {
	code := r.URL.Query().Get("code")
	if code == "" {
		return utils.BadRequestError("Code not found in the request", nil)
	}

	token, err := googledrive.OauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		return utils.BadRequestError("Failed to exchange token", nil)
	}

	tokenPath := config.Env.GOOGLE_DRIVE_TOKEN_PATH
	if tokenPath == "" {
		log.Fatal("Missing GOOGLE_DRIVE_TOKEN_PATH environment variable")
	}

	if err := saveToken(token, tokenPath); err != nil {
		log.Printf("Failed to save token: %v", err)
		return utils.InternalServerError("Failed to save token: " + err.Error())
	}

	return utils.CreatedResponse("Token successfully saved", nil)
}

func saveToken(token *oauth2.Token, tokenPath string) error {
	dir := filepath.Dir(tokenPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("unable to create directory for token file: %v", err)
		}
	}

	f, err := os.Create(tokenPath)
	if err != nil {
		return fmt.Errorf("unable to create token file: %v", err)
	}
	defer f.Close()

	// Ghi token vào file dưới dạng JSON
	if err := json.NewEncoder(f).Encode(token); err != nil {
		return fmt.Errorf("unable to write token to file: %v", err)
	}

	log.Printf("Token saved to file: %s", tokenPath)
	return nil
}

func getClient() *http.Client {
	tokenPath := config.Env.GOOGLE_DRIVE_TOKEN_PATH
	if tokenPath == "" {
		log.Fatal("Missing GOOGLE_DRIVE_TOKEN_PATH environment variable")
	}

	f, err := os.Open(tokenPath)
	if err != nil {
		log.Fatalf("Unable to open token file: %v", err)
	}
	defer f.Close()

	token := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(token); err != nil {
		log.Fatalf("Unable to decode token: %v", err)
	}

	return googledrive.OauthConfig.Client(context.Background(), token)
}

func handleGoogleDriveUpload(w http.ResponseWriter, r *http.Request) utils.Response {
	client := getClient()
	srv, err := drive.New(client)
	if err != nil {
		return utils.InternalServerError("Unable to create Drive client: " + err.Error())
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return utils.BadRequestError("Unable to create Drive client"+err.Error(), nil)
	}
	defer file.Close()

	driveFile := &drive.File{Name: header.Filename}
	_, err = srv.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		http.Error(w, "Unable to upload file: "+err.Error(), http.StatusInternalServerError)
		return utils.InternalServerError("Unable to upload file: " + err.Error())
	}

	return utils.SuccessResponse("File uploaded successfully", nil)
}

func GoogleDriveRoutes(r *mux.Router) {
	googleDriveRouter := r.PathPrefix("/googleDrives").Subrouter()

	// Google Drive routes
	googleDriveRouter.HandleFunc("/auth/google", handleGoogleDriveAuth).Methods(http.MethodGet)
	googleDriveRouter.HandleFunc("/auth/google/callback", utils.WrapHandler(handleGoogleDriveCallback)).Methods(http.MethodGet)
	googleDriveRouter.HandleFunc("/upload", utils.WrapHandler(handleGoogleDriveUpload)).Methods(http.MethodPost)
}
