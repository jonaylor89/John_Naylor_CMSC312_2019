# CMSC312

The operating system simulation that is built in the class
will live in the repository. 


# Building 

without docker
```sh
~$ make
```

with docker
```
~$ make docker-build
```

# Execution

without docker
```
~$ ./jose
```

with docker
```
~$ docker run -it jose:latest
```

# Usage

When executed, the OS Shell is shown. 

```
OS Shell
-------------------
==> 
```

The available commands for the shell are:
- load
    - Load in template file and create processes from it
    - e.g. load ProgramFiles/temp1.txt 1000
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
- [_] GUI (or TUI)

------------------------

### TODO
- Comment everything
- Critical section for multithreading
- Memory caching
- Parent + child
    - pipes
- Two CPUs and Schedulers
- return when IO from process.execute to kernel
- Config for switching schedulers
    - `- sched="rr" || sched="fcfs"`
- Init methods 
    - CPU
    - Memory
    - Scheduler
