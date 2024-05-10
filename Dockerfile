# build stage
FROM golang:alpine AS builder
WORKDIR /go/src/eniqilo-store
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM builder AS build
RUN go build -o /bin/api ./cmd/api/main.go

# final stage
FROM alpine:latest
COPY --from=build /bin/api ./eniqilo-store
CMD ["./eniqilo-store"]
EXPOSE 8080


