# 🚀 Distributed Autoscaler — High-Performance Hashing System

A horizontally scalable, observable, locally deployable distributed system demonstrating production-grade engineering patterns.

This project implements a full distributed architecture featuring:

- **API service (Go)** — receives hashing jobs via REST
- **Worker service (Go)** — consumes and processes jobs from Kafka
- **Kafka** — durable job queue with KRaft consensus
- **Redis** — storage for computed hashes and job status
- **Horizontal Pod Autoscaling (HPA)** — scales workers based on CPU load
- **Prometheus** — comprehensive system metrics
- **Grafana** — real-time visualization dashboard
- **k6** — performance and load testing

Everything is fully containerized and orchestrated on Kubernetes with proper observability and autoscaling capabilities.

---

## 📈 Performance Achieved

Load-test conducted on an Apple M2 Pro (12 CPU / 19 GPU):

| Metric | Result |
|--------|--------|
| **Request Throughput** | 7,600 requests per second sustained |
| **Worker Autoscaling** | Scaled from 1 → 10 pods automatically |
| **CPU Utilization** | Sustained >300% per worker during peak load |
| **System Stability** | Zero crashes, no message loss |
| **Response Time** | <100ms average under load |

This demonstrates real production-grade scaling behavior — a strong portfolio signal for senior engineering roles.

---

## 📦 Architecture Overview

```mermaid
graph TD
    A[Client] -->|HTTP POST /hash| B[API Service]
    A -->|HTTP GET /job/{id}| B
    B -->|Publish Job| C[Kafka Queue]
    C -->|Consume Messages| D[Worker Pool]
    D -->|Store Results| E[Redis]
    B -->|Query Results| E
    
    F[Prometheus] -->|Scrape Metrics| B
    F -->|Scrape Metrics| D
    G[Grafana] -->|Query Data| F
    
    H[Kubernetes HPA] -->|Scale Workers| D
    F -->|CPU Metrics| H
```

**Data Flow:**
1. Client submits hash job to API
2. API publishes job to Kafka topic `jobs`
3. Workers consume jobs and perform SHA256 mining
4. Results stored in Redis with job ID
5. Client queries job status via API

This is a real distributed streaming system with autoscaling & observability built from scratch.

---

## 🗂 Project Structure

```
distributed-autoscaler/
├── api/                    # Go API service
│   ├── cmd/api/main.go    # Application entrypoint
│   ├── internal/
│   │   ├── handlers/      # HTTP handlers (hash, status)
│   │   │   ├── hash_handler.go
│   │   │   └── status_handler.go
│   │   ├── kafka/         # Kafka producer
│   │   │   └── producer.go
│   │   ├── redis/         # Redis client & store
│   │   │   ├── client.go
│   │   │   └── store.go
│   │   ├── metrics/       # Prometheus metrics
│   │   │   └── metrics.go
│   │   ├── server/        # HTTP server setup
│   │   │   └── server.go
│   │   └── config/        # Configuration management
│   └── Dockerfile
├── worker/                 # Go worker service
│   ├── cmd/worker/main.go # Worker entrypoint
│   ├── internal/
│   │   ├── consumer/      # Kafka consumer logic
│   │   │   └── consumer.go
│   │   ├── hashing/       # SHA256 mining algorithm
│   │   │   └── hash.go
│   │   ├── redis/         # Result storage client
│   │   │   └── client.go
│   │   └── metrics/       # Worker performance metrics
│   │       └── metrics.go
│   └── Dockerfile
├── k6/                     # Load testing
│   └── hash_test.js       # k6 performance test script
├── k8s/                    # Kubernetes manifests
│   ├── api/               # API deployment & service
│   │   ├── api-deployment.yaml
│   │   └── api-service.yaml
│   ├── worker/            # Worker deployment & service
│   │   ├── worker-deployment.yaml
│   │   └── worker-service.yaml
│   ├── redis/             # Redis deployment & service
│   │   ├── redis-deployment.yaml
│   │   └── redis-service.yaml
│   ├── kafka/             # Strimzi Kafka cluster configs
│   │   ├── kafka.yaml     # Main Kafka cluster
│   │   ├── kafka-cluster.yaml
│   │   ├── kafka-nodepools.yaml
│   │   └── strimzi-operator.yaml
│   ├── hpa/               # Horizontal Pod Autoscalers
│   │   └── worker-hpa.yaml
│   ├── monitoring/        # Observability stack
│   │   ├── additional-scrape-config.yaml
│   │   └── dashboards/
│   │       └── grafana-dashboard.json
│   └── cluster-config.yaml
├── deployments/           # CI/CD configurations
│   └── github-actions/
├── docs/                  # Project documentation
│   ├── KAFKA_ISSUE_ANALYSIS.md
│   └── readme.md
├── backup-kafka.yaml      # Kafka configuration backup
├── LICENSE                # Apache 2.0 license
├── go.mod                 # Go module dependencies
├── go.sum                 # Go dependency checksums
└── README.md              # This file
```

---

## ▶️ Local Deployment

### Prerequisites
- Kubernetes cluster (Minikube, Kind, or any K8s)
- kubectl configured
- Helm (for monitoring stack)
- Docker (for building images)

### 1. Start Kubernetes Cluster

```bash
# Using Minikube
minikube start --cpus=6 --memory=8g

# Or using Kind
kind create cluster --config k8s/cluster-config.yaml
```

### 2. Build and Deploy Services

```bash
# Build Docker images
docker build -t distributed-autoscaler-api:latest ./api
docker build -t distributed-autoscaler-worker:latest ./worker

# Deploy Kafka (Strimzi)
kubectl apply -f k8s/kafka/kafka.yaml

# Wait for Kafka to be ready
kubectl wait --for=condition=Ready kafka/my-cluster -n kafka --timeout=300s

# Deploy Redis
kubectl apply -f k8s/redis/

# Deploy API and Worker
kubectl apply -f k8s/api/
kubectl apply -f k8s/worker/
```

### 3. Deploy Monitoring Stack

```bash
# Apply Prometheus scrape configuration
kubectl apply -f k8s/monitoring/additional-scrape-config.yaml

# Install Prometheus stack via Helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install monitoring prometheus-community/kube-prometheus-stack -n monitoring --create-namespace

# Enable Metrics Server (required for HPA)
kubectl patch deployment metrics-server -n kube-system \
  --type='json' \
  -p='[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'

# Apply Worker HPA
kubectl apply -f k8s/hpa/worker-hpa.yaml
```

### 4. Verify Deployment

```bash
# Check all pods
kubectl get pods --all-namespaces

# Check HPA status
kubectl get hpa worker-hpa

# Test API endpoint
curl -X POST http://localhost:30080/hash \
  -H "Content-Type: application/json" \
  -d '{"input":"hello world"}'
```

---

## 📊 Accessing Services

### API Service
```bash
kubectl port-forward svc/api 30080:8080
# API available at http://localhost:30080
```

### Prometheus
```bash
kubectl port-forward svc/monitoring-kube-prometheus-prometheus -n monitoring 9090
# Open: http://localhost:9090
```

### Grafana
```bash
kubectl port-forward svc/monitoring-grafana -n monitoring 3000:80
# Open: http://localhost:3000
# Credentials: admin / <get password from secret>
kubectl get secret monitoring-grafana -n monitoring -o jsonpath="{.data.admin-password}" | base64 -d
```

**Import Dashboard:**
1. Navigate to Dashboards → Import
2. Upload `k8s/monitoring/dashboards/grafana-dashboard.json`

---

## 🔥 Load Testing

### Install k6
```bash
brew install k6
# or: apt-get install k6
```

### Run Performance Test
```bash
k6 run k6/hash_test.js
```

**Expected Behavior:**
- Immediate increase in CPU utilization
- HPA triggers automatically (1 → 10 workers)
- Kafka queue depth increases then stabilizes
- System maintains <100ms response times
- Zero message loss or crashes

### Custom Load Test
```bash
# Test different load patterns
k6 run --vus 100 --duration 2m k6/hash_test.js
```

---

## 🛠 Technical Implementation

### Core Components

**API Service Features:**
- RESTful endpoints for job submission and status queries
- Async job processing via Kafka
- Request validation and error handling
- Prometheus metrics for request latency and throughput
- Graceful error responses for queue failures

**Worker Service Features:**
- Kafka consumer group for parallel processing
- SHA256 proof-of-work mining algorithm
- Redis result storage with TTL
- Comprehensive metrics (job duration, nonce attempts)
- Automatic recovery from failures

**Autoscaling Logic:**
- CPU-based scaling with 70% target utilization
- Min/max replicas: 1-10
- 30-second stabilization period
- Custom metrics integration ready

### Monitoring & Observability

**Key Metrics:**
- API request rate and latency
- Worker job processing time
- Kafka consumer lag
- CPU/memory utilization by service
- HPA scaling events

**Alerting Ready:**
- High error rate thresholds
- Kafka consumer lag alerts
- Pod restart notifications
- Resource utilization warnings

---

## 🧪 Development & Testing

### Local Development
```bash
# Run API locally
go run ./api/cmd/api/main.go

# Run worker locally
go run ./worker/cmd/worker/main.go

# Run tests
go test ./...
```

### Integration Testing
```bash
# Test end-to-end flow
curl -X POST http://localhost:8080/hash -d '{"input":"test"}'
# Returns: {"job_id":"uuid-here"}

# Check status
curl http://localhost:8080/job/uuid-here
# Returns: {"status":"done","hash":"0000...","nonce":12345}
```

---

## 🔒 Security Considerations

- **Authentication**: Ready for JWT/OAuth integration
- **Input Validation**: JSON schema validation for all inputs
- **Resource Limits**: Pod resource constraints and limits
- **Network Policies**: Kubernetes network isolation ready
- **Secrets Management**: Kubernetes secrets for credentials

---

## 🚀 Production Readiness

### Scalability Features
- Horizontal pod autoscaling
- Kafka partition scaling
- Redis clustering support
- Load balancer ready

### Reliability Features
- Health checks and readiness probes
- Graceful shutdown handling
- Retry logic with exponential backoff
- Circuit breaker patterns ready

### Monitoring Features
- Structured logging
- Distributed tracing ready
- Custom business metrics
- SLA monitoring dashboard

---

## 🏁 Key Achievements

This project demonstrates senior-level engineering capabilities:

✅ **System Design**: Event-driven architecture with proper separation of concerns  
✅ **Performance**: 7.6k RPS sustained on laptop hardware  
✅ **Scalability**: Automatic horizontal scaling based on load  
✅ **Observability**: Comprehensive monitoring with Prometheus/Grafana  
✅ **Reliability**: Zero message loss, graceful error handling  
✅ **DevOps**: Full CI/CD ready with containerization and orchestration  
✅ **Testing**: Performance testing with realistic load patterns  

---

## 📚 Learning Outcomes

Building this system provides experience with:
- Distributed systems design patterns
- Message queue architecture (Kafka)
- Container orchestration (Kubernetes)
- Microservices communication
- Performance optimization
- Monitoring and observability
- Autoscaling strategies
- Load testing methodologies

---

## 🔮 Future Enhancements

- **GCP Deployment**: GKE, Cloud Load Balancer, Cloud SQL
- **Advanced Autoscaling**: Custom metrics based on Kafka lag
- **Security**: mTLS, RBAC, API gateway
- **Performance**: Redis clustering, Kafka partition tuning
- **Features**: Job priorities, batch processing, web UI

---

## 📄 License

Apache License 2.0 - Perfect for professional portfolios and commercial use.

---

## 🤝 Contributing

This is a portfolio project demonstrating engineering capabilities. Feel free to fork, modify, and use as a learning resource.

---

**Built with ❤️ to showcase distributed systems engineering expertise**

*This project represents end-to-end system development from architecture design through deployment and scaling validation.*
