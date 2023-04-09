# Phoenix

![Screenshot](screenshot.webp)

Self-hosted start page without the extra stuff.

## Features
- No javascript
- Tiny footprint

## Configuration
Service settings can be set through environment variables.
- `PHOENIX_DB_PATH` - path to the sqlite database.

## Docker-compose example
```yml
services:
  phoenix:
    image: ghcr.io/ordinary-dev/phoenix
    volumes:
      - phoenix:/var/lib/phoenix
    ports:
      - 80:8080
    restart: unless-stopped

volumes:
  phoenix:
```
