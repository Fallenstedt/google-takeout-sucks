package download

import (
	"context"
	"errors"
	"fmt"

	"github.com/Fallenstedt/google-takeout-sucks/internal/util"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)


var authenticationUrl = "https://google-takeout-sucks.fallenstedt.com/login"

func NewGoogleDriveService(ctx context.Context) (*drive.Service, error) {
	shouldAuthenticate := util.AskYesNoQuestionDefaultYes(`
This tool will open your browser to authenticate with Google. You will be given an access token,
and asked to paste it into the terminal. 

Authentication is handled by an open-source service built for this project:
https://github.com/Fallenstedt/google-takeout-sucks-auth

Only minimal, read-only permissions are requested to download your files.
No files are modified, deleted, or stored.

You may revoke access at any time from your Google Account security settings.

Continue to google and authenticate? [Y/n]:
`)

	if shouldAuthenticate != 1 {
		return nil, errors.New("Not opening authentication window. Exiting")
	}

	err := util.OpenBrowser(authenticationUrl)
	if err != nil {
		fmt.Println("Unable to open browser. Visit link to authenticate and paste in code")
	}

	token, err := util.WaitForResponse(`
After signing in, you will receive an access token.

Paste the token below to continue.
`)

	if err != nil || token == "" {
		return nil, errors.New("Unable to find access token. Exiting")
	}

	oauthToken := &oauth2.Token{
		AccessToken: token,
		TokenType: "Bearer",
	}

	
	svc, err := drive.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(oauthToken)))

	return svc, err
}
