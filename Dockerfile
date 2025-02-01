### Build stage (make image smaller with AS)
FROM golang:1.23-alpine3.21 AS builder
WORKDIR /app

# Copy all of current files to current working directory in container
COPY . . 

# Build executable file for this package
RUN go build -o main main.go

# Container listens on this port (doesn't actually publish the port, just documents it)
EXPOSE 8080

### Run Stage (makes image smaller with just executable file)
# Goes from 500MB to 23MB
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

# Run this command when the container starts
CMD [ "/app/main" ]