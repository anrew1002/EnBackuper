package telnet

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tlnt "github.com/reiver/go-telnet"
)

type TelnetConnector struct {
	modelName     string
	conn          *tlnt.Conn
	socket        string
	commandPrompt string
}

func NewTelnetConnetor(modelName string, socket string) (*TelnetConnector, error) {
	conn, err := tlnt.DialTo(socket)
	if err != nil {
		return nil, err
	}
	// Временная стандартная строчка на приглашение
	commandPrompt := "#"
	return &TelnetConnector{modelName: modelName, socket: socket, conn: conn, commandPrompt: commandPrompt}, nil
}

func (tc *TelnetConnector) Backup(selfhost string) error {
	tc.Authenticate()

	fmt.Print(readTelnet(tc.conn))
	fmt.Println("")

	command := tc.getCommand(selfhost)
	tc.WriteRawCommand(command)
	// tc.conn.Write([]byte("q"))
	fmt.Print(readTelnet(tc.conn))
	fmt.Println("")
	return nil
}

func (tc *TelnetConnector) getCommand(selfhost string) string {
	filename := fmt.Sprintf("%s_%s", strconv.FormatInt(time.Now().UTC().UnixNano(), 10), tc.modelName)
	filename = strings.Replace(filename, " ", "_", -1)
	filename = strings.Replace(filename, "\n", "", -1)
	command := fmt.Sprintf("upload cfg_toTFTP %s dest_file %s.cfg", selfhost, filename)
	return command
}

func (tc *TelnetConnector) WriteRawCommand(command string) {
	tc.conn.Write([]byte(command + "\r\n"))
}
func (tc *TelnetConnector) Authenticate() {
	tc.WriteRawCommand("admin")
	tc.WriteRawCommand("admin")
}
