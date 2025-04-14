Self-signed TLS certificate and key for testing. Generated with:

openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
