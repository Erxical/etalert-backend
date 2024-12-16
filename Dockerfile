FROM golang:1.23.3-alpine3.20 AS builder
WORKDIR /service
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o etalert_backend main.go

# Stage 2: Run
FROM debian:bookworm-slim
RUN apt update && apt install -y ca-certificates
WORKDIR /root
COPY --from=builder /service/etalert_backend .
EXPOSE 3000
CMD [ "./etalert_backend" ]