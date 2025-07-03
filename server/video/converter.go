package video

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// var for HLS convert
const (
	hlsTime = "10" // target segment length. default is 2
	hlsListSize = "0" // mazimum number of playlist entries. if 0, list file will contain all the segments.
	hlsSegmentFilename = "segment%03d.ts"
	PlaylistFilename = "playlist.m3u8"
)

// NOTE: This function is call "ffmpeg" to convert video -> HLS files.
func ConvertToHLS(inputFile, outputDir string) error {
	// Ensure input file exist
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file not found: %s", inputFile)
	}

	// Create master playlist file path
	playlistPath := filepath.Join(outputDir, PlaylistFilename)
	segmentPath := filepath.Join(outputDir, hlsSegmentFilename)

	// ffmpeg command arguments
	// reference: https://ffmpeg.org/ffmpeg-formats.html#hls-2, https://ffmpeg.org/ffmpeg.html#toc-Main-options
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-c:a", "copy",
		"-c:v", "copy",
		"-f", "hls",
		"-hls_time", hlsTime,
		"-hls_list_size", hlsListSize,
		"-hls_segment_filename", segmentPath,
		playlistPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute ffmpeg command: %w\nOutput: %s", err, string(output))
	}

	// If successful, return nil
	return nil
}
