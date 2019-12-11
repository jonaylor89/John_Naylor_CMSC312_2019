# CMSC312

### Project Description

In this project, each directory with go code contains a specific part or resource for the operating system simulator. `sched` holds the structures and instructions for the scheduler, `memory` contains code related to physical and virtual memory as well as the cache, etc. etc. `ProgramFiles` contains templates that are available for using while the OS is running. The simulator's front end is a terminal user interface that displays information about the processes running and the memory usage of the system. The user can load program file templates from the TUI which will send requests to the goroutine in charge of adding processes to the appropriate queue, which in turn allocates the appropriate memory as well. From there, the scheduler, running in a seperate goroutine, will pick up processes from those queues and execute them. Processes can run on the cpu, perform io functions, enter the critical section, and communicate with other processes. For interprocess communication, processes are assigned a mailbox at creation which they can store values in or receive from. Processes have pages made for them at their creation that are stored in virtual memory until they need to be accessed in which case they are moved to physical memory. There is a ARC Cache for pages to further speed up memory access.

---------------------------------------

### Building and Running

##### Requirements
- Golang compiler
- GNU make (optional)
- Docker (optional)

##### Building

without docker or make
```sh
go build -o jose .
```

with make
```
make
```

with docker and make
```
make docker-build
```

**[+] Go compiles to a single executable, no linking of libraries necessary**

##### Execution

without docker
```
./jose
```

with docker
```
docker run -it jose:latest
```

**[+] Because the frontend for this is a TUI, it helps to full screen the terminal you're running the simulator in**


# Usage

When executed, the OS Shell is shown. 

```
[os_simulator]$ 
```

The available commands for the shell are:
- load
    - Load in template file and create processes from it
    - e.g. `load ProgramFiles/cpu.prgm 10`
        - load template 1 and create 1000 processes
- exit || quit
    - Exits simulator

# Testing

To execute all tests for the application:

```
~$ make test
```

and for an individual module, just:

```
~$ go test module
```

# Assignment

### Part one


The requirements for project part 1 (deadline October 6th) are as follows:

- [x] having at least 4 of your own program file templates
- [x] having a procedure that reads these program files and generates user-specified number of processes from them (hence randomization of values from templates must be used)
- [x] assigning a PCB to each process that stores basic metadata, including process state
- [x] having a single scheduler that optimizes the process running cycle 
- [x] having a dispatcher that changes the status of each process in real time

All of this must be within a single application, not multiple separate modules.

---------------------

### Part two

The requirements for project part 2 (deadline November 10th) are as follows:

- [x] adding critical sections to your processes (can be implemented e.g., as enclosing selected instruction within critical section tag)
- [x] implementing one selected critical section resolving algorithm (mutex lock / semaphore / monitor)
- [x] adding memory and basic operations on it + taking memory into account when admitting processes into ready state and scheduler

Please remember that these requirements are minimal requirements for C/D grade. Those of you who aim for A/B grades must be aware that these require much more functionalities to be implemented. You are free to submit additional functionalities within project part 2 for evaluation.

---------------------------

### Part three

- [x] Multithreading
- [x] GUI (or TUI)

------------------------

### TODO
- Parent + child
    - pipes
- Two CPUs and Schedulers
    - Load balancer to control which processes go where
    - Critical section for multithreading
- return when IO from process.execute to kernel
- Sorting process table
- Config for switching schedulers
    - `- sched="rr" || sched="fcfs"`
- Kernel go module
    - Wrapper for:
        - Sched
        - Memory
        - CPU
- Add shell prompt to config file
        
### Known Bugs
- Race condition between the scheduler and the tui
    - TUI tries to display processes that have already have been deleted (Big bad)
    - *FIX* I should be a good person and add locks
- Race condition between scheduling algorithm and recvProc worker
    - Patched by going through the processes backwards but should be properly fixed
    - *FIX* Add locks to the processes 
