package google_internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/sotatek-dev/hyper-automation-chatbot/internal/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Credentials struct {
	Type                string `json:"type"`
	ProjectID           string `json:"project_id"`
	PrivateKeyID        string `json:"private_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientID            string `json:"client_id"`
	AuthURI             string `json:"auth_uri"`
	TokenURI            string `json:"token_uri"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url"`
	ClientCertURL       string `json:"client_x509_cert_url"`
}

func GetSheetService(config *config.GoogleConfig) *sheets.Service {
	b, err := os.ReadFile(config.Credentials)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	var credential Credentials
	err = json.Unmarshal(b, &credential)
	if err != nil {
		log.Fatalf("Unable to parse service account file: %v", err)
	}

	// Configure OAuth2
	conf, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse service account config: %v", err)
	}
	ctx := context.Background()
	client := conf.Client(ctx)

	// Create a new Sheets service
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}
	return srv
}

func GetDriveService(config *config.GoogleConfig) *drive.Service {
	b, err := os.ReadFile(config.Credentials)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	var credential Credentials
	err = json.Unmarshal(b, &credential)
	if err != nil {
		log.Fatalf("Unable to parse service account file: %v", err)
	}

	// Configure OAuth2
	conf, err := google.JWTConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse service account config: %v", err)
	}
	ctx := context.Background()
	client := conf.Client(ctx)

	// Create a new Drive service
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Drive service: %v", err)
	}
	return srv
}

func ReadSheet(srv *sheets.Service, spreadsheetId, readRange string) {
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from spreadsheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			fmt.Printf("%s, %s\n", row[0], row[1])
		}
	}
}

func ExtractSheetIdFromUrl(url string) (string, error) {
	re := regexp.MustCompile(`spreadsheets/d/([^/]+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid Google Sheet URL")
	}
	return matches[1], nil
}
