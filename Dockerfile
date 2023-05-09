FROM golang:1.20-alpine

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY src/go.mod src/go.sum ./
#RUN go mod download && go mod verify

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN mkdir -p /usr/local/bin/
RUN go mod tidy
RUN go build -o /usr/local/bin/app ./server.go

EXPOSE 8080

ENV DB_HOST=host.docker.internal
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=password
ENV DB_NAME=postgres
ENV DISABLE_RATE_LIMITER=false

CMD ["/usr/local/bin/app"]
