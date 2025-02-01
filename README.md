# mp3Persist

**mp3Persist** is a lightweight Go-based command-line tool designed to play mp3 files on macOS.

It recursively scans a given directory for mp3 files, optionally shuffles them, and persists the playback order between runs.

This means that if the application is interrupted and restarted, it resumes from where it left off.

The tool leverages the macOS `afplay` command for audio playback and supports customizable break durations between tracks.

This tool is particularly suitable for listening to lessons, especially language lessons, where maintaining the playback order and taking regular breaks can enhance the learning experience.

## What is it?

- **Persistent Playback:** Remembers the current track and playback order across sessions.
- **Directory Scanning:** Recursively searches for mp3 files in a specified directory and its subdirectories.
- **Shuffle Support:** Optionally shuffles the playback order.
- **Break Intervals:** Allows you to specify a pause duration (in minutes) between tracks.
- **Cross-Platform Build for macOS:** Easily compile the binary for both Intel and ARM-based Macs.

## How to Use It

1. **Basic Usage:**

   ```bash
   ./mp3Persist -dir="/path/to/mp3/files" -break=2 -shuffle=1
   ```

-dir specifies the directory that contains mp3 files or subdirectories with mp3 files.

-break specifies the break duration in minutes between playing each file.

-shuffle is a flag (0 or 1) that determines if the list of files should be shuffled.

2. **Resuming Playback:**

The tool saves its state in a file named playorder.txt located in the specified directory.
On subsequent runs, if this file exists and the file list remains unchanged, playback resumes from the last saved position.

## How to Build It

### Requirements

- [Go](https://golang.org/dl/) (version 1.16 or higher recommended)
- macOS environment with `afplay` installed (default on macOS)

### Building for a Single Architecture

To build the binary for your current macOS architecture, run:

`go build -o mp3Persist main.go`

### Building Universal Binaries for Intel and ARM Macs

To support both Intel (`amd64`) and ARM (`arm64`) architectures, run:

    `./build`

The resulting `mp3Persist` binary will run natively on both Intel and ARM-based Macs.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests for enhancements and bug fixes.
