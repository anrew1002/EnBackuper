package snmp

import (
	"fmt"
	"strings"

	snmp "github.com/gosnmp/gosnmp"
)

// Получение информации о коммутаторе по протоколу SNMPD
func GetSNMPDescription(ip string) (string, error) {

	snmp.Default.Target = ip
	snmp.Default.ExponentialTimeout = false
	// snmp.Default.Logger = snmp.NewLogger(log.New(os.Stdout, "", 0))

	err := snmp.Default.Connect()
	if err != nil {
		return "", fmt.Errorf("%s error getting desc: %w", ip, err)
	}
	defer snmp.Default.Conn.Close()

	oids := []string{".1.3.6.1.2.1.1.1.0"}
	result, err := snmp.Default.Get(oids)
	if err != nil {
		// log.Fatalf("Get() err: %v", err)
		return "", fmt.Errorf("%s error telnet: %w", snmp.Default.Target, err)
	}
	hostDescr := result.Variables[0].Value.([]byte)

	output := string(hostDescr)
	output = strings.Replace(output, " ", "_", -1)
	output = strings.Replace(output, "/", "_", -1)
	output = strings.Replace(output, "\n", "", -1)

	return output, nil
}
