# syncCalendarScheduler

## 🧰 Usage

### POST /

- Executes `syncCalendar` for each calendar.

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
| Cron              | `*/15 * * * *`|
| Timeout (Seconds) | 900           |
| Scopes            | `execution.write`, `databases.read` |

## 🔒 Environment Variables

No environment variables required.
