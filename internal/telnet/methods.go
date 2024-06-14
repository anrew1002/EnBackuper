package telnet

import (
	"strings"

	tlnt "github.com/reiver/go-telnet"
)

// Читает telnet буфер до тех пор пока не встретит знак приглашения на исполнение
func readTelnet(conn *tlnt.Conn) string {
	buff := ""
	for buff = ""; !strings.Contains(buff, "#"); {
		b := []byte{0}
		conn.Read(b)
		buff += string(b[0])
	}
	return buff
}
