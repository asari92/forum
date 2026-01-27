# 🌐 Web Forum Project

Welcome to the **Web Forum Project**! This project aims to create an interactive web forum that enables users to communicate, share thoughts, and engage in discussions.

Project docs and implementation notes: https://zread.ai/asari92/forum

## 🚀 Features

- **🗣️ User Communication**: Facilitate discussions between users through posts and comments.
- **📂 Post Categories**: Associate categories with your posts for better organization.
- **👍👎 Likes and Dislikes**: Users can like or dislike posts and comments.
- **🔍 Post Filtering**: Easily filter posts by categories, created posts, and liked posts.

---

## 📦 Running the Server Locally (Without Docker)

1. **Install Go**: Ensure you have Go(v.1.22) installed on your machine. You can download it from the [official Go website](https://golang.org/dl/).

2. **Clone the Repository**: Use Git to clone the project repository to your local machine.
   ```bash
   git clone https://github.com/asari92/forum
   cd forum
   ```
3. **Initialize project**:
   Execute the following command to initialize database and create TLS certificates:
   ```bash
   make init
   ```

4. **Run the Server**:
   Execute the following command to start the server locally:
   ```bash
   go run ./cmd/web/main.go
   ```

   The server will be available at [https://localhost:4000](https://localhost:4000).

---

## 🐳 Docker Integration

This project is Dockerized for easy deployment. Follow the steps below to run the server using Docker.

### Running the Server with Docker

1. **Build the Docker Image**:
   Run the following command to build the Docker image:
   ```bash
   make build
   ```

2. **Run the Docker Container**:
   Start the server inside a Docker container with:
   ```bash
   make run
   ```

   The server will be available at [https://localhost:4000](https://localhost:4000).

3. **Stop the Docker Container**:
   If you need to stop the running container, use:
   ```bash
   make stop
   ```
