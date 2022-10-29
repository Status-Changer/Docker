# Docker

## A Simple Docker Container Implementation in Golang

### Suggested running environment:
- Operating system: Ubuntu 14.04 *(Exactly)*
- Core Version: 3.10.0-83-generic
- Golang 1.7.1 *(at least)*

### Deployment Steps:
1. Clone the source code `git clone https://github.com/Status-Changer/Docker.git`;
2. If you are using a Linux machine for development, go to step 4; or
3. Sync your code to a Linux environment first;
4. `cd` to the code's root directory, then execute `go build -o docker` to build the project;
5. Run the executable file like the REAL docker, have a good time enjoying this project!

### Parameters Supporting Now:
- `run` runs a container with the following sub-parameters:
    - `-ti` running at interactive mode
    - `-m` memory space limitation
    - `-cpushare` CPU share limitation
    - `-cpuset` CPU set limitation
    - `-v` mounting a host directory to a container directory
- `commit` zips a container image to a `.tar` file
- `ps` displays all containers