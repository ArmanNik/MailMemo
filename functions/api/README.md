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

- Add new calendar for an user

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

### DELETE /v1/subscriptions

- Mark account unsubscribed from all emails

> This endpoint is public and does not require user session

**Request**

```json
{
  "email": "matej@appwrite.io"
}
```

**Response**

Sample `200` Response:

```text
OK
```

Sample `4XX` or `5XX` Response:

```text
Could not mark as unsubscribed
```

## ⚙️ Configuration

All settings of a function can be found in `appwrite.json` file in root of this repository.

## 🔒 Environment Variables

No environment variables required.
