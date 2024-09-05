# Video Processing Service

## What it does

This service processes video files using VP9 encoding and generates DASH (Dynamic Adaptive Streaming over HTTP) compatible output. It:

1. Consumes video file paths from a Kafka topic
2. Encodes the video into multiple resolutions using VP9 codec
3. Generates video segments for adaptive streaming
4. Creates an MPD (Media Presentation Description) file for DASH streaming
5. Writes the processing results back to another Kafka topic

## How to run

### Using Docker

1. Ensure you have Docker installed on your system.

2. Build the Docker image:
   ```
   docker build -t video-processing-service .
   ```

3. Run the container:
   ```
   docker run -d --name video-processor video-processing-service
   ```

### Manual Setup

If you prefer to run the application without Docker:

1. Ensure you have Go installed on your system.

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Install Av1an and VP9 encoder:
   ```
   sudo apt-get install av1an libvpx-dev ffmpeg 
   ```

4. Set up Kafka:
   - Ensure Kafka is running on `localhost:9092`
   - Create two topics: `video-uploads` and `processed-videos`

5. Run the application:
   ```
   go run cmd/main.go
   ```

## Usage

To process a video, publish a message to the `video-uploads` topic with the following JSON structure:
```json
{
  "video_path": "/[path]/video.mp4"
}
```

7. The service will process the video and publish the results to the `processed-videos` topic.

Note: Adjust Kafka configuration in `cmd/main.go` if needed.
