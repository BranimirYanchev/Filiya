# Frontend API Handoff

## Base URL

`http://localhost:8080/api`

## Authentication

Protected endpoints expect:

`Authorization: Bearer <token>`

Important:

- `login` and `register` return both `token` and `refresh_token` in JSON
- auth is expected via `Authorization: Bearer <token>`
- `refresh_token` should be sent in the JSON body to `/auth/refresh`
- the backend does not set auth cookies unless `AUTH_COOKIES_ENABLED=true`
- some endpoints return `data`, others return direct arrays or `message`

---

## Auth

### `POST /auth/register`

Body:

```json
{
  "email": "user@example.com",
  "full_name": "John Doe",
  "password": "securePassword123",
  "repeated_password": "securePassword123"
}
```

### `POST /auth/login`

Body:

```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

### `POST /auth/refresh`

Body:

```json
{
  "refresh_token": "..."
}
```

Response:

```json
{
  "token": "...",
  "refresh_token": "..."
}
```

### `POST /auth/logout`

No body.

### `POST /auth/reset-password-request`

Body:

```json
{
  "email": "user@example.com"
}
```

### `POST /auth/reset-password`

Body:

```json
{
  "token": "reset_token",
  "new_password": "newPassword123"
}
```

### `GET /auth/google/login`

Starts Google OAuth flow.

### `GET /auth/google/callback?code=...`

Returns Google-authenticated user data plus `refresh_token`.

---

## Users

### `GET /users/`

Returns all users.

### `GET /users/public/:id`

Returns public user data.

### `GET /users/me`

Auth required.

### `PUT /users/profile/sensitive`

Auth required.

Body:

```json
{
  "email": "new@example.com",
  "full_name": "New Name",
  "password": "NewPassword123",
  "old_password": "CurrentPassword123"
}
```

### `PUT /users/profile/basic`

Auth required.

Body:

```json
{
  "bio": {
    "valid": false,
    "value": "New bio text"
  }
}
```

Important: this endpoint expects this exact nested `bio` object shape.

### `POST /users/profile/avatar`

Auth required.

`multipart/form-data`

Field:

`profile_picture=<file>`

### `DELETE /users/me`

Auth required.

### `PUT /users/:id/role`

Auth required.

Body:

```json
{
  "role_name": "RoleModerator"
}
```

### `GET /users/:id/posts`

Intended to return posts by user id.

Important: current backend implementation looks incorrect and may not use the path id properly.

---

## Posts

### `GET /posts?limit=20&offset=0`

Returns paginated posts.

### `GET /posts/:id`

Returns single post.

### `POST /posts`

Auth required.

Body:

```json
{
  "title": "Post title",
  "content": "Post content",
  "category_ids": [1, 2],
  "tags": ["tag1", "tag2"],
  "tagged_users": [2, 3],
  "is_private": false
}
```

### `PUT /posts/:id`

Auth required.

Body:

```json
{
  "title": "Updated title",
  "content": "Updated content",
  "category_ids": [1, 2],
  "tags": ["tag1", "tag2"],
  "tagged_users": [2, 3],
  "is_private": false
}
```

### `DELETE /posts/:id`

Auth required.

### `POST /posts/:id/like`

Auth required.

No body.

### `GET /posts/:id/comments`

Returns comments for a post.

### `POST /posts/:id/comments`

Auth required.

Body:

```json
{
  "content": "Great post!"
}
```

---

## Comments

### `GET /comments/:id`

Returns single comment.

### `PUT /comments/:id`

Auth required.

Body:

```json
{
  "content": "Updated comment"
}
```

### `DELETE /comments/:id`

Auth required.

### `POST /comments/:id/like`

Auth required.

No body.

---

## Categories

### `GET /categories?limit=10&offset=0`

Returns paginated categories.

Important: this endpoint currently returns HTTP `201 Created` even though it is a `GET`.

### `GET /categories/:id`

Returns single category.

### `GET /categories/:id/posts?limit=10&offset=0`

Returns posts in category.

### `GET /categories/name`

Auth required.

Body:

```json
{
  "name": "Technology"
}
```

Important: this is implemented as `GET` with JSON body, which is non-standard for frontend clients.

### `POST /categories`

Auth required.

Body:

```json
{
  "name": "Technology",
  "parent_id": 1
}
```

### `PUT /categories/:id`

Auth required.

Body:

```json
{
  "name": "Updated category",
  "parent_id": 1
}
```

### `DELETE /categories/:id`

Auth required.

---

## Search

### `GET /search/users?q=john&limit=20&offset=0`

### `GET /search/posts?q=history&limit=20&offset=0`

### `GET /search/categories?q=art&limit=20&offset=0`

All return search results in `data`.

---

## Friends

### `POST /users/friends/requests`

Auth required.

Body:

```json
{
  "recipient_email": "friend@example.com"
}
```

### `GET /users/friends/requests/pending`

Auth required.

### `GET /users/friends/requests/sent`

Auth required.

### `POST /users/friends/requests/accept/:id`

Auth required.

### `POST /users/friends/requests/cancel/:id`

Auth required.

This is the decline endpoint.

### `DELETE /users/friends/requests/cancel/:id`

Auth required.

Deletes the friend request.

### `GET /users/friends`

Auth required.

### `DELETE /users/friends/:id`

Auth required.

---

## Notifications

### `GET /notifications/?limit=20&offset=0&unread_only=true`

Auth required.

### `GET /notifications/unread-count`

Auth required.

### `PUT /notifications/read-all`

Auth required.

### `PUT /notifications/:id/read`

Auth required.

### `DELETE /notifications/:id`

Auth required.

Notification example:

```json
{
  "id": 1,
  "user_id": 2,
  "type": "friend_request",
  "title": "New friend request",
  "message": "John sent you a request",
  "read": false,
  "entity_id": 10,
  "entity_type": "friend_request",
  "created_at": "2026-03-10T10:00:00Z"
}
```

---

## Roles

### `GET /roles?limit=20&offset=0`

Auth required.

### `GET /roles/me`

Auth required.

### `POST /roles`

Auth required.

Body:

```json
{
  "name": "Editor",
  "permissions": ["user:read", "user:update"]
}
```

### `DELETE /roles/:id`

Auth required.

---

## Known Backend Issues

- `login/register` do not return access token directly in JSON
- protected endpoints rely on `Authorization` header
- `GET /categories/name` is implemented as GET with body
- `GET /users/:id/posts` appears buggy in current backend code
- response shape is inconsistent across endpoints
- `GET /categories` returns `201 Created` instead of `200 OK`
