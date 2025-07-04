# Player Activity Report Service

This is a Golang microservice that generates player activity reports. It fetches betting statistics from a MySQL database, enriches them with country information from the REST Countries API, and provides the combined data through a REST endpoint.

## How to Run

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make

Alternatively, if you have Nix, run `nix develop` to enter a shell with all dependencies.

### Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/Nikola-Milovic/vyking
    cd vyking
    ```

2. **Configure environment:**

    ```bash
    cp .env.example .env
    # Edit .env with your configuration
    ```

3. **Run the app and database:**

    This command will start the Docker containers, run database migrations, and seed the database.

    ```bash
    make dev-setup
    ```

The API will be available at `http://localhost:8080`. API documentation is available via the Swagger UI at `http://localhost:8080/docs`.

You can run `make help` to see all available commands.

## Testing

For a clean and reliable testing experience, the project uses `mockgen` to isolate components during unit tests and Testcontainers to run integration tests against a real, containerized MySQL database.

**Generate mocks:**

```bash
make generate-mocks
```

**Run tests:**

```bash
make test
```

**Run tests with coverage:**

```bash
make test-coverage
```

## Example

To get the top 3 countries by player activity, run:

```bash
curl "http://localhost:8080/country-player-stats?limit=3"
```

The response will be a JSON object containing player statistics and country details. If the external country API is unavailable, `country_info` will be `null`.

```json
{
  "stats": [
    {
      "country_code": "RS",
      "player_count": 10,
      "total_bets": 13388.83,
      "avg_bet_per_player": 1338.883,
      "country_info": {
        "name": "Serbia",
        "region": "Europe",
        "borders": [
          "BIH",
          "BGR",
          "HRV",
          "HUN",
          "UNK",
          "MKD",
          "MNE",
          "ROU"
        ]
      }
    },
    {
      "country_code": "BR",
      "player_count": 7,
      "total_bets": 11231,
      "avg_bet_per_player": 1604.428571,
      "country_info": {
        "name": "Brazil",
        "region": "Americas",
        "borders": [
          "ARG",
          "BOL",
          "COL",
          "GUF",
          "GUY",
          "PRY",
          "PER",
          "SUR",
          "URY",
          "VEN"
        ]
      }
    },
    {
      "country_code": "DE",
      "player_count": 5,
      "total_bets": 6171.37,
      "avg_bet_per_player": 1234.274,
      "country_info": {
        "name": "Germany",
        "region": "Europe",
        "borders": [
          "AUT",
          "BEL",
          "CZE",
          "DNK",
          "FRA",
          "LUX",
          "NLD",
          "POL",
          "CHE"
        ]
      }
    }
  ]
}
```

## Caching Implementation

To speed up responses, the service uses a simple inmemory cache with a time-to-live (TTL) and a least recently used (LRU) eviction policy. I chose this approach over something like Redis to keep the project lightweight and free of (not critical) external dependencies. Since the country data doesn't change often, this simple cache is a reasonable fit. Plus, it's built behind an interface, so swapping it out later would be straightforward. Using a third party dependency makes no sense in this case and if there was a need for one, I would still keep it behind an internal interface.

## Retrospective

### Challenges

- I had to get reacquainted with MySQL specific syntax for stored procedures, as I primarily work with PostgreSQL.
- Finding the right balance of features and simplicity in Go's OpenAPI/Swagger tooling took some experimentation. From what I've researched, the tooling is quite wonky and there are a ton of opinionated libraries, I landed on `swaggest/rest`, which fits the requriments and has an okay DX. I took the code -> schema approach instead of first generating openapi schema and then using code generation for the http handler.
- A key focus was to try and fit in as many good practices as possible and to squeeze in as many interesting approaches as possible as I would for a real project. But of course considering the simplicity of the task there is only so much that can be done.

### Time

I spent approximately 9 hours working on this task, spread throughout 4 days.

