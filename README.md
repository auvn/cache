##
Desired features:
- [x] Key-value storage with string, lists, dict support
- [x] Per-key TTL
- [x] Operations:
  - Get
  - Set
  - Update
  - Remove
  - Keys
- [x] Custom operations(Get i element on list, get value by key from dict, etc)
- [x] Golang API client
- [x] Telnet-like/HTTP-like API protocol

- [x] Provide some tests, API spec, deployment docs without full coverage, just a few cases and some examples of telnet/http calls to the server.

Optional features:
- [x] persistence to disk/db
- [x] scaling(on server-side or on client-side, up to you)
- [x] auth
- [x] perfomance tests

## Server

### Building and running
```
$ wget https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz -O ./go.tar.gz
$ tar -xf go.tar.gz -C ./
$ export GOROOT="$(pwd)/go"
$ export GOPATH="$(pwd)"
$ export PATH="$PATH:$(pwd)/go/bin:$(pwd)/bin"
$ go get -v github.com/auvn/go.cache
$ go.cache -http :1235 -telnet :1234
2017/02/28 21:48:43 serving telnet at: :1234
2017/02/28 21:48:43 serving http at: :1235
```

### Usage

```
Usage of cache:
  -http string
        Address to listen http on. Optional.
  -journal string
        Journal file for persistence. Optional.
  -pass string
        Password for cache authentication. Optional.
  -telnet string
        Address to listen telnet on (default "0.0.0.0:1234")
```

### Examples

#### Telnet

```
$ echo "A3\r\nV3\r\nSET\r\nV6\r\n界世\r\nV3\r\n界\r\n" | netcat localhost 1234
B1

$ echo "A5\r\nV5\r\nLPUSH\r\nV6\r\nMyList\r\nV2\r\n32\r\nV3\r\n256\r\nV2\r\n\r\n\r\n" | netcat localhost 1234
I3

$ echo "A4\r\nV6\r\nLRANGE\r\nV6\r\nMyList\r\nI0\r\nI-1\r\n" | netcat localhost 1234
A3
V2


V3
256
V2
32

$ echo "A1\r\nV4\r\nKEYS\r\n" | netcat localhost 1234
A2
V6
MyList
V6
界世

$ echo "A4\r\nV6\r\nLRANGE\r\nV6\r\n界世\r\nI0\r\nI-1\r\n" | netcat localhost 1234
Eaccessing a key holding the wrong type of value

```

#### HTTP


```
$ echo "A3\r\nV3\r\nSET\r\nV6\r\n界世\r\nV3\r\n界\r\n" | curl http://localhost:1235/ --data-binary @-
B1

$ echo "A1\r\nV4\r\nKEYS\r\n" | curl http://localhost:1235/ --data-binary @-
A2
V6
MyList
V6
界世

$ echo "A4\r\nV6\r\nLRANGE\r\nV6\r\n界世\r\nI0\r\nI-1\r\n" | curl http://localhost:1235/ --data-binary @-
Eaccessing a key holding the wrong type of value

```

Also it's possible to specify a command name using url path and pass arguments in a body:

```
$ echo "A4\r\nV6\r\nMyList\r\nV2\r\n32\r\nV3\r\n256\r\nV2\r\n\r\n\r\n" | curl http://localhost:1235/lpush --data-binary @-
I3

$ echo "A3\r\nV6\r\nMyList\r\nI0\r\nI-1\r\n" | curl http://localhost:1235/lrange --data-binary @-
A3
V2


V3
256
V2
32
```

## Golang client

### Examples

```golang
package main

import (
    "log"
    "time"
    // ideally it should be separated from the main cache's server source
    "github.com/auvn/go.cache/client"
)

func main() {
    c := client.New(&client.Options{
        Auth:        "password",
        PoolSize:    10,
        DialTimeout: 5 * time.Second,
        // if >1 addrs specified - the client will use a hash function (e.g. hash(key) % len(Addrs))
        // to determine what server should be used for a request
        Addrs:       []string{"localhost:1234"},
    })

    isSet, err := c.Set("key", []byte("some value")).Bool()
    // ...

    listSize, err := c.RPush("mylist", []byte("value1"), []byte("value2")).Int()
    // ...

    isSet, err := c.HSet("myhash", []byte("hashKey"), []byte("hashValue")).Bool()
    // ...
}
```

### Performance tests

Tests are done using b.RunParallel and client implementation.

First run:

```
BenchmarkClientSet-4               50000             28075 ns/op             475 B/op         26 allocs/op
BenchmarkClientGet-4               50000             26822 ns/op             484 B/op         26 allocs/op
BenchmarkClientSetGet-4            20000             55914 ns/op             975 B/op         52 allocs/op
BenchmarkClientSetDel-4            30000             55236 ns/op             878 B/op         50 allocs/op
BenchmarkClientLPush10-4           30000             45730 ns/op            1424 B/op         63 allocs/op
BenchmarkClientLPush50-4           10000            101813 ns/op            5816 B/op        223 allocs/op
BenchmarkClientLPush100-4          10000            177830 ns/op           10950 B/op        423 allocs/op
BenchmarkClientRPush10-4           30000             46899 ns/op            1424 B/op         63 allocs/op
BenchmarkClientRPush50-4           10000            104658 ns/op            5815 B/op        223 allocs/op
BenchmarkClientRPush100-4          10000            188094 ns/op           10954 B/op        423 allocs/op
BenchmarkClientRPop-4              50000             28659 ns/op             483 B/op         26 allocs/op
BenchmarkClientLPop-4              50000             28955 ns/op             483 B/op         26 allocs/op
BenchmarkClientLPushLPop-4         20000             63522 ns/op            1008 B/op         54 allocs/op
BenchmarkClientLPushRPop-4         20000             62626 ns/op            1008 B/op         54 allocs/op
BenchmarkClientRPushRPop-4         20000             60984 ns/op            1008 B/op         54 allocs/op
BenchmarkClientRPushLPop-4         20000             61045 ns/op            1007 B/op         54 allocs/op
BenchmarkClientHSet-4              50000             31182 ns/op             580 B/op         30 allocs/op
BenchmarkClientHGet-4              50000             30439 ns/op             652 B/op         30 allocs/op
BenchmarkClientHSetHGet-4          20000             65649 ns/op            1249 B/op         60 allocs/op
BenchmarkClientHSetHDel-4          20000             64515 ns/op            1103 B/op         58 allocs/op
```

Second run (cache server is not restarted):

```
BenchmarkClientSet-4               50000             34677 ns/op             474 B/op         26 allocs/op
BenchmarkClientGet-4               50000             28859 ns/op             484 B/op         26 allocs/op
BenchmarkClientSetGet-4            20000             59364 ns/op             974 B/op         52 allocs/op
BenchmarkClientSetDel-4            20000             65293 ns/op             879 B/op         50 allocs/op
BenchmarkClientLPush10-4           30000             47621 ns/op            1433 B/op         63 allocs/op
BenchmarkClientLPush50-4           20000             99319 ns/op            5824 B/op        223 allocs/op
BenchmarkClientLPush100-4          10000            133432 ns/op           10956 B/op        423 allocs/op
BenchmarkClientRPush10-4           30000             74067 ns/op            1434 B/op         63 allocs/op
BenchmarkClientRPush50-4           20000            113246 ns/op            5825 B/op        223 allocs/op
BenchmarkClientRPush100-4          10000            137455 ns/op           10957 B/op        423 allocs/op
BenchmarkClientRPop-4              50000             32630 ns/op             483 B/op         26 allocs/op
BenchmarkClientLPop-4              50000             34574 ns/op             484 B/op         26 allocs/op
BenchmarkClientLPushLPop-4         10000            115599 ns/op            1010 B/op         54 allocs/op
BenchmarkClientLPushRPop-4         20000             63410 ns/op            1007 B/op         54 allocs/op
BenchmarkClientRPushRPop-4         20000             68052 ns/op            1008 B/op         54 allocs/op
BenchmarkClientRPushLPop-4         20000             63370 ns/op            1007 B/op         54 allocs/op
BenchmarkClientHSet-4              50000             31886 ns/op             580 B/op         30 allocs/op
BenchmarkClientHGet-4              50000             32150 ns/op             652 B/op         30 allocs/op
BenchmarkClientHSetHGet-4          20000             64285 ns/op            1248 B/op         60 allocs/op
BenchmarkClientHSetHDel-4          20000             64192 ns/op            1103 B/op         58 allocs/op
```

## API spec

### Protocol

There was implemented a protocol similar to [RESP](https://redis.io/topics/protocol). The protocol supports the following types:
* Int:

    ```
    I256\r\n
    I-234\r\n
    ```

* Value (slice of bytes)

    ```
    V10\r\n1234567890\r\n
    V0\r\n
    V6\r\n界世\r\n
    V2\r\n\r\n\r\n
    ```

* Array

    ```
    A0\r\n
    A3\r\nV6\r\nEXPIRE\r\nV3\r\nkey\r\nI10\r\n
    A2\r\nI1\r\nI2\r\n
    ```

* Bool

    ```
    B0\r\n
    B1\r\n
    ```

* Error

    ```
    Ewrong number of arguments\r\n
    Eunknown command\r\n
    ```
* Null

    ```
    N\r\n
    ```

### Available commands
Examples were made by using telnet util.

#### AUTH pass
The server supports simple authorization (-pass flag with non-empty string value). By using this command a client able to perform authentication on the server.

Example:
```
A1
V4
KEYS
Eauth required

A2
V4
AUTH
V2
11
Eforbidden

A2
V4
AUTH
V10
HardPassWD
B1

A1
V4
KEYS
A0
```

#### KEYS
Prints all stored keys in the cache.

Example:

```
A1
V4
KEYS

A1
V3
KeY
```

#### EXPIRE key seconds
Sets key's TTL.

Example:

```
A3
V6
EXPIRE
V3
KeY
I10

B1
```

#### DEL [keys...]
Deletes specified keys.

Example:

```
A2
V3
DEL
V5
12345

I1
```

#### TTL key
Returns key's TTL value if its set.

Example:

```
A2
V3
TTL
V3
KeY

I87
```

#### SET key value
Sets the value to the specified key

Example:

```
A3
V3
SET
V3
KeY
V3
val

B1
```

#### GET key
Gets a value of the specified key

Example:

```
A2
V3
GET
V3
KeY

V3
val
```

#### LPUSH key [values...]
Prepends values to a list with the specified key

Example:

```
A4
V5
LPUSH
V6
MyList
V2
12
V3
345

I2
```


#### RPUSH key [values...]
Appends values to a list with the specified key

Example:

```
A3
V5
RPUSH
V6
rpush_
V2
12

I1
```

#### LPOP key
Removes and gets the first element in a list with the specified key.

Example:

```
A2
V4
LPOP
V6
rpush_

V2
12
```

#### RPOP key
Removes and gets the last element in a list with the specified key.

Example:

```
A2
V4
RPOP
V6
MyList

V5
VALUE
```

#### LRANGE key start stop
Gets a range of elements from a list with the specified key.

Example:

```
A4
V6
LRANGE
V6
MyList
I0
I1

A2
V3
345
V2
12

```

#### LINDEX key index
Gets an element from a list by its index.

Example

```
A3
V6
LINDEX
V6
MyList
I0

V3
345
```

#### HSET key hashKey value
Sets the value of the specified hashKey in hash with the key.

Example:

```
A4
V4
HSET
V4
hash
V8
hash_key
V10
hash_value

B1
```

#### HGET key hashKey

Gets the value of the specified hashKey from a hash with the key.

Example:

```
A3
V4
HGET
V4
hash
V8
hash_key

V10
hash_value
```


#### HDEL key [hashKeys...]
Deletes the specified hashKey in a hash with the key.

Example:

```
A3
V4
HDEL
V4
hash
V8
hash_key

I1
```


#### HKEYS key
Prints all keys in a hash with the key.

Example:

```
A2
V5
HKEYS
V4
hash

A1
V8
hash_key
```
