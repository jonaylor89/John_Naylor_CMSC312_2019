# CMSC312

The operating system simulation that is built in the class
will live in the repository. 

# Available Commands

### new 
create new process

### exit
exit simulator


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
