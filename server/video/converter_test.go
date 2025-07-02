package video

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConvertToHLS(t *testing.T) {
	inputFile := filepath.Join("testdata", "sample.mp4")
	outputDir, err := os.MkdirTemp("", "hls_test_output_*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	// the temporary directory is cleaned up after the test
	t.Cleanup(func() {
		os.RemoveAll(outputDir)
		t.Logf("Cleaned up temporary directory: %s", outputDir)
	})

	// NOTE: This function is not implemeted yet.
	err = ConvertToHLS(inputFile, outputDir)
	if err != nil {
		t.Fatalf("expected playlist file was not found: %v", err)
	}

	playlistFile := filepath.Join(outputDir, "playlist.m3u8")
	if _, err := os.Stat(playlistFile); os.IsNotExist(err) {
		t.Errorf("expected playlist file was not found: %s", playlistFile)
	}

	segmentFiles, err := filepath.Glob(filepath.Join(outputDir, "*.ts"))
	if err != nil {
		t.Fatalf("Error while searching for segment files: %v", err)
	}
	if len(segmentFiles) == 0 {
		t.Errorf("expected at least one TS segment file, but none were found")
	}

	t.Logf("Found %d segment files.", len(segmentFiles))
}
