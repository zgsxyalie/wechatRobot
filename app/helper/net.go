package helper

import (
	"fmt"
	"net"
	"time"
)

// IsPortOpen 检查端口是否打开
func IsPortOpen(port uint, timeout time.Duration) bool {
	target := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
