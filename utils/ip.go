package utils

import "net"

// IsIPv4 判断是否合法IPv4
func IsIPv4(s string) bool {
	ip := net.ParseIP(s)
	if ip == nil {
		return false
	}
	// ParseIP返回16字节数组，IPv4是后面4字节
	return ip.To4() != nil
}
