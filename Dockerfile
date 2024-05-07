FROM golang:1.22.2

WORKDIR /go/src/eniqilo
COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE 8080

CMD ["go", "run", "cmd/", "api/", "main.go"]

