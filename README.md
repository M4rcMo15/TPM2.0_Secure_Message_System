# TPM 2.0 Secure Message Board System

## How to use

```bash
# run the server
go run cmd/server/main.go

# run the client
go run cmd/client/main.go
```
The client connects to the server and prints the authentication result.

- Implement a simple message board to which clients can post messages
- Clients authenticate to the board with their TPM
  - Register on the first interaction
    - At this point only the server needs to be authenticated (e.g., known certificate)
  - Further communication is fully authenticated
    - Both sides need to be authenticated (the server does not need to use TPM)
    - The client should not be able to connect from a different device after the registration (e.g., remote attestation

# build all packages in cmd directory
go build -o bin/ cmd/...

# run the compiled binary
./bin/poc

go build -o bin/client/ cmd/client/main.go && sudo bin/client/main

go build -o bin/server/ cmd/server/main.go && bin/server/main
```

<!-- links -->

[1]: https://decide.nolog.cz/#/poll/LjjKu4fayG/participation?encryptionKey=yagKvnA75GLrktCthJX4xKa08pgmN32CLtWDg0Tw
