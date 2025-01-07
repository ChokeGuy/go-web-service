package googledrive

import (
	"io/ioutil"
	"log"
	"web-service/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

var (
	OauthConfig *oauth2.Config
)

func Init() {
	credentialsPath := config.Env.GOOGLE_DRIVE_CREDENTIALS_PATH
	if credentialsPath == "" {
		log.Fatal("Missing GOOGLE_DRIVE_CREDENTIALS_PATH environment variable")
	}

	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	OauthConfig, err = google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
}
