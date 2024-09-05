package codec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	mpd "github.com/Haeven/codec/pkg/mpd"
)

// GenerateSegments generates video segments encoded with VP9 using Av1an.
func GenerateSegments(videoFile, outputDir string) error {
	baseName := filepath.Base(videoFile)

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	resolutions := []string{"144p", "240p", "720p", "1080p", "1440p", "2160p"}

	for _, resolution := range resolutions {
		bitrate := CalculateVP9Bitrate(resolution)
		segmentPattern := filepath.Join(outputDir, fmt.Sprintf("%s_%s_segment_%%03d.webm", baseName, resolution))

		cmd := exec.Command("av1an", "-i", videoFile, "-o", outputDir, "-c", "vp9", "-b", bitrate, "-p", segmentPattern)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error encoding video segments: %w", err)
		}
	}

	return nil
}

// CalculateVP9Bitrate calculates the ideal bitrate for VP9 encoding based on resolution.
func CalculateVP9Bitrate(resolution string) string {
	switch resolution {
	case "144p":
		return "150k"
	case "240p":
		return "300k"
	case "720p":
		return "1500k"
	case "1080p":
		return "3000k"
	case "1440p":
		return "6000k"
	case "2160p":
		return "12000k"
	default:
		return "1500k"
	}
}

// GenerateMPD generates an MPD file for the video segments
func GenerateMPD(outputDir string) error {
	return mpd.GenerateMPD(outputDir)
}
