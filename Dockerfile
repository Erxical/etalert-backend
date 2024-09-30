FROM golang:1.20.4-alpine3.18 AS builder
WORKDIR /service
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o etalert_backend main.go

# Stage 2: Run
FROM alpine:3.18.4
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /service/etalert_backend .
EXPOSE 3000
CMD [ "./etalert_backend" ]