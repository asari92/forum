Here's the updated README file with Docker instructions:

---

# 🌐 Web Forum Project

Welcome to the **Web Forum Project**! This project aims to create an interactive web forum that enables users to communicate, share thoughts, and engage in discussions.

## 🚀 Features

- **🗣️ User Communication**: Facilitate discussions between users through posts and comments.
- **📂 Post Categories**: Associate categories with your posts for better organization.
- **👍👎 Likes and Dislikes**: Users can like or dislike posts and comments.
- **🔍 Post Filtering**: Easily filter posts by categories, created posts, and liked posts.

## 🛠️ Technology Stack

- **Database**: [SQLite](https://www.sqlite.org/docs.html)

SQLite is an embedded database choice widely used for local storage in applications. It allows for efficient database management through SQL queries.

## 📊 Database Structure

To structure your database effectively, consider creating an **Entity Relationship Diagram (ERD)**. Make sure to implement at least one `SELECT`, one `CREATE`, and one `INSERT` query.

## 🔒 Authentication

Users must register to access the forum. The registration process includes:

1. **Email Input**: 
   - Must ask for an email.
   - If the email is already taken, return an error response.
   
2. **Password Input**: 
   - Must ask for a password.
   - The password must be encrypted when stored (bonus task).

The system verifies if the provided email exists and checks if all credentials are correct. If the password is incorrect, an error response is returned.

## 💬 User Communication

- **Registered Users**: Can create posts and comments.
- **Category Association**: Users can associate one or more categories with their posts.
- **Visibility**: Posts and comments are visible to all users (registered or not). Non-registered users can only view posts and comments.

## 👍🏿👎🏿 Likes and Dislikes

- **Registered Users Only**: Can like or dislike posts and comments.
- **Visibility**: The count of likes and dislikes is visible to all users.

## 🔍 Filtering Mechanism

Implement a filtering mechanism that allows users to filter displayed posts by:

- **Categories**
- **Created Posts**
- **Liked Posts**

*Note*: The last two filters are only available for registered users and refer to the logged-in user.

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

   The server will be available at [http://localhost:4000](http://localhost:4000).

3. **Stop the Docker Container**:
   If you need to stop the running container, use:
   ```bash
   make stop
   ```

---

## 📦 Running the Server Locally (Without Docker)

1. **Install Go**: Ensure you have Go installed on your machine. You can download it from the [official Go website](https://golang.org/dl/).

2. **Clone the Repository**: Use Git to clone the project repository to your local machine.
   ```bash
   git clone <repository-url>
   cd forum
   ```

3. **Run the Server**:
   Execute the following command to start the server locally:
   ```bash
   go run ./cmd/web/main.go
   ```

   The server will be available at [http://localhost:4000](http://localhost:4000).

---

### 📑 Additional Notes

- Ensure `forum.db` is properly initialized if running locally. You can use the `make initDB` command.
- For generating certificates locally, use the `make generateCerts` command.

Enjoy building and improving your Web Forum Project! 🚀