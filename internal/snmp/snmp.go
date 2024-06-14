package snmp

import (
	"fmt"
	"log"
	"os"

	snmp "github.com/gosnmp/gosnmp"
)

func GetSNMPDescription(ip string) (string, error) {

	snmp.Default.Target = ip
	snmp.Default.Logger = snmp.NewLogger(log.New(os.Stdout, "", 0))
	err := snmp.Default.Connect()
	if err != nil {
		return "", fmt.Errorf("%s error getting desc: %w", snmp.Default.Target, err)
	}
	defer snmp.Default.Conn.Close()

	oids := []string{".1.3.6.1.2.1.1.1.0"}
	result, err := snmp.Default.Get(oids)
	if err != nil {
		// log.Fatalf("Get() err: %v", err)
		return "", fmt.Errorf("%s error telnet: %w", snmp.Default.Target, err)
	}
	hostDescr := result.Variables[0].Value.([]byte)

	return string(hostDescr), nil
}
