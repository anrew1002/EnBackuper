package telnet

import (
	"entelekom/backuper/internal/models"
	"fmt"
	"strconv"
	"time"

	tlnt "github.com/reiver/go-telnet"
)

type TelnetConnector struct {
	device        models.Device
	conn          *tlnt.Conn
	socket        string
	commandPrompt string
}

func NewTelnetConnector(device models.Device, socket string) (*TelnetConnector, error) {
	conn, err := tlnt.DialTo(socket)
	if err != nil {
		return nil, err
	}
	// Временная стандартная строчка на приглашение
	commandPrompt := "#"
	return &TelnetConnector{device: device, socket: socket, conn: conn, commandPrompt: commandPrompt}, nil
}

// Backup с помощью telnet отправляет команды на бэкапирование коммутатору
// при этом не гарантируется что бэкап будет создан
// Возращает имя файла бэкапа и ошибку.
func (tc *TelnetConnector) Backup(tftp string) (string, error) {
	tc.Authenticate()

	// fmt.Print(readTelnet(tc.conn))
	// fmt.Println("")
	filename := fmt.Sprintf("%s_%s", strconv.FormatInt(time.Now().UTC().UnixNano(), 10), tc.device.Name)
	command := tc.getBackupCommand(tftp, filename)
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
	return filename, nil
}

func (tc *TelnetConnector) WriteRawCommand(command string) {
	tc.conn.Write([]byte(command + "\r\n"))
}

func (tc *TelnetConnector) Authenticate() {
	tc.WriteRawCommand("admin")
	tc.WriteRawCommand("admin")
}
