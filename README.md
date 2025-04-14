# PV204-Security-Technologies-Project

## How to use

```bash
# run the server
go run cmd/server/main.go

# run the client
go run cmd/client/main.go
```

The client connects to the server and prints the authentication result.

## Meetings

### 2023-03-15

- Adam and Alexandre work on TPM
- Honza and MarcMo work on HTTP client/server communication

We agreed on the following:

- [Go](https://go.dev/learn/) for implementation
- HTTP protocol for sending messages to / from client / server
- TLS is not necessary for a prototype
- Database is not necessary for a prototype (keep all data in memory)
- We write the report in Markdown format on GitHub
- If we get stuck or have the work done early, concact the other group for / to
  help

We dedice on next meeting by voting in the [poll][1]. Most of the work should
be done on wednesday.

## Requirements

- Implement a simple message board to which clients can post messages
- Clients authenticate to the board with their TPM
  - Register on the first interaction
    - At this point only the server needs to be authenticated (e.g., known certificate)
  - Further communication is fully authenticated
    - Both sides need to be authenticated (the server does not need to use TPM)
    - The client should not be able to connect from a different device after the registration (e.g., remote attestation

## Phases

### II (24. 3. 2024)

Deliver 3-4 page report with:

  - description of TPM
  - software architecture
  - progress with individual contribution

### III (14. 4. 2024)

Deliver:

- presentation includes:
  - issues
  - solutions
  - demo
  - individual contribution

- final product, we need to finish:
  - sending messages
  - reading messages
  - CLI on the client

Roadmap:

- finish code and presentation during the week
- record presentation on the weekend (once the code is ready)

Schedule:

- Alex: Mon-Wen
- Adam: Mon, Wen
- Honza: Fri-Sun
- Marc: ?

### IV (15. 5. 2024)

TBD

## Basics of Go

```bash
# build a package from sources
go build -o bin/poc cmd/poc/

# build all packages in cmd directory
go build -o bin/ cmd/...

# run the compiled binary
./bin/poc

go build -o bin/client/ cmd/client/main.go && sudo bin/client/main

go build -o bin/server/ cmd/server/main.go && bin/server/main
```

<!-- links -->

[1]: https://decide.nolog.cz/#/poll/LjjKu4fayG/participation?encryptionKey=yagKvnA75GLrktCthJX4xKa08pgmN32CLtWDg0Tw
