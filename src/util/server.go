package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/sirupsen/logrus"
)

func CallMethod(i interface{}, methodName string, args []interface{}) (interface{}, error) {
	logrus.Info("Calling method dynamically")

	var ptr, value, finalMethod reflect.Value

	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
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
		for _, arg := range args {
			params = append(params, reflect.ValueOf(arg))
		}
		return finalMethod.Call(params)[0].Interface(), nil
	}

	// return or panic, method not found of either type
	return "", errors.New("error calling method")
}

func HandleOptions(pc net.PacketConn, addr net.Addr, buf []byte, n int) {
	fmt.Println("Handle options")

	var requestRPCCall RPCBase
	err := json.Unmarshal(buf[:n], &requestRPCCall)
	if err != nil {
		logrus.Error("Error unMarshalling")
	}

	// Handle rpc options
	response, err := CallMethod(&CallerStruct{}, requestRPCCall.MethodName, []interface{}{requestRPCCall.FirstArg})
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

func ListenServer(addr string) {
	logrus.Info("Listen server at " + addr)
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		logrus.Fatal(err)
	}
	defer pc.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue

		}
		go HandleOptions(pc, addr, buf, n)

	}

}
