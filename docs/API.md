# Filia Project API Documentation

This document provides detailed information about the Filia Project API endpoints, request/response formats, and error handling.

## Base URL

All API endpoints are prefixed with `/api`.

```
http://localhost:8080/api
```

## Authentication

Currently, the API does not implement authentication tokens. User identification is done by providing email or username in the request body.

## Endpoints

### Home

Returns all users in the system.

**URL**: `/`  
**Method**: `GET`  
**Auth required**: No

#### Success Response

**Code**: `200 OK`  
**Content example**:

```json
{
  "users": [
    {
      "id": 1,
      "username": "user1",
      "full_name": "User One",
      "email": "user1@example.com",
      "google_id": 0,
      "password": "[hashed password]",
      "created_at": "2023-07-01T12:00:00Z"
    },
    {
      "id": 2,
      "username": "user2",
      "full_name": "User Two",
      "email": "user2@example.com",
      "google_id": 0,
      "password": "[hashed password]",
      "created_at": "2023-07-02T12:00:00Z"
    }
  ]
}
```

### Get User

Retrieves a specific user by email or username.

**URL**: `/user/`  
**Method**: `GET`  
**Auth required**: No  
**Data constraints**:

```json
{
  "email": "[valid email address]",
  "username": "[username string]"
}
```

**Note**: At least one of email or username must be provided.

#### Success Response

**Code**: `200 OK`  
**Content example**:

```json
{
  "user": {
    "id": 1,
    "username": "user1",
    "full_name": "User One",
    "email": "user1@example.com",
    "google_id": 0,
    "password": "[hashed password]",
    "created_at": "2023-07-01T12:00:00Z"
  }
}
```

#### Error Response

**Condition**: If user does not exist.  
**Code**: `404 NOT FOUND`  
**Content**:

```json
{
  "msg": "user not found"
}
```

**Condition**: If request body is malformed.  
**Code**: `400 BAD REQUEST`  
**Content**:

```json
{
  "error": "[error details]",
  "msg": "Error while binding json to object"
}
```

### Get All Users

Retrieves all users in the system.

**URL**: `/users/`  
**Method**: `GET`  
**Auth required**: No

#### Success Response

**Code**: `200 OK`  
**Content example**:

```json
{
  "users": [
    {
      "id": 1,
      "username": "user1",
      "full_name": "User One",
      "email": "user1@example.com",
      "google_id": 0,
      "password": "[hashed password]",
      "created_at": "2023-07-01T12:00:00Z"
    },
    {
      "id": 2,
      "username": "user2",
      "full_name": "User Two",
      "email": "user2@example.com",
      "google_id": 0,
      "password": "[hashed password]",
      "created_at": "2023-07-02T12:00:00Z"
    }
  ]
}
```

### Create User

Creates a new user in the system.

**URL**: `/users/new/`  
**Method**: `POST`  
**Auth required**: No  
**Data constraints**:

```json
{
  "email": "[valid email address]",
  "username": "[username string]",
  "password": "[password string]"
}
```

**Note**: Email is required. Username is optional but recommended.

#### Success Response

**Code**: `200 OK`  
**Content**: No content is returned on successful creation.

#### Error Response

**Condition**: If email is missing.  
**Code**: `400 BAD REQUEST`  
**Content**:

```json
{
  "msg": "email is required"
}
```

**Condition**: If user with email already exists.  
**Code**: `404 NOT FOUND`  
**Content**:

```json
{
  "msg": "user with this email already exists"
}
```

**Condition**: If user with username already exists.  
**Code**: `404 NOT FOUND`  
**Content**:

```json
{
  "msg": "user with this username already exists"
}
```

**Condition**: If request body is malformed.  
**Code**: `400 BAD REQUEST`  
**Content**:

```json
{
  "error": "[error details]",
  "msg": "Error while binding json to object"
}
```

**Condition**: If password hashing fails.  
**Code**: `400 BAD REQUEST`  
**Content**:

```json
{
  "error": "[error details]",
  "msg": "Error while hashing password"
}
```

## Data Models

### User

```json
{
  "id": "integer",
  "username": "string",
  "full_name": "string",
  "email": "string",
  "google_id": "integer",
  "password": "string (hashed)",
  "created_at": "datetime"
}
```

## Error Handling

The API returns appropriate HTTP status codes and error messages in JSON format:

```json
{
  "error": "Error details (if applicable)",
  "msg": "Human-readable error message"
}
```

## Future Enhancements

1. Token-based authentication
2. User profile management
3. Social features (friends, posts, comments)
4. Media upload capabilities