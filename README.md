# Project README

## Overview

This "BeesBiz" project focuses on the literal management of beekeeping operations. It provides tools for tracking and monitoring bee health, hive conditions, and honey production.

## Architecture Components

### Front End

- **Grafana**: Visualization and monitoring tool that interfaces with Prometheus to display metrics and analytics.
- **Web UI**: User interface that interacts with the backend services.

### Core Services

- **Backend Services**: The core logic and processing unit of the application, handling business operations and data processing.
- **API Gateway**: Acts as a single entry point for the frontend to access backend services and routes requests accordingly. (TO BE IMPLEMENTED)

### Data Layer

- **RabbitMQ (RMQ)**: A message broker that facilitates asynchronous communication between services.
- **PostgreSQL (PG)**: Relational database for persistent data storage.
- **TiKV Cache**: A distributed transactional key-value storage system for caching and fast access to frequently used data. (TO BE IMPLEMENTED)

### Infrastructure (TO BE IMPLEMENTED)

- **Hashicorp Vault**: Manages secrets and sensitive data, ensuring secure access to credentials and configuration.
- **Consul**: Provides service discovery and configuration capabilities, enabling services to find each other in a dynamic environment.

### Monitoring

- **Prometheus**: Monitoring and alerting toolkit that collects metrics from configured services and stores them for analysis.

### Data Mockup

- **hive.rs**: Custom genserver (via actix) rust service that simulates the hive sensor data.

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
---
config:
  layout: elk
---
graph TB
    subgraph "Frontend Plane"
        subgraph "Frontend"
            Grafana[Grafana]
            UI[Web UI]
        end
    end

    subgraph "Backend Plane"
        subgraph "Core Services"
            BE[Backend Services]
        end
        subgraph "Data Layer"
            RMQ[(RabbitMQ)]
            PG[(PostgreSQL)]
            TiKV[(TiKV Cache)]
        end
        subgraph "Infrastructure"
            Gateway[Consul]
            Vault[Hashicorp Vault]
        end
        subgraph "Monitoring"
            Prometheus[Prometheus]
        end
    end

    subgraph "Regional Cluster"
        direction TB
        Gateway --> AP1[Regional Apiary Pod 0]
        AP1 --> H11[Hive Container ...] & H12[Hive Container N] & H13[Hive Container 1]
    end

    subgraph "Regional Cluster"
        direction TB
        Gateway --> AP2[Regional Apiary Pod 0]
        Gateway --> AP3[Regional Apiary Pod N]
        AP2 --> H21[Hive Container ...] & H22[Hive Container N] & H23[Hive Container 1]
        AP3 --> H31[Hive Container ...] & H32[Hive Container N] & H33[Hive Container 1]
    end

    UI --> Gateway
    BE <--> RMQ & PG & TiKV
    Vault --> RMQ & PG & TiKV
    Gateway --> BE
    Vault <--> BE
    Prometheus -->  BE & PG & RMQ & TiKV
    Grafana --> Prometheus

    classDef default fill:#F7F7F7,stroke:#303030,stroke-width:1px;
    classDef cluster fill:#326CE5,stroke:#303030,stroke-width:2px,color:#fff;
    classDef apiary fill:#A1C4FD,stroke:#326CE5,stroke-width:3px;
    classDef hive fill:#F7F7F7,stroke:#326CE5,stroke-width:1px;
    class GC1,GC2,GC3 cluster;
    class AP1,AP2,AP3 apiary;
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
