# Blog Aggregator CLI

A command-line RSS feed aggregator that allows users to manage and follow RSS feeds, storing posts in a PostgreSQL database.

## Prerequisites

- Go 1.23.6 or later
- PostgreSQL database server
- [goose](https://github.com/pressly/goose) for database migrations

## Installation

```bash
go install github.com/isotronic/blog-aggregator@latest
```

## Setup

1. Create a PostgreSQL database for the application

2. Create a config file at `~/.gatorconfig.json` with the following structure:

```json
{
  "db_url": "postgres://username:password@localhost:5432/dbname?sslmode=disable",
  "current_user_name": ""
}
```

3. Run database migrations:

```bash
goose -dir sql/schema postgres "your-db-url" up
```

## Usage

### User Management

- `register <username>` - Create a new user
- `login <username>` - Log in as an existing user
- `users` - List all users
- `reset` - Delete all users from the database

### Feed Management

- `addfeed <name> <url>` - Add a new RSS feed (requires login)
- `feeds` - List all feeds
- `follow <feed-url>` - Follow an existing feed (requires login)
- `following` - List feeds you're following (requires login)
- `unfollow <feed-url>` - Unfollow a feed (requires login)

### Post Management

- `browse [limit] [offset]` - View posts from followed feeds (requires login)
  - `limit` (optional): Number of posts to show (default: 2)
  - `offset` (optional): Number of posts to skip

### Feed Aggregation

- `agg <duration>` - Start the feed aggregator that fetches new posts
  - `duration`: Time between fetches (e.g., "1h", "30m", "1m")

## Examples

```bash
# Create a new user
blog-aggregator register johndoe

# Log in
blog-aggregator login johndoe

# Add a new feed
blog-aggregator addfeed "Golang Blog" "https://go.dev/blog/feed.atom"

# Start aggregator to fetch posts every hour
blog-aggregator agg 1h

# View latest 5 posts
blog-aggregator browse 5
```
