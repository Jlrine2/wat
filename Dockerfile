FROM golang:latest as buildserver

WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download the Go module dependencies
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s"  -o /wat ./cmd/wat



FROM node:latest as buildclient

WORKDIR /app/client

# Copy package files
COPY client/package*.json ./

# Install dependencies
RUN npm install

# Copy client source code
COPY client/ .

# Build the Svelte app
RUN npm run build


FROM scratch as run

# Copy the built Svelte app into the final image
COPY --from=buildclient /app/client/build /client

# Copy the application executable from the build image
COPY --from=buildserver /wat /wat
COPY --from=buildserver /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/wat"]