FROM haproxy:1.9

# Copy over the binary
COPY ./haproxy-check-api /haproxy-check-api

# Update the container and install openssl
RUN apt-get update -y
RUN apt-get install openssl -y

# Generate self-signed certivicates
RUN mkdir /etc/certs
RUN openssl req -x509 -newkey rsa:4096 -nodes -keyout /etc/certs/key.pem -out /etc/certs/cert.pem -days 365 -subj "/C=GB/OU=HAPROXY-CHECK-API"
RUN cat /etc/certs/cert.pem /etc/certs/key.pem > /etc/certs/fullKey.pem

# haproxy-check-api runs on port 8223
EXPOSE 8223

# Run the haproxy binary
CMD ./haproxy-check-api
