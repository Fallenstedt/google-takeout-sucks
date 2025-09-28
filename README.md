> 🚨 This is a work in progress

# Google Takeout Sucks

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/fallenstedt)

Google Takeout allows you to generate zip files of your google data. The problem is that it will generate hundreds of zip files of your data, and it sucks downloading it if these files are 50GB.

I made this script to download all of my stuff from Google Takeout

Before you begin, export all of your google photos to google drive using [Google Takeout](https://takeout.google.com/settings/takeout/custom/photos).

### Setup

1. Follow the [quickstart](https://developers.google.com/drive/api/quickstart/go) to setup the google cloud project. This is needed so you can generate a token via Oauth2 with Google. Store the `credentials.json` file at root level. Skip the part of setting up a golang project. Ensure the following scopes are included `https://www.googleapis.com/auth/drive.readonly https://www.googleapis.com/auth/drive.metadata.readonly`

1. Copy the directory ID of the google takeout folder in google drive. This can be found in the URL when viewing the Takeout folder in google drive. `drive.google.com/drive/folders/{really long id....}`
   ![Image of google drive](images/drive-id.png)

1. Run the `download` script with `go run cmd/download/main.go -directoryId=abc123`. By default, `dryRun` is set to true. You should see your files appear in the console

### Logs

Error and info logs are written to a file and stored in `tmp/*.log` when running a script

### Downloading all of my Google Takeout files

Run the following command locally after setting up the project. When you are confident everything is setup correctly, set `dryRun` to false.

```
go run cmd/download/*.go -directoryId={GoogleTakeout Directory ID} -dryRun=true -outDir={/absolute/path/to/destination/directory}
```

### Unzipping all of my Google Takeout files

After downloading all of your files, you can unzip them all. When you are confident everything is setup correctly, set `dryRun` to false.

```
go run cmd/unzip/*.go -source=/absolute/path/to/directory/with/zipfiles -out=/absolute/path/for/unizpped/files -dryRun=true
```
