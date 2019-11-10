package main

import (
	"github.com/stdevMac/rpc_udp/src/util"
)

func main() {
	// listen to incoming udp packets
	util.ListenServer(":1053")
}
