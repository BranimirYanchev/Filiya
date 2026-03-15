# Filia Project Backend

Filia is a social application backend built with Go, providing a robust API for user management and social interactions.

## Overview

Filia Project Backend is a RESTful API service that powers the Filia social application. It provides endpoints for user management, authentication, posts, comments, categories, friend management, search, and notifications. The backend is designed to be scalable, secure, and maintainable.

## Features

- **User Management**: Complete user profile management, authentication, and role-based access control
- **Authentication**: 
  - Email/password authentication
  - Google OAuth integration
  - JWT token-based authentication
  - Refresh token mechanism (7-day expiration)
  - Password reset flow
  - Logout functionality
- **Posts & Comments**: Create, update, delete, and interact with posts and comments
- **Categories**: Hierarchical category system for organizing content
- **Friends**: Complete friend request system with accept/decline functionality and friends list
- **Post Recommendations**: Intelligent post recommendations based on user viewing history and categories
- **Post View Tracking**: Automatic tracking of post views for personalized recommendations
- **Search**: Full-text search for users, posts, and categories
- **Notifications**: Comprehensive notification system for friend requests, likes, comments, and more
- **Roles & Permissions**: Flexible role-based permission system
- **Security Features**:
  - Rate limiting on authentication endpoints
  - Request ID tracking for debugging
  - Centralized error handling
  - Panic recovery middleware
- **PostgreSQL**: Robust database for data persistence
- **Swagger Documentation**: Interactive API documentation available at `/api/swagger/index.html`

## Technologies Used

- **Go** (version 1.24.4) - Programming language
- **Gin** - Web framework for building the API
- **GORM** - ORM library for database operations
- **PostgreSQL** - Database for storing application data
- **godotenv** - Environment variable management
- **logrus** - Structured logging
- **bcrypt** - Password hashing
- **JWT** - JSON Web Tokens for authentication

## Installation

### Prerequisites

- Go 1.24.4 or higher
- PostgreSQL database
- Git

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/Marionvd/filia-project-backend.git
   cd filia-project-backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file in the project root with the following variables:
   ```env
   # Preferred for managed Postgres providers like Neon
   DATABASE_URL=

   # Fallback for local Postgres / Docker
   DB_USER=your_db_user
   DB_HOST=your_db_host
   DB_PASS=your_db_password
   DB_NAME=filia
   DB_PORT=5432
   DB_SSLMODE=disable
   
   JWT_SECRET=your_jwt_secret_key
   
   # Google OAuth (optional)
   GOOGLE_CLIENT_ID=your_google_client_id
   GOOGLE_CLIENT_SECRET=your_google_client_secret
   GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
   ```

4. Run the application:
   ```bash
   go run cmd/main.go
   ```

   Or use the build command:
   ```bash
   go build -o main cmd/main.go
   ./main
   ```

The server will start on port 8080. Once running, you can access:
- **API Base URL**: `http://localhost:8080/api`
- **Swagger UI**: `http://localhost:8080/api/swagger/index.html`

## API Documentation

### Interactive Documentation

The complete interactive API documentation is available via Swagger UI at:
- **Swagger UI**: `http://localhost:8080/api/swagger/index.html`

After starting the server, visit the Swagger UI to explore all endpoints, try them out, and view request/response schemas.

### Authentication

All authenticated endpoints require a JWT token in the `Authorization` header:
```
Authorization: Bearer <your-jwt-token>
```

The token is automatically set as a cookie upon successful login or registration.

### Base URL

All API endpoints are prefixed with `/api`:
```
http://localhost:8080/api
```

### API Endpoints

#### Authentication (`/api/auth`)

| Method | Endpoint | Description | Auth Required | Rate Limit |
|--------|----------|-------------|---------------|------------|
| GET | `/auth/` | Authentication home endpoint | Optional | 5/min |
| POST | `/auth/register` | Register a new user with email and password | No | 5/min |
| POST | `/auth/login` | Login with email and password | No | 5/min |
| POST | `/auth/logout` | Logout user (clears cookies) | Yes | - |
| POST | `/auth/refresh` | Refresh access token using refresh token | No | 5/min |
| POST | `/auth/reset-password-request` | Request password reset (sends reset token) | No | 5/min |
| POST | `/auth/reset-password` | Reset password using reset token | No | 5/min |
| GET | `/auth/google/login` | Initiate Google OAuth login | No | - |
| GET | `/auth/google/callback` | Handle Google OAuth callback | No | - |

**Example - Register:**
```bash
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "full_name": "John Doe",
  "password": "securePassword123",
  "repeated_password": "securePassword123"
}
```

**Example - Login:**
```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Example - Refresh Token:**
```bash
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your_refresh_token_here"
}
```

**Example - Password Reset Request:**
```bash
POST /api/auth/reset-password-request
Content-Type: application/json

{
  "email": "user@example.com"
}
```

#### Users (`/api/users`)

| Method | Endpoint | Description | Auth Required | Permissions |
|--------|----------|-------------|---------------|-------------|
| GET | `/users/` | Get all users | No | - |
| GET | `/users/public/{id}` | Get user by ID (public info) | No | - |
| GET | `/users/me` | Get current authenticated user | Yes | - |
| DELETE | `/users/me` | Delete current user account | Yes | `user:delete` |
| GET | `/users/{id}/posts` | Get posts by user | No | - |
| PUT | `/users/profile/sensitive` | Update sensitive info (email, password) | Yes | `user:update` |
| PUT | `/users/profile/basic` | Update bio | Yes | `user:update` |
| POST | `/users/profile/avatar` | Upload profile picture | Yes | `user:update` |
| PUT | `/users/{id}/role` | Change user role | Yes | `moderator:update` |

#### Posts (`/api/posts`)

| Method | Endpoint | Description | Auth Required | Permissions |
|--------|----------|-------------|---------------|-------------|
| GET | `/posts` | Get all posts (paginated, with recommendations) | No | - |
| GET | `/posts/{id}` | Get post by ID (tracks view for authenticated users) | No | - |
| GET | `/posts/{id}/comments` | Get comments for a post | No | - |
| POST | `/posts` | Create a new post | Yes | `user:create` |
| POST | `/posts/{id}/comments` | Create a comment on a post | Yes | `user:create` |
| PUT | `/posts/{id}` | Update a post | Yes | `user:update` |
| DELETE | `/posts/{id}` | Delete a post | Yes | `user:delete` |
| POST | `/posts/{id}/like` | Like/unlike a post | Yes | `user:update` |

**Example - Create Post:**
```bash
POST /api/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "My First Post",
  "content": "This is the content of my first post.",
  "category_ids": [1, 2],
  "tags": ["technology", "programming"],
  "tagged_users": [2, 3],
  "is_private": false
}
```

**Post Recommendations:**
- Authenticated users receive personalized post recommendations based on their viewing history
- Posts are prioritized by categories the user has shown interest in
- Falls back to general posts when user has no viewing history

#### Categories (`/api/categories`)

| Method | Endpoint | Description | Auth Required | Permissions |
|--------|----------|-------------|---------------|-------------|
| GET | `/categories` | Get all categories (paginated) | No | - |
| GET | `/categories/{id}` | Get category by ID | No | - |
| GET | `/categories/{id}/posts` | Get posts in category | No | - |
| GET | `/categories/name` | Find category by name | Yes | - |
| POST | `/categories` | Create a new category | Yes | `moderator:create` |
| PUT | `/categories/{id}` | Update a category | Yes | `moderator:update` |
| DELETE | `/categories/{id}` | Delete a category | Yes | `moderator:delete` |

#### Comments (`/api/comments`)

| Method | Endpoint | Description | Auth Required | Permissions |
|--------|----------|-------------|---------------|-------------|
| GET | `/comments/{id}` | Get comment by ID | No | - |
| PUT | `/comments/{id}` | Update a comment | Yes | `user:update` |
| DELETE | `/comments/{id}` | Delete a comment | Yes | `user:update` |
| POST | `/comments/{id}/like` | Like/unlike a comment | Yes | `user:update` |

#### Roles (`/api/roles`)

| Method | Endpoint | Description | Auth Required | Permissions |
|--------|----------|-------------|---------------|-------------|
| GET | `/roles` | Get all roles (paginated) | Yes | `moderator:read` |
| GET | `/roles/me` | Get current user's role | Yes | `user:read` |
| POST | `/roles` | Create a new role | Yes | `moderator:update` |
| DELETE | `/roles/{id}` | Delete a role | Yes | `admin:access` |

#### Friends (`/api/users/friends`)

| Method | Endpoint | Description | Auth Required | Permissions |
|--------|----------|-------------|---------------|-------------|
| GET | `/users/friends` | Get all friends of authenticated user | Yes | `user:update` |
| DELETE | `/users/friends/{id}` | Remove a friend | Yes | `user:update` |
| POST | `/users/friends/requests` | Send a friend request | Yes | `user:update` |
| GET | `/users/friends/requests/pending` | Get pending friend requests | Yes | `user:update` |
| GET | `/users/friends/requests/sent` | Get sent friend requests | Yes | `user:update` |
| POST | `/users/friends/requests/accept/{id}` | Accept friend request | Yes | `user:update` |
| POST | `/users/friends/requests/cancel/{id}` | Decline friend request | Yes | `user:update` |
| DELETE | `/users/friends/requests/cancel/{id}` | Delete friend request | Yes | `user:update` |

**Example - Send Friend Request:**
```bash
POST /api/users/friends/requests
Authorization: Bearer <token>
Content-Type: application/json

{
  "recipient_email": "friend@example.com"
}
```

#### Search (`/api/search`)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/search/users` | Search users by email, name, or bio | No |
| GET | `/search/posts` | Search posts by title or content | No |
| GET | `/search/categories` | Search categories by name | No |

**Example - Search Users:**
```bash
GET /api/search/users?q=john&limit=20&offset=0
```

**Example - Search Posts:**
```bash
GET /api/search/posts?q=technology&limit=20&offset=0
```

#### Notifications (`/api/notifications`)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/notifications` | Get all notifications (with pagination) | Yes |
| GET | `/notifications/unread-count` | Get count of unread notifications | Yes |
| PUT | `/notifications/read-all` | Mark all notifications as read | Yes |
| PUT | `/notifications/{id}/read` | Mark a notification as read | Yes |
| DELETE | `/notifications/{id}` | Delete a notification | Yes |

**Notification Types:**
- `friend_request` - New friend request received
- `friend_accept` - Friend request accepted
- `post_like` - Someone liked your post
- `comment` - New comment on your post
- `comment_like` - Someone liked your comment
- `post_mention` - Mentioned in a post

**Example - Get Notifications:**
```bash
GET /api/notifications?limit=20&offset=0&unread_only=false
Authorization: Bearer <token>
```

### Query Parameters

Several endpoints support pagination:
- `limit`: Number of items to return (default varies by endpoint, max: 100-200)
- `offset`: Number of items to skip (default: 0)

**Example:**
```
GET /api/posts?limit=10&offset=20
```

### Response Format

Successful responses typically follow this format:
```json
{
  "data": { ... },
  "error": ""
}
```

Error responses with request ID tracking:
```json
{
  "error": "error message",
  "request_id": "unique-request-id",
  "path": "/api/endpoint"
}
```

### Permissions

The API uses a role-based permission system. Common permissions include:

- `user:read` - Read user data
- `user:create` - Create content
- `user:update` - Update content
- `user:delete` - Delete content
- `moderator:read` - Read moderator data
- `moderator:create` - Create categories
- `moderator:update` - Update categories/users
- `moderator:delete` - Delete categories
- `admin:access` - Full admin access

### Status Codes

- `200` - Success
- `201` - Created
- `202` - Accepted
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `429` - Too Many Requests (rate limit exceeded)
- `500` - Internal Server Error

### Rate Limiting

Rate limiting is enforced on certain endpoints:

- **Authentication endpoints**: 5 requests per minute per IP address
  - `/auth/login`
  - `/auth/register`
  - `/auth/refresh`
  - `/auth/reset-password-request`
  - `/auth/reset-password`
  
- **General endpoints**: 100 requests per minute per IP address (if configured)

Rate limit headers are included in responses:
- `X-RateLimit-Limit`: Maximum number of requests allowed
- `X-RateLimit-Remaining`: Number of requests remaining
- `X-RateLimit-Reset`: Time when the rate limit resets (for exceeded limits)

## Project Structure

```
filia-project-backend/
├── cmd/
│   └── main.go           # Application entry point
├── config/
│   └── google.go         # Google OAuth configuration
├── database/
│   ├── database.go       # Database connection
│   └── seed.go           # Database seeding
├── docs/
│   ├── docs.go           # Swagger documentation (auto-generated)
│   ├── swagger.json      # Swagger JSON specification
│   └── swagger.yaml      # Swagger YAML specification
├── internal/
│   ├── handler/          # Request handlers (controllers)
│   │   ├── api.go
│   │   ├── authentication_handler.go
│   │   ├── category_handler.go
│   │   ├── comment_handler.go
│   │   ├── friend_handler.go
│   │   ├── notification_handler.go
│   │   ├── post_handler.go
│   │   ├── role_handler.go
│   │   ├── search_handler.go
│   │   └── user_handler.go
│   ├── helper/           # Helper functions
│   │   ├── api_helpers.go
│   │   ├── jwt_helpers.go
│   │   ├── role_helpers.go
│   │   └── user_helpers.go
│   ├── middleware/       # HTTP middleware
│   │   ├── error_handler.go      # Error handling & request ID tracking
│   │   ├── jwt_middleware.go     # JWT authentication
│   │   ├── permission_authorizer.go  # Permission checking
│   │   └── rate_limiter.go       # Rate limiting
│   ├── model/            # Data models and DTOs
│   │   ├── api.go
│   │   ├── category.go
│   │   ├── comment.go
│   │   ├── friends.go
│   │   ├── notification.go
│   │   ├── permission.go
│   │   ├── posts.go
│   │   ├── roles.go
│   │   └── user.go
│   └── service/          # Business logic layer
│       ├── category_service.go
│       ├── comment_service.go
│       ├── friend_service.go
│       ├── notification_service.go
│       ├── post_service.go
│       ├── role_service.go
│       ├── search_service.go
│       └── user_service.go
├── router/               # API route definitions
│   ├── api_router.go
│   ├── authentication_router.go
│   ├── category_router.go
│   ├── comment_router.go
│   ├── notification_router.go
│   ├── post_router.go
│   ├── role_router.go
│   ├── search_router.go
│   └── user_router.go
├── tests/                # Test files
│   ├── authentication_test.go
│   ├── http/
│   └── test_utils.go
├── .env                  # Environment variables (create this)
├── go.mod                # Go module definition
└── go.sum                # Go module checksums
```

## Development

### Generating Swagger Documentation

To regenerate Swagger documentation after adding new endpoints or updating annotations:

1. Install swag:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. Generate documentation:
   ```bash
   swag init -g cmd/main.go -o docs
   ```

3. Restart the server to see updates in Swagger UI.

### Testing

You can use the HTTP request examples in the `tests/http/` directory to test the API endpoints, or use the Swagger UI for interactive testing.

### Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# Database Configuration
# Preferred for managed Postgres providers like Neon
DATABASE_URL=

# Fallback for local Postgres / Docker
DB_USER=your_db_user
DB_HOST=your_db_host
DB_PASS=your_db_password
DB_NAME=filia
DB_PORT=5432
DB_SSLMODE=disable

# JWT Secret (required)
JWT_SECRET=your_secret_key_here_minimum_32_characters

# Google OAuth (optional)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/auth/google/callback
```

### Running Tests

```bash
go test ./...
```

## Features in Detail

### Post Recommendations

The backend includes an intelligent recommendation system that:
- Tracks post views for authenticated users
- Extracts preferred categories from viewing history
- Prioritizes posts in user's favorite categories
- Falls back to popular/general posts when user has no history
- Ensures variety by mixing recommended and general content

### Notifications System

The notification system provides:
- Real-time notifications for social interactions
- Automatic notifications for friend requests and accepts
- Mark as read/unread functionality
- Notification deletion
- Unread count endpoint for badge display
- Filtering by read/unread status

### Search Functionality

Search supports:
- Full-text search across users, posts, and categories
- Case-insensitive matching
- Pagination support
- For authenticated users, includes their private posts in search results

### Account Deletion

When a user deletes their account:
- User data is anonymized (not hard-deleted) to preserve data integrity
- Related data is cleaned up (friends, requests, likes)
- Posts are marked as "[Deleted]" but remain for reference
- Email is anonymized to prevent reuse

### Error Handling

The backend includes comprehensive error handling:
- Request ID tracking for debugging
- Standardized error response format
- Panic recovery middleware
- Centralized error handling
- Detailed logging with structured fields

### Rate Limiting

Rate limiting protects the API from abuse:
- IP-based rate limiting
- Configurable limits per endpoint type
- Rate limit headers in responses
- Automatic cleanup of old rate limit data

## Security Considerations

- **Password Security**: Passwords are hashed using bcrypt
- **JWT Tokens**: Secure token-based authentication with expiration
- **Refresh Tokens**: Separate refresh tokens with longer expiration
- **Rate Limiting**: Protection against brute force attacks
- **Input Validation**: All inputs are validated
- **SQL Injection Protection**: GORM provides parameterized queries
- **CORS**: Configurable CORS settings

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

Project Link: [https://github.com/Marionvd/filia-project-backend](https://github.com/Marionvd/filia-project-backend)

## Changelog

### Latest Updates

- ✅ Added logout endpoint
- ✅ Implemented refresh token mechanism
- ✅ Added password reset functionality
- ✅ Fixed friend service bugs and added friends list endpoint
- ✅ Implemented unfriend functionality
- ✅ Added post view tracking
- ✅ Completed post recommendation algorithm
- ✅ Implemented search functionality
- ✅ Created comprehensive notifications system
- ✅ Added rate limiting middleware
- ✅ Implemented centralized error handling with request ID tracking
- ✅ Added account deletion with proper data anonymization
