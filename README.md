# HLS Downloader

A simple command-line tool written in Go that downloads an HLS (HTTP Live Streaming) media playlist and merges all `.ts` segment files into a single output file.

## Features

- Fetches and parses an HLS **media playlist** (`.m3u8`).
- Resolves relative segment URLs against the playlist URL.
- Concatenates all `.ts` segments in order into one output file.
- Shows progress during download.
- Handles HTTP errors gracefully
