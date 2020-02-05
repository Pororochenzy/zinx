package utils

import (
	"fmt"
	"testing"

)

func TestServer(t *testing.T) {
	GlobalObject.Reload()

	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		GlobalObject.Version,
		GlobalObject.MaxConn,
		GlobalObject.MaxPacketSize)
}
