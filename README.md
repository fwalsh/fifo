#  fifo – Items API with Go, Postgres, Docker & CircleCI

**fifo** is a nickname for Fiona (me!).  

In this case, it’s also the name of a simple **Items API** written in Go — a simple but complete reference app designed to showcase a modern CI/CD pipeline.  

The service exposes a few endpoints:

- `GET /` → Landing page (HTML with styled welcome)  
- `GET /health` → Health check (JSON: `{"status":"ok"}`)  
- `GET /items` → List stored items (JSON array)  
- `POST /items` → Add a new item (accepts JSON or form input)  


Behind the scenes:  
- **Go** for the app logic  
- **Postgres** as the database  
- **Docker** for containerization  
- **CircleCI** to test, build, and publish images to AWS ECR using **OIDC** (no static creds)  

It’s both a **reference pipeline** and a way to show how everything maps together end-to-end.

---


##  Architecture

```text
             ┌─────────────┐
             │   Client    │
             │ (curl, app) │
             └──────┬──────┘
                    │  HTTP
                    ▼
            ┌───────────────┐
            │   FIFO App    │   (Go REST API)
            │   (Docker)    │
            └──────┬────────┘
                   │  SQL
                   ▼
           ┌─────────────────┐
           │   Postgres DB   │
           │  (Docker side)  │
           └─────────────────┘

   CircleCI Pipeline:
   ┌───────────────┐   ┌───────────────┐   ┌───────────────┐
   │    Test Job   │ → │   Build Job   │ → │   Push to ECR │
   └───────────────┘   └───────────────┘   └───────────────┘
                           (only on main)

```
 Running Locally

Make sure you have Docker and docker-compose installed.

Open a terminal (bash/zsh on macOS/Linux, or PowerShell on Windows) and run the following commands:


# Clone the repo
git clone https://github.com/fwalsh/fifo.git
cd fifo

# Start services (database, migrations, and app)
docker-compose up -d


Once the services are running:

# Health check
curl http://localhost:8080/health
```text
# -> {"status":"ok"}
```


# Create an item (JSON)
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \
  -d '{"name":"peach"}'


# List items
curl http://localhost:8080/items
```text
# -> [{"id":1,"name":"peach","created_at":"..."}]
```

Or open your browser and visit:
 http://localhost:8080/

You’ll see a simple landing page, and you can interact with the API from there.


 CI/CD Pipeline

The CircleCI config (.circleci/config.yml) includes:

Test job:
Runs Go tests against a Postgres sidecar, collects JUnit results in CircleCI.

Build job:
Builds Docker image, halts if only docs changed, stores fifo-app.tar as artifact.

Push job:
On merge to main, uses OIDC to authenticate to AWS and push the image to ECR.


 Branch-based conditions:

Tests run on all branches (PRs, feature branches, main).

Docker build and push to AWS ECR run only on main.
This ensures quick feedback everywhere while limiting heavier jobs to production-ready changes.


 Features & Optimizations

 Postgres sidecar for realistic testing

 Multi-stage Docker builds for lean images

 OIDC → AWS: no static credentials

 Skips unnecessary builds for docs-only changes

 Publishes artifacts to both CircleCI + AWS ECR

 Automatic migrations ensure schema is ready (via `migrate` sidecar)


 Future Improvements

Lock IAM trust to main branch only

Replace broad ECR IAM policy with least-privilege custom role

Add ECS (Fargate) deployment stage for full CD

Parallelize/conditionalize tests for speedier builds


👩🏼‍💻 Author

👋 Hi, I’m Fiona. This project was built as part of a CircleCI field engineer challenge.


