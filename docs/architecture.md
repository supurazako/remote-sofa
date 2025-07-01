# Architecture Decision for MVP

This document outlines the technology stack and architectural decisions for the Minimum Viable Product (MVP) of this project. The primary goals for the MVP are development speed and laying a solid foundation for future scalability.

## 1. Technology Stack

| Category | Technology / Service | Rationale |
| --------------------- | ---------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| Backend Language  | Go | High performance, strong concurrency support, and excellent compatibility with the `tusd` server for resumable file uploads.  |
| Frontend Framework| React | Robust ecosystem, component-based architecture, and extensive community support. |
| HLS Video Playback| `hls.js` | De-facto standard for HLS playback in browsers, lightweight, and well-documented for integration with React. |
| Database | PostgreSQL | A relational database is suitable for managing session and user data. PostgreSQL is a reliable and feature-rich choice. |
| Deployment | Docker + AWS Fargate/ECS | Containerization provides a consistent environment. Fargate reduces infrastructure management overhead and enables easy scaling. |

## 2. Core Feature: HLS Streaming Workflow

This section details the architecture for the core video upload and streaming feature.

### 2.1. HLS Conversion Strategy

- Decision: Perform HLS conversion directly on the application server using FFmpeg.
- Rationale:
  - Simplicity for MVP: This approach avoids the complexity and cost of integrating a managed service like AWS Elemental MediaConvert at this early stage.
  - Development Efficiency: Allows for easy local testing and development, as the production environment setup is mirrored locally.
  - Future Flexibility: If performance becomes a bottleneck post-MVP, this component can be swapped out for a managed service or a dedicated microservice without altering the core application logic significantly.
- Alternatives Considered:
  - AWS Elemental MediaConvert: A powerful managed service, but deemed overly complex and costly for the MVP. It remains a primary candidate for future scaling.

### 2.2. Data Flow

1. Upload: The host user uploads a video file. The `tusd` server, running alongside the Go backend, handles the resumable upload and stores the original file directly in an AWS S3 bucket (e.g., `uploads/`).
2. Trigger: Upon successful upload, `tusd` sends a post-finish hook to the Go backend API.
3. Conversion: The Go backend receives the hook, identifies the uploaded file, and initiates a background job. This job executes an `ffmpeg` command to convert the video file into HLS format (`.m3u8` playlist and `.ts` segments).
4. Storage: The resulting HLS files are saved to a different location within the same AWS S3 bucket (e.g., `streams/<session-id>/`).
5. Playback: The client application requests the stream URL from the Go backend. The backend provides the URL to the `.m3u8` file in S3. The client's video player (`hls.js`) uses this URL to begin streaming.
6. Cleanup: A separate scheduled task or event will be responsible for deleting the original and converted files from S3 after a session ends, to manage storage costs and protect data (as per non-functional requirements).
