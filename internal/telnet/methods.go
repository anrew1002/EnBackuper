package telnet

import (
	"fmt"
	"os"
	"strings"

	tlnt "github.com/reiver/go-telnet"
)

var (
	// Сделать в виде key-value storage miteshbsjat/textfilekv (?)
	backupCommandString = map[string]string{
		"DES-1210-28_ME_B2":                     "upload cfg_toTFTP %s %s.cfg",
		"DGS-3420-26SC_Gigabit_Ethernet_Switch": "upload cfg_toTFTP %s dest_file %s.cfg",
	}
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

func (tc *TelnetConnector) getBackupCommand(selfhost string, filename string) string {

	cmdStr, ok := backupCommandString[tc.device.Model]
	if !ok {
		return ""
	}
	command := fmt.Sprintf(cmdStr, selfhost, filename)
	return command
}

func writetoFile(d2 string) {
	f, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	// d2 := []byte{104, 101, 108, 108, 111, 32, 98, 121, 116, 101, 115}
	n2, err := f.WriteString(d2)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return

	}
	fmt.Println(n2, " bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
