# Delivery REST API Service

A service that registers couriers, adds new orders, and assigns them to couriers.

## Features

1. REST API with 7 basic methods.
2. Courier rating calculation.
3. Rate limiter for incoming requests.

## REST API

- `POST /couriers`: Load couriers list.
- `GET /couriers/{courier_id}`: Get specific courier info.
- `GET /couriers`: Get all couriers info.
- `POST /orders`: Add orders list.
- `GET /orders/{order_id}`: Get specific order info.
- `GET /orders`: Get all orders info with pagination.
- `POST /orders/complete`: Mark orders as completed.

## Courier Ratings

- `GET /couriers/meta-info/{courier_id}`: Get courier's earnings and rating.

## Rate Limiter

- Limits load to 10 RPS for each endpoint.

## Requirements

1. Dockerfile for build, configuration, and launch.
2. Service on port 8080.
3. PostgreSQL 15.2 on port 5432, with user=postgres and password=password.

