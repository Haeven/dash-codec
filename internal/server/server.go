package server

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/Haeven/codec/pkg/codec"
	"github.com/Haeven/codec/pkg/kafka"
)

type Server struct {
	kafkaClient *kafka.KafkaClient
}

func NewServer(kafkaClient *kafka.KafkaClient) *Server {
	log.Println("Creating new server instance")
	return &Server{
		kafkaClient: kafkaClient,
	}
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Starting server")
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, shutting down server")
			return nil
		default:
			log.Println("Waiting for new message")
			m, err := s.kafkaClient.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}
			if m == nil {
				log.Println("No message received, continuing")
				continue
			}

			log.Println("Message received, processing")
			var videoUpload struct {
				VideoPath string `json:"video_path"`
			}
			if err := json.Unmarshal(m.Value, &videoUpload); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}

			log.Printf("Processing video: %s", videoUpload.VideoPath)
			outputDir := filepath.Join(os.TempDir(), "video_processing")
			if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
				log.Printf("Error creating output directory: %v", err)
				continue
			}

			log.Println("Generating video segments")
			if err := codec.GenerateSegments(videoUpload.VideoPath, outputDir); err != nil {
				log.Printf("Error generating segments: %v", err)
				continue
			}

			log.Println("Generating MPD file")
			if err := codec.GenerateMPD(outputDir); err != nil {
				log.Printf("Error generating MPD: %v", err)
				continue
			}

			mpdPath := filepath.Join(outputDir, "output.mpd")
			log.Printf("Reading MPD file: %s", mpdPath)
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

			log.Println("Marshaling result")
			resultJSON, err := json.Marshal(result)
			if err != nil {
				log.Printf("Error marshaling result: %v", err)
				continue
			}

			log.Println("Writing result to Kafka")
			if err := s.kafkaClient.WriteMessage(ctx, nil, resultJSON); err != nil {
				log.Printf("Error writing result to Kafka: %v", err)
			}

			log.Printf("Processed video: %s", videoUpload.VideoPath)
		}
	}
}
