# API-Plan für das Frontend

## Aktuell vorhandene Routen

| Methode | Route          | Request                                        | Response                       |
|---------|----------------|------------------------------------------------|--------------------------------|
| GET     | `/helloworld`  | kein JSON                                      | Text: `Hello, world!`          |
| POST    | `/user/create` | `{ "name": "max", "password_hash": "geheim" }` | Text: `User created`           |
| POST    | `/user/login`  | aktuell kein JSON                              | setzt `session_token`-Cookie   |
| GET     | `/user/auth`   | Cookie `session_token`                         | aktuell Token als Text         |
| GET     | `/users`       | kein JSON                                      | `[{ "id": 1, "name": "max" }]` |

## Empfohlener Zielvertrag

### `POST /user/register`

Request:

```json
{ "name": "max", "password": "geheim" }
```

Response `201 Created`:

```json
{ "message": "User created", "user": { "id": 1, "name": "max" } }
```

### `POST /user/login`

Request:

```json
{ "name": "max", "password": "geheim" }
```

Response `200 OK`: `Set-Cookie: session_token=...; HttpOnly; Secure; SameSite=Lax`

```json
{ "message": "Login successful", "user": { "id": 1, "name": "max" } }
```

### `GET /user/auth`

Request: kein JSON, Cookie wird automatisch gesendet.

Response `200 OK`:

```json
{ "authenticated": true, "user": { "id": 1, "name": "max" } }
```

Ohne gültige Session: `401 Unauthorized` und `{ "error": "unauthorized" }`.

### `POST /user/logout`

Request: kein JSON. Session serverseitig löschen und Cookie mit `Max-Age=0` überschreiben.

Response:

```json
{ "message": "Logout successful" }
```

### `GET /users`

Nur für authentifizierte Benutzer. Response `200 OK`:

```json
{ "users": [{ "id": 1, "name": "max" }] }
```

Das Frontend sendet alle Requests mit `credentials: "include"`, damit Cookies mitgeschickt werden. Für lokale Entwicklung über HTTP muss `Secure: false` gesetzt werden; in Produktion immer HTTPS und `Secure: true` verwenden.
