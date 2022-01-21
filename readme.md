# vRPC

vrpc is a rpc based on lemon-mint's [vstruct](https://github.com/lemon-mint/vstruct).

## install

```bash
go get github.com/snowmerak/vrpc
```

## Frame

vrpc has a frame for communication.

|Service|Method|Sequence|BodySize|Body|
|---|---|---|---|---|
|uint32|uint32|uint32|uint32|bytes|

Body size can't over then `BodySize`.  
`BodySize`'s limitation equals `2^32-1`.

`Sequence` is set by random value.  
response's `Sequence` is set by request's `Sequence` + 1.

## Server

### new server

```go
server := vrpc.NewServer(log.New(os.Stdout, "server: ", log.LstdFlags))
```

Calling `NewServer`, create vrpc server instance.  
vrpc server must get a `*log.Logger` instance.

### example vstruct

```vstruct
struct Request {
    string Name;
}

struct Response {
    string Reply;
}

```

You can compile this sample vstruct code from [vstruct repo](https://github.com/snowmerak/lemon-mint/vstruct).

### register handler

```go
server.Register(1, 1, func(_ vrpc.EmptyValue) test.Response {
	response := test.New_Response("Hello" + "!")
	return response
})

server.Register(1, 2, func(request test.Request) test.Response {
	response := test.New_Response("Hello, " + request.Name() + "!")
	return response
})
```

We can register handler to server calling `Register` method.  
First parameter is `Service` number.  
Second parameter is `Method` number.

You must write function with one vstruct parameter and one vstruct return.  
Function's body is free for you.

vrpc offers `vrpc.EmptyValue` for empty parameter.  
can not use `nil`.

### serve

```go
if err := server.Serve("localhost:8080"); err != nil {
	panic(err)
}
defer server.Shutdown()
```

## client

### new client

```go
client, err := vrpc.NewClient("localhost:8080", log.New(os.Stdout, "client: ", log.LstdFlags))
if err != nil {
	panic(err)
}
defer client.Close()
```

Calling `NewClient`, create vrpc client.  
`Client` must have server's address:port and `*log.Logger`.

### request to server

```go
rs, err := client.Request(1, 1, vrpc.Empty())
if err != nil {
	panic(err)
}
fmt.Println(test.Response(rs).Reply())
```

client's `Request` method request with `Service`, `Method`, vstruct instance.  
`Request`'s return type is `[]byte`, you can use this for vstruct by wrapping `[]byte`.

```go
rs, err = client.Request(1, 2, test.New_Request("snowmerak"))
if err != nil {
	panic(err)
}
fmt.Println(test.Response(rs).Reply())
```

You can use client many times.

# Thanks

[@lemon-mint](https://github.com/lemon-mint)
