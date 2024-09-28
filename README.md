# Project README

## Overview

This "BeesBiz" project focuses on the literal management of beekeeping operations. It provides tools for tracking and monitoring bee health, hive conditions, and honey production.

## Architecture Components

### Front End
- **Grafana**: Visualization and monitoring tool that interfaces with Prometheus to display metrics and analytics.
- **Web UI**: User interface that interacts with the backend through the API Gateway.

### Core Services
- **Backend Services**: The core logic and processing unit of the application, handling business operations and data processing.
- **API Gateway**: Acts as a single entry point for the frontend to access backend services and routes requests accordingly.
- **Configuration Management**: Manages application settings and configurations, ensuring the right parameters are used across services.

### Data Layer
- **RabbitMQ (RMQ)**: A message broker that facilitates asynchronous communication between services.
- **PostgreSQL (PG)**: Relational database for persistent data storage.
- **TiKV Cache**: A distributed transactional key-value storage system for caching and fast access to frequently used data.

### Infrastructure
- **Dapr**: A distributed application runtime that simplifies microservices development and communication.
- **Hashicorp Vault**: Manages secrets and sensitive data, ensuring secure access to credentials and configuration.
- **Consul**: Provides service discovery and configuration capabilities, enabling services to find each other in a dynamic environment.

### Monitoring
- **Prometheus**: Monitoring and alerting toolkit that collects metrics from configured services and stores them for analysis.

## Diagram Legend

- **Clusters**:
  - Front End
  - Core Services
  - Data Layer
  - Infrastructure
  - Monitoring
  - Geo-Regions

- **Component Classifications**:
  - **Apiary**: Denotes pod-level services that handle specific group of hives.
  - **Hive**: Represents individual hive.

```mermaid
graph TB
    subgraph "Front End"
        Grafana[Grafana]
        UI[Web UI]
    end

    subgraph "Core Services"
        BE[Backend Services]
        Gateway[API Gateway]
        ConfigMgmt[Configuration Management]
    end

    subgraph "Data Layer"
        RMQ[(RabbitMQ)]
        PG[(PostgreSQL)]
        TiKV[(TiKV Cache)]
    end

    subgraph "Infrastructure"
        Dapr[Dapr]
        Vault[Hashicorp Vault]
        Consul[Consul]
    end

    subgraph "Monitoring"
        Prometheus[Prometheus]
    end

    subgraph "Geo-Regions"
        GC1[Region 1]
        GC2[Region 2]
        GC3[Region 3]
    end

    UI --> Gateway
    Gateway --> Consul
    BE <--> Dapr
    Dapr <--> RMQ & PG & TiKV
    Prometheus --> GC1 & GC2 & GC3 & BE & Dapr
    BE <--> Vault & Consul
    ConfigMgmt --> RMQ

    Consul --> Prometheus
    Grafana --> Prometheus

    RMQ --> GC1 & GC2 & GC3

    subgraph "Region Details"
        direction TB
        GC1 --> AP1[Apiary Pod 1]
        AP1 --> H11[Hive 1] & H12[Hive 2] & H13[Hive 3]
    end

    classDef default fill:#e0f7fa,stroke:#006064,stroke-width:2px;
    classDef cluster fill:#bbdefb,stroke:#0d47a1,stroke-width:2px;
    classDef apiary fill:#c8e6c9,stroke:#1b5e20,stroke-width:2px;
    classDef hive fill:#ffe082,stroke:#f57f17,stroke-width:1px;
    class GC1,GC2,GC3 cluster;
    class AP1 apiary;
    class H11,H12,H13 hive;
```

## Getting Started

To set up the project locally:

1. Clone the repository.
2. Install dependencies.
3. Configure the environment variables for the database and messaging services.
4. Start the application using Docker or your preferred method.
5. Access the Web UI at `http://localhost:PORT` and monitor with Grafana at `http://localhost:GRAFANA_PORT`.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any enhancements or bugs.

## License

See LICENSE
