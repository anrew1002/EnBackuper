package workers

import (
	"entelekom/backuper/internal/models"
	"entelekom/backuper/internal/sl"
	"entelekom/backuper/internal/snmp"
	"entelekom/backuper/internal/telnet"
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

func worker(log *slog.Logger, tftpAddr string, jobs <-chan string, results chan<- string) {
	fmt.Println("start worker")
	for ip := range jobs {
		fmt.Println(ip)
		modelName, err := snmp.GetSNMPDescription(ip)
		if err != nil {
			switch {
			case errors.Is(err, snmp.ErrConnect):
				log.Debug("worker couldnt connect", sl.Err(err))
			case errors.Is(err, snmp.ErrGetOID):
				log.Error("worker error", sl.Err(err))
			default:
				log.Error("Unexpected error", sl.Err(err))
			}
			continue
		}

		tc, err := telnet.NewTelnetConnector(modelName, ip+":23")
		if err != nil {
			log.Error("Failed obtain telnet connection", sl.Err(err))
			continue
		}
		backupFilename, err := tc.Backup(tftpAddr)
		if err != nil {
			log.Error("Failed backup", sl.Err(err))
			continue
		}
		results <- backupFilename
	}
}

// ConcurrentBackup Начинает мультипоточное выполнение бэкапирования
//
// Возращает названия файлов бэкапов, которые должны будут загрузить коммутаторы
func ConcurrentBackup(log *slog.Logger, networks []models.Network, tftpAddr string) []string {
	const numJobs = 3
	jobs := make(chan string, numJobs)
	results := make(chan string, numJobs)
	var wg sync.WaitGroup

	for w := 1; w <= numJobs; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(log, tftpAddr, jobs, results)
		}()
	}
	for i := range networks {
		for _, ip := range networks[i].GetIPs() {
			jobs <- ip
		}
	}
	backups := make([]string, 0, 40)
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for res := range results {
			backups = append(backups, res)
		}
	}()

	close(jobs)
	wg.Wait()
	close(results)
	wg2.Wait()
	return backups

}
