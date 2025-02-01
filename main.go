package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const orderFileName = "playorder.txt"

func main() {
	// Parse command-line arguments
	dirPath := flag.String("dir", "", "Path to the directory containing mp3 files or subdirectories")
	breakMinutes := flag.Int("break", 0, "Duration of break between files (in minutes)")
	shuffleFlag := flag.Int("shuffle", 0, "Shuffle flag (0 or 1)")
	flag.Parse()

	if *dirPath == "" {
		log.Fatal("Please specify the directory path using -dir")
	}

	absDir, err := filepath.Abs(*dirPath)
	if err != nil {
		log.Fatalf("Error obtaining absolute path: %v", err)
	}

	// Recursively collect mp3 files
	var files []string
	err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".mp3" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking the directory: %v", err)
	}
	if len(files) == 0 {
		log.Fatal("No mp3 files found.")
	}

	fmt.Printf("Found %d mp3 files in %s", len(files), absDir)

	// Path for the playback order file
	orderFilePath := filepath.Join(absDir, orderFileName)

	var playOrder []string
	var currentIndex int
	var currentMode int

	playOrder = nil
	// Attempt to load existing playback order
	if fileExists(orderFilePath) {
		playOrder, currentIndex, currentMode, err = loadOrder(orderFilePath)
		if err != nil {
			log.Printf("Error reading playback order, a new order will be created: %v", err)
			playOrder = nil
		}
		if currentMode != *shuffleFlag {
			log.Printf("Mode has changed, a new order will be created")
			playOrder = nil
		}
	}

	// If no order was loaded or file list has changed, create a new order
	if playOrder == nil || len(playOrder) != len(files) {
		playOrder = files
		if *shuffleFlag == 1 {
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(playOrder), func(i, j int) {
				playOrder[i], playOrder[j] = playOrder[j], playOrder[i]
			})
		}
		currentIndex = 0
		if err := saveOrder(orderFilePath, playOrder, currentIndex, *shuffleFlag); err != nil {
			log.Fatalf("Error saving playback order: %v", err)
		}
	}

	// Playback loop
	for {
		// Restart from beginning if at end of list
		if currentIndex >= len(playOrder) {
			if *shuffleFlag == 1 {
				rand.Seed(time.Now().UnixNano())
				rand.Shuffle(len(playOrder), func(i, j int) {
					playOrder[i], playOrder[j] = playOrder[j], playOrder[i]
				})
			}
			currentIndex = 0
		}

		currentFile := playOrder[currentIndex]
		fmt.Printf("Playing: %s\n", currentFile)

		// Use afplay to play the mp3 file (available on macOS)
		cmd := exec.Command("afplay", currentFile)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("Error playing file %s: %v", currentFile, err)
		}

		// Update current index and save state
		currentIndex++
		if err := saveOrder(orderFilePath, playOrder, currentIndex, *shuffleFlag); err != nil {
			log.Printf("Error saving playback order: %v", err)
		}

		// Wait for the specified break duration
		if *breakMinutes > 0 {
			fmt.Printf("Break for %d minute(s)...\n", *breakMinutes)
			time.Sleep(time.Duration(*breakMinutes) * time.Minute)
		}
	}
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// loadOrder loads the playback order from a file.
// File format: first line is the current index, second line is the current mode, following lines are file paths.
func loadOrder(path string) (order []string, currentIndex int, currentMode int, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		if lineNum == 0 {
			// First line is the index
			currentIndex, err = strconv.Atoi(line)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("error converting index: %v", err)
			}
		} else if lineNum == 1 {
			// Second line is the index
			currentMode, err = strconv.Atoi(line)
			if err != nil {
				return nil, 0, 0, fmt.Errorf("error converting mode: %v", err)
			}
		} else {
			order = append(order, line)
		}
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		return nil, 0, 0, err
	}
	return order, currentIndex, currentMode, nil
}

// saveOrder saves the playback order and current index to a file.
func saveOrder(path string, order []string, currentIndex int, currentMode int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	// First line: current index
	_, err = writer.WriteString(fmt.Sprintf("%d\n", currentIndex))
	if err != nil {
		return err
	}
	// Second line: current mode
	_, err = writer.WriteString(fmt.Sprintf("%d\n", currentMode))
	if err != nil {
		return err
	}
	for _, line := range order {
		_, err = writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
