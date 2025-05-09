# Event Streaming Application


### Challenge Prompt done via Cursor
```
Challenge Description

Task: Implement a simplified, event-driven streaming application with the following functionality…

HTTP API endpoint to accept incoming events (REST or gRPC optional).

Each event contains:

ID (string)

Timestamp (int64, Unix timestamp)

Payload (string)

Store events in memory safely for concurrent use.

Background worker(s) concurrently consume events, apply a simple transformation (e.g., uppercase payload, add metadata), and "publish" transformed events.

Publish by printing to stdout, logging, or simulating Pub/Sub via channels.

Demonstrate graceful shutdown (clean shutdown upon interrupt).

Bonus Points

These areas are optional and intended to show deeper understanding and broader perspective. The more the candidate get through them, the merrier, as they can help us differentiate between strong and excellent candidates.

1. Written, Working Tests:

Meaningful unit tests verifying core logic (event handling, API endpoints, data transformations, concurrency handling).

2. Infrastructure & Hosting Considerations (AWS):

In your README, briefly describe how you'd deploy this app using AWS infrastructure (no actual deployment required).

Include:

AWS services you would consider (e.g., ECS/EKS, Fargate, EC2, Lambda, DynamoDB, RDS, API Gateway, ELB).

How you'd manage scaling, high availability, and fault tolerance.

Deployment strategies (CI/CD, IaC tools like Terraform or CloudFormation).

3. Showcase System Under Stress/Load:

Briefly demonstrate or describe how your app performs under high volumes of incoming events:

What do you expect to be the bottlenecks or breaking points?

Have you considered how your choice of concurrency strategy affects performance under load?

What adjustments would you recommend to handle significantly increased event throughput?

4. Observability & Monitoring Considerations:

Explain in your README how you'd add observability to your system:

Logging strategy, metrics you'd track, monitoring and alerting systems (e.g., Prometheus, CloudWatch, Grafana).

Brief explanation of how you'd debug or troubleshoot issues in production.
```

A concurrent event-driven streaming application that accepts events via HTTP, processes them with background workers, and publishes transformed events.

## Features

- HTTP API endpoint to accept incoming events
- Concurrent event processing with multiple workers
- Thread-safe in-memory event storage
- Event transformation (uppercase payload)
- Graceful shutdown handling

## Architecture

The application consists of the following components:

1. **API Server**: HTTP server exposing endpoints for event submission and retrieval
2. **Event Store**: Thread-safe in-memory store with a channel-based notification system
3. **Worker Pool**: Multiple workers processing events concurrently
4. **Main App**: Coordinates startup and shutdown of all components

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Docker and Docker Compose (optional, for containerized deployment)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/event-streaming-app.git
cd event-streaming-app

# Install dependencies
go mod download

# Build the application
go build -o event-processor ./cmd
```

### Running the Application

#### Using Go directly

```bash
./event-processor
```

#### Using Docker Compose

```bash
docker-compose up --build
```

The HTTP server will start on port 8080.

### API Endpoints

- `POST /events` - Submit a new event
  ```json
  {
    "id": "event-123",
    "timestamp": 1625097600,
    "payload": "example payload"
  }
  ```

- `GET /events` - Retrieve all events
- `GET /events/{id}` - Retrieve a specific event by ID
- `GET /health` - Health check endpoint

### Sending Test Events

You can run the test client to send test events:

```bash
go run scripts/main.go -type=client
```

### Running Load Tests

To run a load test that simulates multiple clients sending events concurrently:

```bash
go run scripts/main.go -type=load
```

### Running Tests

To run the unit tests:

```bash
go test ./...
```

## Development Notes

### Project Structure

```
.
├── app
│   ├── api         # HTTP API server implementation
│   └── processor   # Event processing workers
├── cmd             # Application entry point
├── config          # Configuration files
├── internal
│   └── models      # Data models and event store
├── scripts         # Test client and load testing scripts
└── docker-compose.yml
```

### Adding More Worker Types

To add new types of event processors:

1. Create a new processor in the `app/processor` directory
2. Implement the processing logic in the new processor
3. Update the main application to start the new processor type

## AWS Infrastructure Considerations

### Services

For deploying this application in AWS, the following services would be appropriate:

#### **Compute Options**
- **ECS with Fargate**: Good for containerized microservices with auto-scaling
- **Kubernetes on EKS**: For complex orchestration requirements
- **Lambda with API Gateway**: For serverless approach (would require refactoring)

#### **Data Storage**
- **DynamoDB**: For event storage with high throughput (replacing in-memory store)
- **Amazon SQS**: For reliable message queuing between components
- **Amazon Kinesis**: For real-time streaming data processing

#### **API Fronting**
- **API Gateway**: To secure, manage, and route API requests
- **Application Load Balancer**: For load balancing HTTP traffic

### Scaling & High Availability

- **Auto Scaling Groups**: To adjust capacity based on demand
- **Multi-AZ Deployment**: For high availability
- **Read Replicas**: For database scaling (if using RDS)
- **DynamoDB Global Tables**: For multi-region data replication

### Deployment Strategy

- **CI/CD Pipeline**: Using AWS CodePipeline or GitHub Actions
- **Infrastructure as Code**: Using Terraform or AWS CloudFormation
- **Blue/Green Deployment**: For zero-downtime updates
- **Canary Releases**: For gradual rollouts

### AWS Architecture Diagram

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│                 │     │                  │     │                 │
│   API Gateway   ├────►│  ALB / NLB       ├────►│  ECS Fargate    │
│                 │     │                  │     │  (Container)    │
└─────────────────┘     └──────────────────┘     └────────┬────────┘
                                                          │
                                                          ▼
                         ┌───────────────┐     ┌──────────────────┐
                         │               │     │                  │
                         │   CloudWatch  │◄────┤   DynamoDB       │
                         │               │     │   (Event Store)  │
                         └───────────────┘     └──────────────────┘
```

## Performance Under Load

### Potential Bottlenecks

1. **In-Memory Event Storage**: The current implementation stores all events in memory, which limits scalability based on available RAM.

2. **Channel Buffer Size**: If events arrive faster than workers can process them, the channel buffer may fill up.

3. **HTTP Server Throughput**: Under high concurrency, the HTTP server might become a bottleneck.

### Scaling Strategies

1. **Replace In-Memory Store**: Switch to a distributed database like DynamoDB or Cassandra.

2. **Add Load Balancing**: Distribute traffic across multiple API server instances.

3. **Implement Backpressure**: Rate-limit API requests when the system is overloaded.

4. **Horizontal Scaling**: Run multiple instances of the application behind a load balancer.

5. **Worker Pool Tuning**: Adjust the number of workers based on CPU cores and workload.

## Observability & Monitoring

### Logging Strategy

- **Structured Logging**: Use JSON-formatted logs with consistent fields
- **Log Aggregation**: Collect logs with AWS CloudWatch Logs or ELK stack
- **Log Levels**: Implement different log levels (DEBUG, INFO, ERROR) for filtering

### Key Metrics to Track

- **Request Rate**: Incoming events per second
- **Processing Latency**: Time to process each event
- **Queue Depth**: Number of events waiting to be processed
- **Error Rate**: Failed event processing attempts
- **System Resources**: CPU, memory, network usage

### Monitoring & Alerting Tools

- **Prometheus**: For metrics collection and alerting
- **Grafana**: For visualization and dashboards
- **AWS CloudWatch**: For AWS-integrated monitoring
- **X-Ray**: For distributed tracing

### Troubleshooting Strategy

1. **Distributed Tracing**: Implement trace IDs across components
2. **Health Check Endpoints**: Add detailed health checks for each component
3. **Circuit Breakers**: Implement to prevent cascading failures
4. **Correlation IDs**: Track event flow through the system

## License

MIT 
