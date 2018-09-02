# cloudwatch

This is a fork of [this library](https://github.com/ejholmes/cloudwatch) which allows treating CloudWatch Log streams as `io.WriterClosers` and `io.Readers`.

## Usage

```go
session := session.Must(session.NewSession(nil))
group := NewGroup(cloudwatchlogs.New(session), "groupName")
w, err := group.Create("streamName")

io.WriteString(w, "Hello World")

r, err := group.Open("streamName")
io.Copy(os.Stdout, r)
```

## Dependencies

This library depends on [aws-sdk-go](https://github.com/aws/aws-sdk-go/).
