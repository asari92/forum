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
   
2. **Username Input**: 
   - Must ask for a username.

3. **Password Input**: 
   - Must ask for a password.
   - The password must be encrypted when stored (bonus task).

The system verifies if the provided email exists and checks if all credentials are correct. If the password is incorrect, an error response is returned.

## 💬 User Communication

- **Registered Users**: Can create posts and comments.
- **Category Association**: Users can associate one or more categories with their posts.
- **Visibility**: Posts and comments are visible to all users (registered or not). Non-registered users can only view posts and comments.

## 👍👎 Likes and Dislikes

- **Registered Users Only**: Can like or dislike posts and comments.
- **Visibility**: The count of likes and dislikes is visible to all users.

## 🔍 Filtering Mechanism

Implement a filtering mechanism that allows users to filter displayed posts by:

- **Categories**
- **Created Posts**
- **Liked Posts**

*Note*: The last two filters are only available for registered users and refer to the logged-in user.

## 🐳 Docker Integration

This project must utilize Docker for deployment. Familiarize yourself with Docker basics through the ASCII Art Web Dockerize subject.

### 📝 Instructions

- Use SQLite for data storage.
- Handle website errors and HTTP status codes appropriately.
- Ensure all technical errors are managed effectively.
- Follow best coding practices.
- Include unit testing with test files.

### 📦 Allowed Packages

- All standard Go packages.
- `sqlite3`
- `bcrypt`
- `UUID`

**Important**: No frontend libraries or frameworks (e.g., React, Angular, Vue) are allowed.

## 📚 Learning Outcomes

By completing this project, you will learn about:

- Basics of web development:
  - HTML
  - HTTP
  - Sessions and cookies
- Setting up and using Docker:
  - Containerization
  - Compatibility and dependency management
  - Creating Docker images
- SQL language fundamentals
- Database manipulation techniques
- Basic encryption principles

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.

---

Thank you for checking out the **Web Forum Project**! We hope you enjoy building and learning through this experience. If you have any questions or feedback, feel free to reach out! 💬
