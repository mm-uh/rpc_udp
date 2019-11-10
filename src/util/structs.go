package util

type RPCBase struct {
	MethodName string
	FirstArg   string
}

type CallerStruct struct {
}

type ResponseRPC struct {
	Response string
	Error    error
}
