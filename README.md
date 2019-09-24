# CMSC312

The operating system simulation that is built in the class
will live in the repository. 


# Building 

without docker
```sh
~$ go build -o OS
```

with docker
```
~$ docker build -t jonaylor/operatingsystem:os .
```

# Exectuion

without docker
```
~$ ./OS
```

with docker
```
~$ docker run -it jonaylor/operatingsystem:os
```

# Assignment

### Part one


The requirements for project part 1 (deadline October 6th) are as follows:

+ having at least 4 of your own program file templates
+ having a procedure that reads these program files and generates user-specified number of processes from them (hence randomization of values from templates must be used)
+ assigning a PCB to each process that stores basic metadata, including process state
+ having a single scheduler that optimizes the process running cycle 
+ having a dispatcher that changes the status of each process in real time

All of this must be within a single application, not multiple separate modules.

---------------------