# Filia Project Frontend

React frontend for Filia.

## Local Setup

1. Install dependencies:
```bash
npm install
```

2. Create `.env` from `.env.example`:
```bash
cp .env.example .env
```

3. Set the backend API base URL:
```env
REACT_APP_API_BASE_URL=http://localhost:8080/api
```

4. Start the app:
```bash
npm start
```

## Render Setup

Frontend Render environment:
```env
REACT_APP_API_BASE_URL=https://<your-backend-service>.onrender.com/api
```

Backend Render environment:
```env
FRONTEND_URL=https://<your-frontend-service>.onrender.com
CORS_ALLOWED_ORIGINS=https://<your-frontend-service>.onrender.com
COOKIE_SECURE=true
COOKIE_SAME_SITE=none
GOOGLE_REDIRECT_URL=https://<your-backend-service>.onrender.com/api/auth/google/callback
FRONTEND_AUTH_SUCCESS_URL=https://<your-frontend-service>.onrender.com/auth/callback
```

Google OAuth redirect URI:
```text
https://<your-backend-service>.onrender.com/api/auth/google/callback
```

Render static site rewrite rule:
- Source: `/*`
- Destination: `/index.html`
- Action: `Rewrite`

This repo also includes [render.yaml](/Users/branimiryanchev/Documents/GitHub/filia-project-frontend/render.yaml) with the same SPA rewrite.

Important:
- If your Render site was created manually in the dashboard and is not managed by Blueprint sync, adding `render.yaml` to the repo does not automatically update the existing service.
- In that case you must add the rewrite rule manually in the Render Dashboard under `Static Site > Redirects/Rewrites`.

Without that rewrite, direct visits or refreshes on `/feed`, `/profile`, or `/auth/callback` will return `404 Not Found`.

## Scripts

- `npm start`
- `npm run build`
- `npm test`
