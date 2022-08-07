package helpers

import "strconv"

// JoinAddressAndPort 使用 : 拼接 address 和 port。
func JoinAddressAndPort(address string, port int) string {
	return address + ":" + strconv.Itoa(port)
}
