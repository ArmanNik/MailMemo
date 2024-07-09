# syncCalendar

## 🧰 Usage

### POST /

- Downloads calendar from URL and synchronize `events` collection with events from calendar

Request body must be `$id` of document from `calendars` collection.

**Response**

Sample `200` Response:

```text
OK
```

## ⚙️ Configuration

| Setting           | Value         |
| ----------------- | ------------- |
| Runtime           | Go (1.22)     |
| Entrypoint        | `main.go`     |
| Timeout (Seconds) | 900           |
| Scopes            | `databases.write`, `databases.read` |

## 🔒 Environment Variables

No environment variables required.
