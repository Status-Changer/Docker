<font size=5>**Docker**:</font> <font size=3>A Simple Container Implementation in Golang</font>

### Suggested running environment:
- Operating System: Ubuntu 14.04 *(Exactly)*
- Core Version: 3.10.0-83-generic
- Golang 1.7.1 *(at least)*

### Deployment Steps:
1. Clone the source code `git clone https://github.com/Status-Changer/Docker.git`;
2. If you are using a Linux machine for development, go to step 4; or
3. Sync your code to a Linux environment first;
4. `cd` to the code's root directory, then execute `go build -o docker` to build the project;
5. Run the executable file like the REAL docker, and have a good time enjoying it!

### Parameters Supporting Now:
- `run` runs a container with the following (optional) sub-parameters:
    - `-ti` running at interactive mode
    - `-m` memory space limitation
    - `-cpushare` CPU share limitation
    - `-cpuset` CPU set limitation
    - `-v` mounting a host directory to a container directory
    - `-d` running in the background
- `commit` zips a container image to a `.tar` file
- `ps` displays all containers
