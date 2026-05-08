package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/grafov/m3u8"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <playlist_url> <output_filename.ts>")
		return
	}

	playlistURL := os.Args[1]
	outputFile := os.Args[2]

	err := downloadHLS(playlistURL, outputFile)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nDownload complete!")
}

func downloadHLS(m3u8URL, outputName string) error {
	// 1. Fetch the playlist
	resp, err := http.Get(m3u8URL)
	if err != nil {
		return fmt.Errorf("failed to fetch playlist: %w", err)
	}
	defer resp.Body.Close()

	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), true)
	if err != nil {
		return fmt.Errorf("failed to decode playlist: %w", err)
	}

	// We expect a Media Playlist (the one with .ts segments)
	// If it's a Master Playlist, this logic would need to select a variant.
	if listType != m3u8.MEDIA {
		return fmt.Errorf("not a media playlist; master playlists not supported in this simple example")
	}

	medialist := p.(*m3u8.MediaPlaylist)

	// Create output file
	out, err := os.Create(outputName)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	baseURL, err := url.Parse(m3u8URL)
	if err != nil {
		return err
	}

	fmt.Printf("Downloading %d segments...\n", medalist.Count())

	for i, segment := range medalist.Segments {
		if segment == nil {
			break
		}

		// Resolve segment URL (handles relative paths)
		segmentURL, err := resolveURL(baseURL, segment.URI)
		if err != nil {
			return fmt.Errorf("failed to resolve segment URL: %w", err)
		}

		fmt.Printf("\rDownloading segment [%d/%d]", i+1, medalist.Count())

		// Download and append
		err = appendSegment(segmentURL, out)
		if err != nil {
			return fmt.Errorf("\nfailed to download segment %d: %w", i, err)
		}
	}

	return nil
}

func resolveURL(base *url.URL, uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(u).String(), nil
}

func appendSegment(segmentURL string, out io.Writer) error {
	resp, err := http.Get(segmentURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	return err
}
