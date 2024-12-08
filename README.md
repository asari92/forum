# Forum Image Upload Feature

This is a forum feature that allows registered users to create posts containing both images and text. When viewing a post, users (both registered and guests) can see the image associated with it. The image upload supports JPEG, PNG, and GIF formats with a maximum file size of 20MB. 

## Features:
- Image upload support for JPEG, PNG, and GIF formats.
- Automatic rejection of images exceeding 20MB in size.
- Error handling for unsupported formats or oversized images.
- Display the uploaded image alongside the post for both registered users and guests.

## Technologies Used:
- Go (Golang) - Backend logic and image handling.
- SQLite3 - Database for storing user data and post information.
- bcrypt - User authentication and password hashing.
- UUID - Unique identification for users, posts, and images.

## File Upload:
The uploaded images must be:
- JPEG, PNG, or GIF formats.
- No larger than 20MB.
  
If the image exceeds 20MB, the backend will return an error message stating "The image is too large."

## Error Handling:
The backend will handle errors related to:
- Unsupported image file types.
- Image size exceeding 20MB.
- Any other operational errors.

## Backend Implementation:

### Dependencies:
- Go standard libraries
- `github.com/gofrs/uuid`
- `github.com/mattn/go-sqlite3`
- `golang.org/x/crypto/bcrypt`



## Installation

### Prerequisites:
1. Install Go (version 1.22 or higher).
2. Install SQLite3 and make sure it's accessible.
3. Install required Go packages:
   ```bash
   go get github.com/gofrs/uuid
   go get github.com/mattn/go-sqlite3
   go get golang.org/x/crypto/bcrypt


## 👨‍💻 Authors
- [adulmaev](https://01.alem.school/git/adulmaev)
- [dkurmant](https://01.alem.school/git/dkurmant)
