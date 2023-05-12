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

## Project Setup

1. Clone the project repository.

    ```bash
    git clone <project-repo-url>
    ```

2. Build the Docker image and start the Docker container using Docker Compose.

    ```bash
    docker-compose up --build
    ```

The application should now be running at `http://localhost:8080`.

## Running Tests

`test.sh` will run the project tests in a Docker environment. Rate limiter is disabled while running the tests.

1. Ensure that the script is executable:

    ```bash
    chmod +x test.sh
    ```

2. Run the tests:

    ```bash
    bash test.sh
    ```

This script will:

1. Rename any existing `db` and `app` containers to avoid conflict.
2. Create a Docker network called `enrollment`.
3. Build the Docker images and start the Docker containers using Docker Compose.
4. Execute the test cases.
5. Cleanup the Docker environment.
6. Rename the `db` and `app` containers back to their original names.

## Running the Application using Dockerfile

In case running the application using only Docker (without Docker Compose), follow the steps below.
1. Set up network `my-network` separately.

    ```bash
    docker network create my-network
    ```

2. Set up the database service `db` separately.

    ```bash
    docker run -d --name db -e POSTGRES_PASSWORD=password -p 5432:5432 --network my-network postgres
    ```

3. Build a Docker image from the Dockerfile. You can do this using the `docker build` command.

    ```bash
    docker build -t my-app .
    ```

4. After the image is built, you can run it using the `docker run` command.

    ```bash
    docker run -d --name app -p 8080:8080 --network my-network my-app
    ```