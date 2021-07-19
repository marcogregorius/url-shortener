# URL Shortener
This is a simple URL Shortener backend service written in Go using [gorilla/mux](https://github.com/gorilla/mux) as the HTTP router and PostgreSQL as the database.
The shortened link generation is using shortuuid.

# Features
## Create a shortlink
Generate a shortlink by calling the POST API endpoint.
```
curl 'localhost:8080/api/shortlinks' -d '{"source_url": "https://blog.golang.org/"}' -H "Content-Type: application/json" -s | jq
{
  "created_at": "2021-07-19T16:58:42.806091+08:00",
  "id": "XMRP2TAd",
  "last_visited_at": null,
  "source_url": "https://blog.golang.org/",
  "visited": 0
}
```
The `id` field denotes the id for the shortened URL.

## Visit a shortlink
Go to your browser and visit http://localhost:8080/XMRP2TAd
You will be redirected to the original URL. Each visit will increase the `visited` field by 1, and updates the `last_visited_at` timestamp.

##  Get a shortlink object
Retrieves the shortlink object.
```
curl 'http://localhost:8080/api/shortlinks/XMRP2TAd' -s | jq
{
  "created_at": "2021-07-19T16:58:42.806091+08:00",
  "id": "XMRP2TAd",
  "last_visited_at": "2021-07-19T17:06:51.012013+08:00",
  "source_url": "https://blog.golang.org/",
  "visited": 1
}
```


# Setup
1. Have a local PostgresSQL running
2. Set environment variables below which are required for connecting to Postgres (either with export or chaining them in the `go run .` command below):
	- `DB_HOST`
	- `DB_PORT`
	- `DB_USER`
	- `DB_NAME`
	- `DB_PASSWORD`
3. Initialize DB schema with going into `sql/` directory and run `bash setup_db.sh`
4. Run program with `go run . --port={your_desired_port}`. Change the flag `--port` as needed. Default at port :8080

# Tests
Test will also require the same environment variables as in the Setup.
Run test with `go test -v`
