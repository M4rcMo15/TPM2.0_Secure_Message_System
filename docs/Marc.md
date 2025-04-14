# Separation into client and server

In our basic implementation, the communication between client and server is based on a secure data exchange, where the client initiates the process by generating and sending its public key and TPM attestation parameters to the server. The server, after receiving this data, generates a cryptographic challenge that only the legitimate client device can solve, thanks to its TPM. 

This strategy ensures that each party fulfills its specific role: the client proves its authenticity and the server validates this authenticity before allowing access or interaction. 

In our system, we use an HTTP server to handle communications between the client and the server, with data exchange taking place via POST requests. The data, including the client's public key and attestation parameters, is serialized to JSON for sending, and the server deserializes this JSON to process the information received. 
This process allows for a structured and efficient data exchange. 

The server then responds to the client also sending a generated cryptographic challenge that the client must solve, thus proving the authenticity of its identity and the integrity of its TPM.

## Contributions

I have contributed to the development of a secure registration and authentication mechanism within our TPM message board project. Drawing inspiration from a proof of concept initially programmed by, Adam, I dedicated myself to expanding upon these ideas. 

My primary role involved architecting and implementing a clear separation of responsibilities between client and server components. 

This development entailed establishing a robust and functional communication protocol over HTTP, which facilitated efficient data exchange and also granted the  standards of data integrity. 

Through careful programming and testing, I ensured that our system manages succesfully TPM functionalities using its potential, thereby enhancing the security and reliability of our authentication processes. This work presents a good advancement in our project, setting a solid foundation for future expansions and enhancements to the project.
 