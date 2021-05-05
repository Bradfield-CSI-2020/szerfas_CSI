# kvstore
A learn-by-doing implementation of a distributed key-value store

## Usage
Open a terminal window for each desired replica/server, one load balancer, and each desired client, in that order.

For replicas, after building, run each with a unique integer command > 0.

```bash
~/.../kvstore/server @ szerfas-mbp161 (szerfas)
| => go build server.go && ./server 1
kvstore server listening at localhost:8081
```
```bash
~/.../kvstore/server @ szerfas-mbp161 (szerfas)
| => go build server.go && ./server 2
kvstore server listening at localhost:8082
```
```bash
~/.../kvstore/server @ szerfas-mbp161 (szerfas)
| => go build server.go && ./server 3
kvstore server listening at localhost:8083
```

For the load_balancer, pass in `localhost:8080` followed by the listening addresses of each of your replicas.
For example:
```bash
~/.../kvstore/load_balancer @ szerfas-mbp161 (szerfas)
| => go build load_balancer.go && ./load_balancer localhost:8080 localhost:8081 localhost:8082 localhost:8083
```
You'll notice that the load_balancer automatically promotes a replica to a leader.
The leader then uses a semi-synchronous replication scheme -- reaching out synchronously to the first replica and asynchronously to all other replicas upon receiving an update request.

Then simply build and launch the client without any arguments (it's currently hardcoded to connect to localhost:8080, which should be your load_balancer).
```bash
~/.../kvstore/client @ szerfas-mbp161 (szerfas)
| => go build client.go && ./client
```

From the client, you can run commands of two formats: `set key=value` or `get key`. 
You'll notice that `get` commands are passed through the load balancer to replicas in round robin fashion.
`Set` commands are only sent to the leader, and propogated via the semi-synchronous replication scheme mentioned above.

Example operations from the client:
```bash
set stephen=zerfas
Successfully set key 'stephen' with value 'zerfas'

get stephen
zerfas

get stephen
zerfas

get stephen
zerfas

get stephen
zerfas

set val=pal
Successfully set key 'val' with value 'pal'

get val
pal
```



## TODO
* Add input validation and tests. This is currently a brittle and mostly untested implementation; incorrect args can throw the system into bugs.
* Automatic failover -- if the leader terminal is killed via a SIGINT
* Binary protocol -- using JSON and text-based messages was a shortcut
* Graceful handling of interrupts

## Additional ideas 
* Switching from statement-based replication to row-based replication
* Build a more robust on-disk format. This implementation assumes everything can fit in memory.
