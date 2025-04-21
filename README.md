# TPM2.0_Secure_Message_System
## Server setup

### PostgreSQL

The server requires PostgreSQL database. You can setup a development instance
with Docker:

```bash
docker run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 postgres:bookworm

# or with Podman
podman run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 docker.io/postgres:bookworm
```

### Credentials

Credentials to the database are stored in the environemnt file `.env`. TLS
certificate and key paths are also configured there.

```bash
# only for development, insecure!
cp .env.example .env
```

The self-signed TLS certificate and key in the `tls` directory are provided for
testing purposes and were generated with the following:

```bash
openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes
```

## Client setup

The client requires
[TPM](https://wiki.archlinux.org/title/Trusted_Platform_Module).

You can disable TLS verification by running the client with the
`IGNORE_TLS_ERRORS=1` environment variable.

```
$ export IGNORE_TLS_ERRORS=1
$ sudo -E ./bin/client
```

## Compilation && Running the tool

The source code is written in [Go](https://go.dev/doc/install). Here's how to
build and run the binaries:

```bash
# compile both binaries with make
make build

# or manually
go build -o bin/server cmd/server/main.go
go build -o bin/client cmd/server/client.go

# run the server
go run cmd/server/main.go

# run the client
go run cmd/client/main.go
```

## Development

```bash
# build binaries
make build

# format the code
make ftm
```

## Credits and Contribution
    • Adam Chovanec
    • Alexandre Deneu
    • Jan Sekanina
    • Marc Monfort Muñoz
