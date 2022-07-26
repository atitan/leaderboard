# Leaderboard

## How to Use

### Build and run executable

```bash
$ go build
$ ./homework
```

### Send HTTP request to endpoints

```bash
$ curl -X GET 'http://localhost:8080/api/v1/leaderboard'
{"topPlayers":[]}

$ curl -X POST -H 'ClientId: Tom' -H 'Content-Type: application/json' -d '{"score":13.37}' 'http://localhost:8080/api/v1/score'
{"status":"ok"}

$ curl -X GET 'http://localhost:8080/api/v1/leaderboard'
{"topPlayers":[{"clientId":"Tom","score":13.37}]}

```

## Storage consideration

The application has the following characteristics:

1. Support multi-server design
2. Short data lifecycle
3. Data does not require strong persistence

Plus, Redis server provides:

1. Built-in sorted set algorithm
2. High performance

Thus Redis server is simple and robust choice for storage.
