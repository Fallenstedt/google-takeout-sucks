# Download google takeout photos

Google Takeout allows you to generate zip files of your google data. The problem is that it will generate hundreds of zip files of you data, all with this information loosely connected

I made this script to download all of my google photos so I can back it up onto a harddrive

Before you begin, export all of your google photos to google drive.

### Setup

1. Follow the [quickstart](https://developers.google.com/drive/api/quickstart/go) to setup the google cloud project. This is needed so you can generate a token via Oauth2 with Google. Store the `credentials.json` file at root level. Skip the part of setting up a golang project. Ensure the following scopes are included `https://www.googleapis.com/auth/drive.readonly https://www.googleapis.com/auth/drive.metadata.readonly`

1. Copy the directory ID of the google takeout folder in google drive. This can be found in the URL when viewing the Takeout folder in google drive. `drive.google.com/drive/folders/{really long id....}`
   ![Image of google drive](images/drive-id.png)

1. Run the `download` script with `go run cmd/download/main.go -directoryId=abc123`. By default, `dryRun` is set to true. You should see your files appear in the console

### Logs

Errors are written to a file and stored in `tmp/error.log` when downloading the files
