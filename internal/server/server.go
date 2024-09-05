package server

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"your-project/pkg/codec"
	"your-project/pkg/kafka"
)

type Server struct {
	kafkaClient *kafka.KafkaClient
}

func NewServer(kafkaClient *kafka.KafkaClient) *Server {
	return &Server{
		kafkaClient: kafkaClient,
	}
}

func (s *Server) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			m, err := s.kafkaClient.Reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			var videoUpload struct {
				VideoPath string `json:"video_path"`
			}
			if err := json.Unmarshal(m.Value, &videoUpload); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			outputDir := filepath.Join(os.TempDir(), "video_processing")
			if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
				log.Printf("Error creating output directory: %v", err)
				continue
			}

			if err := codec.GenerateSegments(videoUpload.VideoPath, outputDir); err != nil {
				log.Printf("Error generating segments: %v", err)
				continue
			}

			if err := codec.GenerateMPD(outputDir); err != nil {
				log.Printf("Error generating MPD: %v", err)
				continue
			}

			mpdPath := filepath.Join(outputDir, "output.mpd")
			mpdContent, err := os.ReadFile(mpdPath)
			if err != nil {
				log.Printf("Error reading MPD file: %v", err)
				continue
			}

			result := struct {
				OriginalVideo string `json:"original_video"`
				MPDContent    string `json:"mpd_content"`
			}{
				OriginalVideo: videoUpload.VideoPath,
				MPDContent:    string(mpdContent),
			}

			resultJSON, err := json.Marshal(result)
			if err != nil {
				log.Printf("Error marshaling result: %v", err)
				continue
			}

			if err := s.kafkaClient.Writer.WriteMessages(ctx, kafka.Message{
				Value: resultJSON,
			}); err != nil {
				log.Printf("Error writing result to Kafka: %v", err)
			}

			log.Printf("Processed video: %s", videoUpload.VideoPath)
		}
	}
}
