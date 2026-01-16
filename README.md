# google-takeout-sucks

[!["Buy Me A Coffee"](https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png)](https://www.buymeacoffee.com/fallenstedt)

[Homepage](https://google-takeout-sucks.fallenstedt.com/)

A command-line tool for downloading Google Takeout exports from Google Drive.

Google Takeout often splits exports into many ZIP files. Downloading them manually through a browser is slow, error-prone, and difficult to resume. This tool automates downloading all Takeout files directly from Google Drive.

## Why This Exists

When exporting large datasets, Google Takeout:

- Splits data into many ZIP archives

- Requires manual downloads for each file

- Provides no way to resume failed downloads

- Expires download links

If you choose “Export to Google Drive” when creating a Takeout, Google places all ZIP files into a Drive folder. This tool downloads everything from that folder automatically.

## Authentication

Authentication only occurs when running the download command.

This project uses a separate, open-source authentication service:

https://github.com/Fallenstedt/google-takeout-sucks-auth

When you run the download command:

1. A browser window opens
1. You sign in with Google using the auth service
1. Google returns an OAuth access token
1. You paste the token into the terminal
1. The tool uses the token to download your files

Only read-only Google Drive access is requested.
No files are modified or deleted.

## Installation

```
go install github.com/Fallenstedt/google-takeout-sucks@latest
```

Or build from source:

```
git clone https://github.com/Fallenstedt/google-takeout-sucks.git
cd google-takeout-sucks
go build
```

## Usage

1. Create a Google Takeout Export

Go to Google Takeout

Select the services you want to export

Choose Export to Google Drive

Wait for the export to finish

Google will create a folder in your Drive containing the ZIP files.

2. Get the Google Drive Folder ID

Open Google Drive in your browser

Open the folder that contains your Takeout files

Look at the URL in the address bar

The folder ID is the long string after /folders/

Example:

```
https://drive.google.com/drive/folders/1MLFBDprqxzMp4zgN_1PPpRABdBYnHfn0
```

Folder ID:

```
1MLFBDprqxzMp4zgN_1PPpRABdBYnHfn0
```

3. Download the Files

Run:

```
google-takeout-sucks download --directoryId <FOLDER_ID> --outDir <ABSOLUTE_PATH_TO_OUTPUT_DIR>
```

The tool will:

- Open a browser to authenticate
- Prompt you to paste the access token
- Download all ZIP files from the folder

4. Unzip all the files

If you'd like to bulk unzip files, you can use the following `unzip` command

Run:

```
google-takeout-sucks unzip --source <SOURCE_DIR_CONTAINING_ALL_ZIP_FIELS> --outDir <OUTPUT_DIR_FOR_UNZIPPED_FILES>
```

## Security and Privacy

Uses OAuth 2.0 access tokens

Requests read-only Google Drive access

Does not store credentials

Does not modify or delete Google data

You can revoke access at any time from your Google Account security settings.

## Notes

This tool is intentionally simple and focused on one task: downloading Google Takeout exports without manual browser downloads.

Issues and contributions are welcome.
