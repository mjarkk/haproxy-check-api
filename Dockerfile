# Create a self signed certificate for testing
FROM alpine as cert

# Add openssl to the container
RUN apk add --update openssl && rm -rf /var/cache/apk/*

# Generate self-signed certivicates
RUN mkdir /etc/certs && \
  openssl req -x509 -newkey rsa:4096 -nodes -keyout /etc/certs/key.pem -out /etc/certs/cert.pem -days 365 -subj "/C=GB/OU=HAPROXY-CHECK-API" && \
  cat /etc/certs/cert.pem /etc/certs/key.pem > /etc/certs/fullKey.pem



# Build the go files
FROM golang:1.11-alpine as build

# Make the directories and copy the needed files
RUN mkdir -p /go/src/github.com/mjarkk/haproxy-check-api
WORKDIR /go/src/github.com/mjarkk/haproxy-check-api
COPY ./ ./

# buidl the program and chmod it
RUN GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -v -a -installsuffix cgo



# Build the api container
FROM haproxy:2.0.1

# Copy over the binary
COPY --from=build /go/src/github.com/mjarkk/haproxy-check-api/haproxy-check-api /haproxy-check-api

# Make the certs folder and chmod the binary
RUN mkdir /etc/certs

# Copy over the certificate
COPY --from=cert /etc/certs/fullKey.pem /etc/certs/fullKey.pem

# haproxy-check-api runs on port 8223
EXPOSE 8223

# Run the haproxy binary
CMD /haproxy-check-api
