# Command-line Interface

A command-line interface (CLI) is a text-based interface that allows users to interact with a computer's operating system by entering commands in the form of text. In our project, the CLI will be used to send and receive messages on a message board.
To implement the client-side functionality of the message board, we will use a CLI that allows users to perform various actions, such as posting messages. The CLI will prompt the user to enter their message, and then send it to the server for posting. To create the CLI, we will use the bufio library. This library is simple to implement and offers a wide range of functionality. In the future, we may also use the CLI to handle authentication and registration. The CLI could also be used to change the destination of a message.
In summary the CLI will provide a simple and efficient way for users to interact with the message board.

# HTTPS

In our  project, the server will use the HTTPS protocol to listen for incoming requests from clients and send responses back to the clients. When a client wants to post a message to the message board, it will send an HTTPS request to the server using a specific method, such as GET, POST, or PUT. The request may also include additional information, such as headers and a body, which can be used to provide additional context or data for the request.
Upon receiving the request. The server will then send an HTTPS response back to the client, indicating whether the request was successful or not. 
Overall, the use of HTTPS in our message board project will provide a simple and flexible way for the client and server to communicate with each other. It will allow the client to send requests to the server, and the server to send responses back to the client. 

# Contribution

During the second phase of the project, I contributed by dividing the code between the client and the server, and I began to create the CLI interface.

