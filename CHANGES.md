# Changes Made to google-takeout-sucks

## Summary
Modified the download functionality to skip files that have already been downloaded by checking both file existence and file size. Enhanced logging to provide detailed progress information for each file.

## Changes

### File: `internal/download/files.go`

1. **Line 6**: Added `log` package import for logging support

2. **Line 18**: Updated `DownloadFileInput` struct to include logger
   - Added `infoLog *log.Logger` field to pass logger to download function

3. **Line 23**: Updated `FetchFiles()` function to retrieve file size from Google Drive
   - Changed: `Fields("nextPageToken, files(id, name)")`
   - To: `Fields("nextPageToken, files(id, name, size)")`
   - This ensures we get the file size from Google Drive API

4. **Line 49**: Updated `DownloadWorker()` function signature to accept logger
   - Added `infoLog *log.Logger` parameter
   - Passes logger to `downloadFileToDisk()` function

5. **Lines 77-128**: Enhanced `downloadFileToDisk()` function with detailed logging
   - Added file existence check using `os.Stat()`
   - Added file size comparison between local and remote files
   - Skips download if file exists and size matches
   - Re-downloads if file exists but size doesn't match (handles incomplete downloads)
   - **New logging statements:**
     - "Processing file: [filename]" - When starting to process a file
     - "File found locally: [filename] (local size: X bytes, remote size: Y bytes)" - When file exists
     - "Skipping file (already complete): [filename]" - When file is already downloaded
     - "File size mismatch, re-downloading: [filename]" - When file needs re-download
     - "New file, downloading: [filename] (size: X bytes)" - When file doesn't exist
     - "Starting download: [filename]" - When download begins
     - "Download completed: [filename]" - When download finishes

### File: `cmd/download.go`

1. **Line 90**: Updated worker initialization to pass logger
   - Changed: `go download.DownloadWorker(w, processCh, errCh, resCh, srv, cfg, &wg)`
   - To: `go download.DownloadWorker(w, processCh, errCh, resCh, srv, cfg, downloadInfoLog, &wg)`
   - Passes the existing `downloadInfoLog` to each worker

## Behavior

### Before Changes
- Always downloaded all files, even if they already existed locally
- Would overwrite existing files
- Limited logging - only showed when files were saved

### After Changes
- Checks if file exists locally before downloading
- Compares local file size with Google Drive file size
- Only downloads if:
  - File doesn't exist locally, OR
  - File exists but size doesn't match (incomplete download)
- Skips download if file exists and size matches
- **Detailed logging** shows:
  - Which file is being processed
  - Whether file was found locally
  - File size comparison (local vs remote)
  - Whether file is being skipped, re-downloaded, or is new
  - When download starts
  - When download completes

## Benefits
- Saves bandwidth by not re-downloading existing files
- Allows resuming interrupted download sessions
- Handles incomplete downloads by detecting size mismatches
- **Enhanced visibility** into download progress with detailed logging
- Easy troubleshooting with file-by-file status updates
- No changes to CLI interface - works transparently

## Example Log Output

```
Processing file: takeout-20230101T120000Z-001.zip
New file, downloading: takeout-20230101T120000Z-001.zip (size: 2147483648 bytes)
Starting download: takeout-20230101T120000Z-001.zip
Download completed: takeout-20230101T120000Z-001.zip

Processing file: takeout-20230101T120000Z-002.zip
File found locally: takeout-20230101T120000Z-002.zip (local size: 2147483648 bytes, remote size: 2147483648 bytes)
Skipping file (already complete): takeout-20230101T120000Z-002.zip

Processing file: takeout-20230101T120000Z-003.zip
File found locally: takeout-20230101T120000Z-003.zip (local size: 1048576 bytes, remote size: 2147483648 bytes)
File size mismatch, re-downloading: takeout-20230101T120000Z-003.zip
Starting download: takeout-20230101T120000Z-003.zip
Download completed: takeout-20230101T120000Z-003.zip
```

## Testing
To test the changes:

1. Build the application: `go build`
2. Run a download: `./google-takeout-sucks download --directoryId=<id> --outDir=<path>`
3. Interrupt the download partway through
4. Run the same command again - it should skip already-downloaded files and resume
5. Check logs (shown in console and in `~/.google_takeout_sucks/logs/download info.log`) - you should see detailed status for each file

## Deployment

### Building
```bash
go build
```

### Installation
```bash
go install
```

Or copy the built binary to your desired location.

### Usage (unchanged)
```bash
takeout download --directoryId=<directoryId> --outDir=<absolute/path/to/save/files>
```
