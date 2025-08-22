### Build stage (make image smaller with AS)
FROM golang:1.23-alpine3.21 AS builder
WORKDIR /app
# Copy all of current files to current working directory in container
COPY . .
# Build executable file for this package
RUN go build -o main main.go
# Install golang migrate from script (need curl first)
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz

### Run Stage (makes image smaller with just executable file)
# Goes from 500MB to 23MB
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
RUN chmod +x /app/start.sh
RUN chmod +x /app/wait-for.sh

# Container listens on this port (doesn't actually publish the port, just documents it)
EXPOSE 8080

# Run this command when the container starts
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
