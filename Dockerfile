# Build stage
FROM golang:1.26 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/upver .

# Runtime stage
FROM alpine:3.23
RUN apk add --no-cache git ca-certificates
COPY --from=build /out/upver /usr/local/bin/upver
ENTRYPOINT ["upver"]

