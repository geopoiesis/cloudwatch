# cloudwatch

[![Godoc](https://godoc.org/github.com/marcinwyszynski/cloudwatch?status.svg)](http://godoc.org/github.com/marcinwyszynski/cloudwatch)
[![CircleCI](https://circleci.com/gh/marcinwyszynski/cloudwatch/tree/master.svg?style=svg)](https://circleci.com/gh/marcinwyszynski/cloudwatch/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/marcinwyszynski/cloudwatch)](https://goreportcard.com/report/github.com/marcinwyszynski/cloudwatch)
[![codecov](https://codecov.io/gh/marcinwyszynski/cloudwatch/branch/master/graph/badge.svg)](https://codecov.io/gh/marcinwyszynski/cloudwatch)

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
