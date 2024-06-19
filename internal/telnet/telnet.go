package telnet

import (
	"entelekom/backuper/internal/models"
	"fmt"
	"time"

	tlnt "github.com/reiver/go-telnet"
)

type TelnetConnector struct {
	device        models.Device
	conn          *tlnt.Conn
	socket        string
	commandPrompt string
}

func NewTelnetConnetor(device models.Device, socket string) (*TelnetConnector, error) {
	conn, err := tlnt.DialTo(socket)
	if err != nil {
		return nil, err
	}
	// Временная стандартная строчка на приглашение
	commandPrompt := "#"
	return &TelnetConnector{device: device, socket: socket, conn: conn, commandPrompt: commandPrompt}, nil
}

func (tc *TelnetConnector) Backup(selfhost string) error {
	tc.Authenticate()

	// fmt.Print(readTelnet(tc.conn))
	// fmt.Println("")

	command := tc.getBackupCommand(selfhost)
	fmt.Println(command)
	tc.WriteRawCommand(command)
	time.Sleep(2 * time.Second)
	writetoFile(readTelnet(tc.conn))
	// fmt.Print(readTelnet(tc.conn))
	// if strings.Contains(readTelnet(tc.conn), "Success") {
	// 	return nil
	// }
	// // fmt.Println("")
	// return errors.New("Backup isnt succesful")
	return nil
}

func (tc *TelnetConnector) WriteRawCommand(command string) {
	tc.conn.Write([]byte(command + "\r\n"))
}

func (tc *TelnetConnector) Authenticate() {
	tc.WriteRawCommand("admin")
	tc.WriteRawCommand("admin")
}
