# Gator CLI – RSS Blog Aggregator

A command-line RSS blog aggregator written in Go. This tool allows users to register, log in, follow RSS feeds, and browse aggregated blog posts from followed feeds. It uses PostgreSQL for persistent storage.

## Requirements

To run this project, you must have the following installed:

- Go (1.20 or newer recommended): https://golang.org/dl/
- PostgreSQL (14 or newer recommended): https://www.postgresql.org/download/

## Installation

1. Clone the repository:
```
git clone https://github.com/EdBurroughes/RSSGator.git
cd gator
```

2. Install the `gator` CLI:
```
go install
```

This will compile the CLI and install it into your `$GOPATH/bin`.

## Usage

Once configured, you can run the CLI using:

### Available Commands

- `register` - Create a new user account
- `login` - Log in to your account
- `reset` - Reset the application state
- `users` - List all registered users
- `agg` - Manually aggregate all feeds
- `addfeed` - Add a new RSS feed (requires login)
- `feeds` - List all available feeds
- `follow` - Follow a feed (requires login)
- `following` - Show feeds you are following (requires login)
- `unfollow` - Unfollow a feed (requires login)
- `browse` - View aggregated posts from feeds you're following (requires login)