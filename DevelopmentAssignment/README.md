# DaSH Lab Development Assignment

## Install and Configure
Requirements:
- Go > 1.12
- Docker > 18.1
- Docker Compose > 1.17

To Run:
1. Install the dependencies 
2. Add API Keys (Can be done in docker-compose.yml, it is NOT secure\. This is due to time constraints.)
3. Add an 'input.txt' file to ./client with prompts.
4. Call:
```
$ chmod +x ./up.sh
$ ./up.sh
```

** The Build Process is Dockerized, and any changes can be built by running: **
```
$ ./up.sh
```

## Architecture
The most appropriate way to explain how the program is structured is to run along a request,
explaining design choices, and limitations along the way. The Hugging Face serverless inference API
was used for the gemma 2B model.

Introduction
The main choice is the medium of communication between the Client and Server; Here I have chosen D-Bus
mostly as an academic exercise, but also because it presents a unique set of challenges and opportunities
leading to a more interesting architecture.
This program is designed with a particular scenario in mind, when 'n' number of Clients connect to the Server
all with different input files but, write to a common 'output.json' and duplication must be avoided.

    Client Entry
1. The Client main() function starts off by connecting over TCP to the Session bus of the Server container.
2. Then, prompts are read from the input.txt file.
3. These prompts are then hashed using the adler32 checksum. (This is not perfect, there is a small chance
    that two prompts may generate the same hash which will lead to the Server ignoring one of them.)
4. Next, the RegisterHashes(hashes) method is called on the Server over D-Bus.

    RegisterHashes [Concurrent]
5. The Server then subtracts the set of hashes that it has already processed from the set of hashes sent by the Client,
    any remaining hashes are unique.
6. The Client is assigned a Unique ID (Unix UTC Millisecond time stamp of request as it is unique). Better alternatives exist,
    but were ignored.)
7. The ID and list of unique hashes is sent back to the Client.

    Back to Client
8. The Client creates and sends requests (containing the prompts) corresponding to the unique hashes sent back by the Server.
9. The Client then goes into a blocked state, waiting for the Server to send a write-to-file command.

    LLM Request [Concurrent]
10. The Server sends API requests to the LLM for each received prompt.
11. Assigns the output to one of the connected Clients.
12. Then, the output is broadcast via a D-Bus interface Signal to _all_ connected Clients.
13. The function then goes into blocked state for one of three events, Success (WriteCheck), Failure or a 15 second timeout. In cases other
    than Success, the output is re-assigned.

    Signal Received
14. On receiving a signal, the Client checks if the output is assigned to it, if not, the output is discarded.
15. If assigned, the output is written to the output.json.
16. The D-Bus interface's WriteCheck method is called.
