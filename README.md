# 🌀 fifo – Items API with Go, Postgres, Docker & CircleCI

**fifo** is a nickname for Fiona (me! 👋).  
Here it’s also the name of a simple **Items API** written in Go — a simple but complete reference app designed to showcase a modern CI/CD pipeline.  

The service exposes a couple of endpoints:  
- `GET /health` → Health check  
- `GET /items` → List stored items  
- `POST /items` → Add a new item (🍎, 🍐, 🍇, 🌀 … you get the idea!)

Behind the scenes:  
- **Go** for the app logic  
- **Postgres** as the database  
- **Docker** for containerization  
- **CircleCI** to test, build, and publish images to AWS ECR using **OIDC** (no static creds 🚫🔑)  

It’s both a **reference pipeline** and a fun way to show how everything maps together end-to-end.

---

## 🌀 Architecture

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
🚀 Running Locally
Make sure you have Docker and docker-compose installed.

bash
Copy code
# Clone repo
git clone https://github.com/fwalsh/fifo.git
cd fifo

# Start services
docker-compose up -d

# Test the app
curl http://localhost:8080/health
curl -X POST http://localhost:8080/items -H "Content-Type: application/json" -d '{"name":"peach"}'
curl http://localhost:8080/items
You’ll see items persisted in Postgres 🎉.

🔄 CI/CD Pipeline
The CircleCI config (.circleci/config.yml) includes:

Test job:
Runs Go tests against a Postgres sidecar, collects JUnit results in CircleCI.

Build job:
Builds Docker image, halts if only docs changed, stores fifo-app.tar as artifact.

Push job:
On merge to main, uses OIDC to authenticate to AWS and push the image to ECR.

🔀 Branch-based conditions:

Tests run on all branches (PRs, feature branches, main).

Docker build and push to AWS ECR run only on main.
This ensures quick feedback everywhere while limiting heavier jobs to production-ready changes.

🐾 Features & Optimizations
🐘 Postgres sidecar for realistic testing

🐳 Multi-stage Docker builds for lean images

🔒 OIDC → AWS: no static credentials

✂️ Skips unnecessary builds for docs-only changes

🌱 Automatic migrations ensure schema is ready

📦 Publishes artifacts to both CircleCI + AWS ECR

🔮 Future Improvements
Lock IAM trust to main branch only

Replace broad ECR IAM policy with least-privilege custom role

Add ECS (Fargate) deployment stage for full CD

Parallelize/conditionalize tests for speedier builds

🧑‍💻 Author
👋 Hi, I’m Fiona. This project was built as part of a CircleCI field engineer challenge.
If you’d like to chat DevOps tooling, CI/CD, or fun side projects — let’s connect! 🌟

