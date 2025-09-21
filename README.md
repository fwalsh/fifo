# fifo

Testing out some tools and functionalities.

# Items API

A simple Go web service with a Postgres database.  
Built for demonstrating CircleCI pipelines, Docker, and ECS deployment.  

## Endpoints
- `GET /health` → Health check  
- `POST /items` → Add a new item  
- `GET /items` → List all items  

## Running Locally
```bash
go run main.go
curl localhost:8080/health

