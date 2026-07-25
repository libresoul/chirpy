# Chirpy API Reference

## Health

### `GET /api/healthz`
Checks service readiness.

**Response:**
- 200 OK if healthy

---

## Users

### `POST /api/users`
Create a new user.

- Body: `{ "email": string, "password": string }`
- Response: User object
- Errors: 400 bad input, 500 internal error

### `PUT /api/users` _(auth required)_
Update email & password.

- Bearer JWT required
- Body: `{ "email": string, "password": string }`
- Response: Updated user object

### `POST /api/login`
Login and receive tokens.

- Body: `{ "email": string, "password": string }`
- Response: User object with tokens
- Errors: 400/401 bad login

---

## Sessions & Auth

### `POST /api/refresh`
Obtain new JWT using refresh token.

- Bearer: refresh token in header
- Response: `{ "token": string }`

### `POST /api/revoke`
Revoke refresh token.

- Bearer: refresh token in header
- Response: 204 No Content

---

## Chirps

### `GET /api/chirps`
List chirps, filter by author, sort optional.

- Query params: `author_id`, `sort` (`asc` or `desc`)
- Response: Array of chirps

### `POST /api/chirps` _(auth required)_
Create a new chirp.

- Bearer JWT required
- Body: `{ "body": string }`
- Response: Chirp object
- Errors: 400/401/500

### `GET /api/chirps/{chirpID}`
Get a chirp by ID.

- Response: Chirp object
- Errors: 400/404

### `DELETE /api/chirps/{chirpID}` _(auth required)_
Delete a chirp by ID (owner only).

- Bearer JWT required
- Response: 204 No Content
- Errors: 401/403/404

---

## Polka Webhooks

### `POST /api/polka/webhooks`
Handle Polka events (e.g., user upgrade).

- Header: `Authorization: ApiKey ...`
- Body: `{ "event": string, "data": { "user_id": string } }`
- Response: 204 No Content
- Errors: 400/401/404/500

---

## Admin Endpoints

### `GET /admin/metrics`
View HTML metrics dashboard.

### `POST /admin/reset`
(Development only) Reset all users.

- Response: 200 OK or 403 if unavailable

---

## Errors

- 400 – Invalid request
- 401 – Unauthorized
- 403 – Forbidden
- 404 – Not found
- 500 – Internal server error
