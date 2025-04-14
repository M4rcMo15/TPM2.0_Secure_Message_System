## Threat Model

The security of the communication stands on several pillars:

    • TMP keys on the client (authentication of the client)
    • TLS keys on the server (secure channel, authentication of the server)
    • TLS certificate signed by CA (secure channel, authentication of the server)
    • Secret used for cookies encryption (session post-TMP authentication)

An attacker may compromise the security of the system, if:

    • The server’s TLS key is compromised and the attacker is able to read the messages, including the authentication cookies.
    • The certificate authority’s TLS keys are compromised – the implications are the same as above
    • The cookies for a single session are compromised – an attacker may act in the name of the logged-in user
    • Server’s secret for encryption of cookies is compromised – an attacker is able to act in the name of any logged-in user

We don’t consider these risks to be significant enough to make us reconsider our design choices. We believe that similar risks would be present in any system. TLS is used universally on the web and we consider them to be a solid foundation to build on. Cookies are rotated every login, and TLS certificate on the server has a validity for 90 days. 

## Architecture and Design Choices

### Registration process

The client initiates a new connection to the server and begins the act of registration. The server sends a challenge to the client. If a client succeeds at the challenge, the registration is successfull and the EK's public key is bound to the username.

### Authentication

For any subsequent authentication request, the server takes advantage of the EK's public key saved during the registration and uses it to generate a new challenge for the client. The client proves their identity by decrypting the challenge. Subsequent communication is authenticated via HTTP cookies and secured by TLS.

## State of the project

We have successfully implemented registration of the client to the server through remote attestation via HTTP protocol. The server does not use TLS yet and data for the registration are stored in-memory instead of pernament storage, such as a database. We have a proof of concept code for the command line interface. 

We believe that the hardest part of the project is behind us and the missing features, such as sending and recieving messages, will be much easier to implement.

## Contributions

I participated in creating the architecture and set up the repository structure and basic instructions on how to use Golang. I have implemented a proof-of-concept of TMP remote attestation. Later I expanded the code to encode the client-server communication in JSON files as a prerequisite for sending them through the Internet.

I have also extracted the common code into libraries, and finished the implementation of registration of a client using HTTP requests and cookies. I have also put together individual contributions into the final report.
