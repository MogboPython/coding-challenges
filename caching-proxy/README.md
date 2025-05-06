# Caching Proxy Server

This project is a simple solution to the [Roadmap.sh Caching Server Problem](https://roadmap.sh/projects/caching-server). It implements a simple HTTP proxy server with caching capabilities. The server forwards requests to a target server, caches the responses, and serves cached responses for subsequent requests to improve performance and reduce load on the target server.

## Features

- **Proxy Functionality**: Forwards HTTP requests to a specified target server.
- **In-Memory Caching**: Stores responses in an in-memory map.

## How It Works

1. The server listens on a specified port and forwards incoming requests to a target server.
2. Responses from the target server are cached in memory for a configurable TTL.
3. If a request matches a cached response and the cache entry is still valid, the server serves the cached response.
4. Expired cache entries are automatically cleaned up in the background using goroutines.

## Usage

1. Clone the repository:
   ```bash
   git clone https://github.com/MogboPython/coding-challenges.git
   cd caching-proxy
   ```

2. Build the project:
    `go build`

3. Run the server:
    `./caching-proxy --port <PORT> --origin <ORIGIN_URL>`

4. Clear the cache
    `./caching-proxy --clear-cache`

## Conclusion
This is a minimal implementation intended as a starting point for more complex caching proxy solutions.
