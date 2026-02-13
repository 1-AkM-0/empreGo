# empreGo: Automated Tech Job Aggregator


#### Video Demo: https://www.youtube.com/watch?v=mKo7JWvq5ns

## Overview

**empreGo** is a high-performance, concurrent CLI tool designed to automate the search for tech internships in Brazil.

Finding a job often requires checking multiple platforms (LinkedIn, Gupy, Indeed) manually every day. `empreGo` solves this fragmentation by aggregating these sources into a single pipeline, deduplicating results to avoid spam, and sending notifications every 2 hours to a private Discord channel.

It was built to run efficiently on low-resource environments (like a Google Cloud `e2-micro` instance) by avoiding heavy headless browsers in favor of direct HTTP requests and HTML parsing.

## Features

- **Multi-Source Aggregation**: Fetches jobs from:
  - **LinkedIn**: Uses `goquery` for server-side HTML parsing (Guest Mode).
  - **Gupy**: Reverse-engineered public API for lightweight JSON retrieval.
- **Smart Deduplication**: Uses SQLite to store job links and prevents sending the same alert twice.
- **Concurrency**: Implements Go Goroutines to scrape multiple sources simultaneously, reducing execution time by 70%.
- **Discord Integration**: Sends formatted alerts directly to a configured channel.
- **Resource Efficient**: Designed to run on <1GB RAM environments.

## Tech Stack 

- **Language**: Go (Golang) 1.25+
- **Database**: SQLite3 (via `modernc.org/sqlite` - CGo-free driver)
- **Libraries**:
  - `goquery` (HTML parsing)
  - `discordgo` (Discord API)
- **Infrastructure**: Google Cloud Platform, Cron Jobs, Linux/Bash.

## Project Structure

Here is a breakdown of the main files and directories in the project:

### `main.go`
The entry point of the application. It orchestrates the entire flow:
1. Initializes the SQLite database connection.
2. Launches concurrent scrapers (LinkedIn, Gupy) using Goroutines.
3. Iterates over the results and checks the database for duplicates.
4. Sends new unique jobs to Discord and saves them to the DB.

### `internal/search/`
Contains the logic for fetching data from external sources. The `scraper.go` file is responsible for:
- Defining the `Job` and `GupyJobs` structs.
- **SearchLinkedin()**: Implementing the scraping logic for LinkedIn using specific User-Agents and CSS selectors to bypass basic bot detection.
- **SearchGupy()**: Consuming the Gupy JSON API directly.

### `internal/storage/`
Handles data persistence.
- **`database.go`**: Manages the SQLite connection, table creation (`CREATE TABLE IF NOT EXISTS jobs...`), and the `AlreadyExists` check to ensure idempotency.

### `internal/discord/`
Manages the bot interaction.
- **`discord.go`**: Wraps the `discordgo` session to send messages to a specific Channel ID defined in the environment variables.

## Installation & Usage

### Prerequisites
- Go 1.25 or higher
- Git

### 1. Clone the repository
```bash
git clone https://github.com/1-AkM-0/empreGo.git
cd empreGo/
```

### 2. Configure Environment variables and Run
You can run by passing the variables inline like this:
```bash
BOT_KEY="your_bot_token_here" CHANNEL_ID="your_channel_id_here" go run .
```
Or create a `.sh` file and export them like this:
```bash
export BOT_KEY="your_bot_token_here"
export CHANNEL_ID="your_channel_id_here"
go run .
```


## Design Choices
During the development of empreGo, several architectural decisions were made to ensure efficiency and reliability:

**Why Go?**

I chose Go because I wanted to learn a language known for its performance and concurrency primitives. The ability to spawn Goroutines allowed the bot to scrape LinkedIn and Gupy in parallel, significantly reducing the total runtime compared to a sequential Python script, for example. Additionally, Go compiles to a single static binary, making deployment to a Linux VM extremely simple.

**Why SQLite?**

Since the bot runs on a schedule (Cron) and the data volume is relatively small (thousands of rows, not millions), a full client-server database like PostgreSQL would be overkill and consume unnecessary resources on the free-tier VM. SQLite offers a serverless, zero-configuration SQL engine that fits perfectly within the application's binary footprint.

**Why Scraping vs. Official APIs?**

Most job platforms (like LinkedIn) do not offer free public APIs for job searching.
- For LinkedIn: I implemented a "Guest Mode" scraper using goquery. To avoid blocking, I mimic a real browser's User-Agent and only scrape public pages.
- For Gupy: Instead of using a heavy tool like Selenium to render the React frontend, I inspected the network traffic and found their internal API endpoints, allowing me to fetch clean JSON data with minimal bandwidth usage.

**Infrastructure**

The project is deployed on a Google Cloud e2-micro instance. I used Cron instead of a continuously running loop to save CPU cycles and prevent memory leaks. The OS wakes up the process, runs the check, and shuts it down, ensuring maximum efficiency.
