package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/sirupsen/logrus"
)

// Define the RPC server struct
type Server struct {
	// This struct is passed as argument, and we made reflection over this interface to call the methods
	caller interface{}
	// If the server get any error handle the information, is represented here
	Error error
	// Represent the addr of the server
	addr string
}

// Create a new instance of *Server struct
func NewServer(i interface{}, addr string) *Server {
	return &Server{caller: i, addr: addr}
}

// Call methods associated to Server.caller using reflection
func (server *Server) callMethod(methodName string, args interface{}) (interface{}, error) {
	logrus.Info("Calling method dynamically")

	var ptr, value, finalMethod reflect.Value

	value = reflect.ValueOf(server.caller)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(server.caller))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if finalMethod.IsValid() {
		params := make([]reflect.Value, 0)
		//for _, arg := range args {
		params = append(params, reflect.ValueOf(args))
		//}
		return finalMethod.Call(params)[0].Interface(), nil
	}

	// return or panic, method not found of either type
	return "", errors.New("error calling method")
}

// Handle the rpc call of methods
func (server *Server) handleOptions(pc net.PacketConn, addr net.Addr, buf []byte, n int) {
	fmt.Println("Handle options")

	var requestRPCCall RPCBase
	err := json.Unmarshal(buf[:n], &requestRPCCall)
	if err != nil {
		logrus.Error("Error unMarshalling")
	}

	// Handle rpc options
	response, err := server.callMethod(requestRPCCall.MethodName, requestRPCCall.Args)
	if err != nil {
		logrus.Error("Couldn't call method " + err.Error())
		return
	}

	str, err := json.Marshal(&ResponseRPC{Response: response.(string), Error: nil})
	if err != nil {
		_, err = pc.WriteTo(str, addr)
		logrus.Error("Couldn't marshal response for rpc " + err.Error())
		return
	}

	_, err = pc.WriteTo(str, addr)
	if err != nil {
		logrus.Error("Couldn't send buffer " + err.Error())
		return
	}
	logrus.Info("Successfully handled")
}

// Listen in server.addr for calls of the rpc
// Recive as params a channel that handle when we get a error
func (server *Server) ListenServer(exit chan bool) {
	logrus.Info("Listen server at " + server.addr)
	pc, err := net.ListenPacket("udp", server.addr)
	if err != nil {
		logrus.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			server.Error = err
			break
		}
		go server.handleOptions(pc, addr, buf, n)

	}
	exit <- true
}
