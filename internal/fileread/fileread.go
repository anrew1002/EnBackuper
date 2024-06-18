package fileread

import (
	"bufio"
	"entelekom/backuper/internal/models"
	"os"
)

func checkFile(filePath string) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		panic("Несуществующий файл сетей для сканирования")
	}
}

func readFile(filePath string) (*bufio.Scanner, *os.File, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, file, err
	}

	scanner := bufio.NewScanner(file)

	return scanner, file, err
	// content, _ := reader.ReadString('\n')

}

func ReadNetworks(filePath string) (networks []models.Network, err error) {
	checkFile(filePath)
	dataBuf, file, err := readFile(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	index := 0
	for dataBuf.Scan() {
		line := dataBuf.Text()
		net, err := models.NewNetwork(line)
		if err != nil {
			return nil, err
		}
		networks = append(networks, net)
		index++
	}

	// Check for errors during the scan
	if err = dataBuf.Err(); err != nil {
		return
	}
	return
}
