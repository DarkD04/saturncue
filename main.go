package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Track struct {
	Filename string
	Number   int
	Type     string
}

func main() {
	// Wholesome welcoming message
	fmt.Println("Saturn Cue Maker by darkn 2025")

	// Default values
	pregap := 1
	dir := ""

	// Directory argument
	if (len(os.Args)) > 1 {
		dir = os.Args[1]
	} else {
		fmt.Println("Directory not set! Defaulting to the application directory.")
	}

	// Pregap argument
	if len(os.Args) > 2 {
		intarg, err := strconv.Atoi(os.Args[2])

		// Invalid integer
		if err != nil {
			fmt.Printf("Invalid integer argument: %s\n", os.Args[2])
			return
		}

		// Valid integer
		pregap = intarg
	} else {
		fmt.Println("Pregap not defined! it is set to 0 seconds by default")
	}

	// Free space
	fmt.Println(" ")

	// Find the iso file
	matches, err := filepath.Glob(filepath.Join(dir, "*.iso"))

	// Error if it can't look for it
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to search for ISO files: %v\n", err)
		return
	}

	// Error if there's nothing
	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "No ISO found in path: %s\n", dir)
		return
	}

	// Target the first ISO
	fmt.Fprintf(os.Stderr, "%s found.\n", filepath.Base(matches[0]))
	iso := filepath.Base(matches[0])

	// Start building cue
	var b strings.Builder
	b.WriteString(fmt.Sprintf("FILE \"%s\" BINARY\n", iso))
	b.WriteString("  TRACK 01 MODE1/2048\n")
	b.WriteString("	INDEX 01 00:00:00\n")

	// Regex to match trackNN.wav/bin (case-insensitive)
	trackRe := regexp.MustCompile(`(?i)^track0?(\d)\.([a-z0-9]+)$`)

	// Read the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Failed to read directory:", err)
		return
	}

	// Look for tracks
	var tracks []Track
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Get the track file
		match := trackRe.FindStringSubmatch(file.Name())

		// Skip the loop if it's not a correct file
		if match == nil {
			continue
		}

		// Separate track data
		var number int
		fmt.Sscanf(match[1], "%d", &number)
		ext := strings.ToLower(match[2])
		trackType := ""

		// Print what was found
		fmt.Printf(fmt.Sprintf("%s found.\n", file.Name()))

		// Assign type and file extension
		switch ext {
		case "wav":
			trackType = "WAVE"

			fmt.Println("WAV is not supported by Mednafen and SAROO")

		case "bin":
			trackType = "BINARY"

		case "raw":
			trackType = "BINARY"

		default:
			fmt.Printf(fmt.Sprintf(".%s is not a supported format\n", ext))
			return

		}

		// Push track's metadata
		tracks = append(tracks, Track{
			Filename: file.Name(),
			Number:   number + 1,
			Type:     trackType,
		})
	}

	// Sort by track number
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Number < tracks[j].Number
	})

	// Write a track to the cue sheet
	for _, t := range tracks {
		b.WriteString(fmt.Sprintf("FILE \"%s\" %s\n", t.Filename, t.Type))
		b.WriteString(fmt.Sprintf("  TRACK %02d AUDIO\n", t.Number))
		b.WriteString("	INDEX 00 00:00:00\n")
		b.WriteString(fmt.Sprintf("	INDEX 01 %02d:%02d:00\n", (pregap / 60), pregap%60))
	}

	// Removing the file extension for the output name
	ext := filepath.Ext(iso)

	// All done, write everything to a cue file
	outputPath := filepath.Join(dir, fmt.Sprintf("%s.cue", strings.TrimSuffix(iso, ext)))
	err = os.WriteFile(outputPath, []byte(b.String()), 0644)
	if err != nil {
		fmt.Println("Failed to write cue file:", err)
		return
	}

	// Yeah it's done
	fmt.Printf(fmt.Sprintf("\n%s.cue has been generated", strings.TrimSuffix(iso, ext)))
}
