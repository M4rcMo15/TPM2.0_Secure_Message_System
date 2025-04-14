# Architecture

## Database

**not necessary for PoC, implement later**

- messages
- users
- stuff for auth

- SQLite (or PostgreSQL in Docker)

## HTTP

- Use [echo](https://echo.labstack.com/docs/quick-start) for the server
- I don't know what to use for the client
- some links: https://www.digitalocean.com/community/tutorials/how-to-make-an-http-server-in-go, https://go.dev/doc/articles/wiki/

## TPM emulator

- link: https://github.com/stefanberger/swtpm 

## CLI

- I don't know
- The prorotype doesn't need to have a CLI necessarily, the finished product
  does

## TPM

[go-attestation](https://github.com/google/go-attestation) library for google seems to do just what we want:

> TPMs can be used to identify a device remotely and provision unique per-device hardware-bound keys.

Also, I found this article which seems to describe what we need to do.

> Usually, TPM vendors provision the TPM device with a Primary Endorsement Key (PEK), and generate a certificate for that key (EK cert) in x.509 format at the time of manufacture. The EK Certificate contains the public part of the PEK, as well as other useful fields, such as TPM manufacturer name, part model, part version, key id, etc. This information can be used to uniquely identify the TPM and if the device OEM securely attaches a TPM to the device it can be used as a device identifier. [tpm2-software.github.io/tpm2-tss][1]

We may just need to use TPM once for the inital authentication and then
establish other keys (or just a cookie) for all subsequent messages, this way
we don't have to work that much with TPM.

[1]: https://tpm2-software.github.io/tpm2-tss/getting-started/2019/12/18/Remote-Attestation.html
