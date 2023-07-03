# net-cat

## Objectives

Net-Cat is the project where we should construct tcp-server in which tcp-clients can chat with each other.


## Usage
- run the server by specifying the port
```console
$ go run .
Listening on the port :8989
$ go run . 8000
Listening on the port :8000
- connect to server
```console
$ nc localhost 8000
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]:
```

### authors

@aomirhan
@rzhampeis