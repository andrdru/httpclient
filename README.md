# httpclient
golang net/http wrapper

## HttpClient

### Examples

```go
client := httpclient.NewHttpClient("jsonplaceholder.typicode.com")
	statusCode, res, err := client.Request(context.Background(), http.MethodGet, "/posts", nil, nil)

	fmt.Println(statusCode, res, err) // 200 [<bytes>] <nil>
```

## JsonClient

### Examples
Post request
```go
client := httpclient.NewJsonClient("jsonplaceholder.typicode.com")
	statusCode, res, err := client.Post(context.Background(), "/posts",
		&map[string]interface{}{
			"test": 101,
		}, nil)

	if r, ok := res.(map[string]interface{}); ok {
		fmt.Println(statusCode, r, err)
	}
```

Returning json object
```go
client := httpclient.NewJsonClient("jsonplaceholder.typicode.com")
	statusCode, res, err := client.Get(context.Background(), "/posts/1", nil, nil)

	if r, ok := res.(map[string]interface{}); ok {
		fmt.Println(statusCode, r, err) // 200 <map> <nil>
	}
```

Returning json array
```go
client := httpclient.NewJsonClient("jsonplaceholder.typicode.com")
	statusCode, res, err := client.Get(context.Background(), "/posts", nil, nil)

	if r, ok := res.([]interface{}); ok {
		fmt.Println(statusCode, r, err) // 200 <slice> <nil>
	}
```
