# API

## 🧰 Usage

### PATCH /v1/scheduler/intervals

- Sets user's label based on his onboarding prefferences.

**Request**

```json
{
  "hours": 9,
  "minutes": 20,
  "format": "am"
}
```

**Response**

Sample `200` Response:

```text
OK
```

Sample `4XX` or `5XX` Response:

```text
Error updating user labels
```


### POST /v1/calendars

- Sets user's label based on his onboarding prefferences.

**Request**

```json
{
  "url": "http://..../.../calendar.ics",
  "name": "Test",
  "color": "mint"
}
```

**Response**

Sample `200` Response:

```text
OK
```

Sample `4XX` or `5XX` Response:

```text
Calendar URL is not valid
```

## ⚙️ Configuration

| Setting           | Value         |
| ----------------- | ------------- |
| Runtime           | Go (1.22)     |
| Entrypoint        | `main.go`     |
| Timeout (Seconds) | 15            |
| Scopes            | `users.write` |

## 🔒 Environment Variables

No environment variables required.
