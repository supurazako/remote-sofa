package video

import (
	"fmt"
	"os/exec"
	"path/filepath"
)

// NOTE: This function is call "ffmpeg" to convert video -> HLS files.
func ConvertToHLS(inputFile, outputDir string) error {
	// Create master playlist file path
	playlistPath := filepath.Join(outputDir, "playlist.m3u8")

	// ffmpeg command arguments
	// reference: https://ffmpeg.org/ffmpeg-formats.html#hls-2, https://ffmpeg.org/ffmpeg.html#toc-Main-options
	// -i: input file
	// -c:a copy -c:v copy: audio and video codecs to copy without re-encoding
	// -f hls: specify the output format as HLS
	// -hls_time [duration]: target segment length. default is 2
	// -hls_list_size [size]: mazimum number of playlist entries. if 0, list file will contain all the segments.
	// -hls_segment_filename [filename]: segment filename.
	cmd := exec.Command("ffmpeg",
		"-i", inputFile,
		"-c:a", "copy",
		"-c:v", "copy",
		"-f", "hls",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-hls_segment_filename", filepath.Join(outputDir, "segment%03d.ts"),
		playlistPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute ffmpeg command: %w\nOutput: %s", err, string(output))
	}

	// If successful, return nil
	return nil
}
