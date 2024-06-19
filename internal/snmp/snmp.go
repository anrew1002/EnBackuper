package snmp

import (
	"entelekom/backuper/internal/models"
	"errors"
	"fmt"
	"strings"
	"time"

	snmp "github.com/gosnmp/gosnmp"
)

var (
	ErrConnect = errors.New("failed to connect SNMP")
	ErrGetOID  = errors.New("failed getting oids")
)

// Получение информации о коммутаторе по протоколу SNMPD
func GetSNMPDescription(ip string) (models.Device, error) {

	//Создаем новое соединение SNMP
	ConnSNMP := &snmp.GoSNMP{
		Port:               161,
		Transport:          "udp",
		Community:          "public",
		Version:            snmp.Version2c,
		Timeout:            time.Duration(2) * time.Second,
		Retries:            3,
		ExponentialTimeout: false,
		MaxOids:            snmp.MaxOids,
		Target:             ip,
	}
	// Так как SNMP работает по протоколу UDP понять работает соединение
	// или нет возможно лишь после отправки первого запроса
	err := ConnSNMP.Connect()
	if err != nil {
		return models.Device{}, fmt.Errorf("%w, %s: %w", ErrConnect, ConnSNMP.Target, err)
	}
	defer ConnSNMP.Conn.Close()
	// .1.3.6.1.2.1.1.5.0 .1.3.6.1.2.1.1.1.0
	oids := []string{
		// модель
		".1.3.6.1.2.1.1.1.0",
		// имя
		".1.3.6.1.2.1.1.5.0",
	}
	result, err := ConnSNMP.Get(oids)
	if err != nil {
		return models.Device{}, fmt.Errorf("%w, %s: %w", ErrGetOID, ConnSNMP.Target, err)
	}
	// fmt.Printf("%+v", result.Variables)
	// fmt.Printf("%+v", string(result.Variables[0].Value.([]byte)))

	device := models.Device{
		Model: sanitizeString(string(result.Variables[0].Value.([]byte))),
		Name:  sanitizeString(string(result.Variables[1].Value.([]byte))),
	}

	return device, nil
}

func sanitizeString(hostDescr string) string {
	output := hostDescr
	output = strings.Replace(output, " ", "_", -1)
	output = strings.Replace(output, "/", "_", -1)
	output = strings.Replace(output, "\n", "", -1)
	return output
}
