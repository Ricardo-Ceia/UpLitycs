# Multi-stage build

# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci 
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend - Use latest Go
FROM golang:alpine AS backend-builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o statusframe main.go

# Stage 3: Final lightweight image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /app/statusframe .

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Copy database migrations
COPY db/ ./db/

EXPOSE 8080

#CMD ["./statusframe"]

ENTRYPOINT ["./statusframe"]