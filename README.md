# RSS Aggregator

This project aggregates RSS feeds from various URLs and saves them in a posts database. It provides an API to create users, authenticate them on the system, and allow them to create posts, follow, and unfollow posts.

## Key Features

- **User Management:** Create and authenticate users, enabling them to interact with posts.
- **Post Management:** Create, follow, and unfollow posts.
- **RSS Feed Aggregation:** Utilize Go routines for multithreaded processing when fetching RSS feeds.
- **Database Management:** Use Goose for database migration and SQLC for query generation.

## Technologies Used

- **Go Routines:** For efficient multithreading during RSS feed fetching.
- **Goose:** To handle database migrations seamlessly.
- **SQLC:** For generating type-safe SQL queries.

This project ensures robust and efficient management of RSS feeds and user interactions.
