
# Smart Home: Legacy Monolith Migration (Strangler Fig Pattern)

This project demonstrates a strategic architectural transition from a legacy Java Spring Boot monolith to a highly scalable, event-driven microservices architecture using **Go**, **Kafka**, and **Kubernetes**.

## 🏗 Architectural Vision

The core objective of this project is to implement the **Strangler Fig** pattern. By intercepting traffic at the API Gateway level, functionality is incrementally "strangled" out of the Java monolith and replaced with high-performance Go microservices.

### C4 Model Documentation

The system design is documented using the C4 model to ensure clarity across different levels of abstraction:

* **Context & Container Diagrams**: High-level system overview and infrastructure boundaries.
* **Component Diagrams**: Detailed view of internal service structures and interactions.
* **ER-Diagrams**: Data modeling for the target state.

You can find all diagrams in the [`/docs`](/docs) directory.

---

## 🛠 Tech Stack

* **Backend**: Go (Microservices), Java 17 (Legacy Monolith), Spring Boot 3.x.
* **Communication**: Asynchronous Event-Driven via **Apache Kafka**.
* **Infrastructure**: Kubernetes (Minikube), Helm (Modular Charts), Terraform (IaaC).
* **Traffic Management**: **Kusk API Gateway** (OpenAPI-driven routing).
* **Storage**: PostgreSQL (Relational persistence).

---

## 🚀 Migration Strategy: From Monolith to Microservices

### Phase 1: The As-Is State (Legacy)

A classic Java monolith handling synchronous HTTP requests for temperature control.

* **Pain Point**: Synchronous interaction with sensors leads to system overhead and high latency.

### Phase 2: The To-Be State (Event-Driven)

Introduction of specialized services to handle load and improve responsiveness:

1. **Temperature Manage Service (Go)**: Decouples device state management from the core logic.
2. **Notification Service (Go)**: Consumes Kafka events to notify devices asynchronously about state changes.
3. **Sensor Stub (Go)**: Simulates hardware behavior for end-to-end integration testing.

---

## 🚦 Getting Started

### Prerequisites

* Minikube
* Terraform
* Helm
* Kusk CLI

### Deployment

1. **Infrastructure**: Provision Kubernetes resources using Terraform.

   ```bash
   cd terraform && terraform apply
   ```

2. **Gateway**: Deploy the API Gateway using the OpenAPI specification.

   ```bash
   kusk deploy -i api.yaml
   ```

3. **Services**: Services are packaged as Helm charts for consistent deployment across environments.

---

## 📈 Roadmap & Future Improvements

To further align with production-grade Senior standards:

* **Observability**: Integration of Prometheus/Grafana for real-time monitoring and OpenTelemetry for distributed tracing across Java and Go services.
* **Resiliency**: Implementation of Circuit Breaker and Retry patterns for inter-service communication.
* **Security**: Migrating the `Users Manage Microservice` to support modern OAuth2/OIDC standards.
