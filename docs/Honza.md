# Trusted Platform Module

Trusted Platform Module is a trusted entity with a purpose of improving the security of its user. It is also referred to as an international standard. It usually takes on a form of a dedicated microdevice set on the motherboard, however other forms of TPMs are common, such as Integrated TPMs or Virtual and Software TPMs. 

It can be used as a secret storage and it provides a true/hardware random number generator. The TPM can be used for remote attestation, absently checking whether the hardware and software setup complies with the according demands. Two versions of TPMs emerged. Version 1.2 and 2.0, which is not backward compatible with the former.

## TPM as part of our project

Our goal is to be able to authenticate one or more clients to a server via a client’s TPM. And such as that it is a vital element to our project. The server then checks and authorizes the client to add something to the board or take a peek at the board, or manipulate the board in unpredictable and peculiar ways. As far as it is allowed by the server. The communication should be therefore binded to a concrete device via the TPM, and without using a registered device, access to the board shall not be granted.

On the very first interaction a registration of the client should take place. Afterwards the communication between the server and the client should be bidirectionally authenticated, although only clients use TPM. This is done by TPM holding a set of secret keys.

TPM holds two types of keys, Endorsement Keys and Attestation Keys. Attestation Keys are certified by Endorsement Keys and then used to directly sign or encrypt data. Each TPM has its own unique identity defined by a set of Endorsement Keys that are provisioned by its manufacturer. These keys are not used for data encryption. It should be noted that none of the keys ever leave the device.

Our choice of programming language Golang was heavily influenced by the libraries and projects widely working with TPMs being available to significantly simplify the coding efforts. The most notable is “https://github.com/google/go-attestation”, which states that the project tries to provide high level primitives for both client and server logic.

In summary TPMs are being forced where possible and it seems to be a good thing.

## Contributions

I took part in engineering architecture, doing research and finding appropriate solutions. I set up a basic html server structure, which in the end (I think) served only as an inspiration. I also participated in managing, merging and cleaning the codebase, dividing the workload among the team members and writing the final TPM report article.
