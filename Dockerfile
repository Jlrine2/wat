FROM golang:latest as build

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s"  -o /wat ./cmd/wat

FROM scratch as run

# Copy the application executable from the build image
COPY --from=build /wat /wat

ENTRYPOINT ["/wat"]