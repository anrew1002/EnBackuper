package workers

import (
	"entelekom/backuper/internal/models"
	"entelekom/backuper/internal/sl"
	"entelekom/backuper/internal/snmp"
	"entelekom/backuper/internal/telnet"
	"errors"
	"log/slog"
	"sync"
)

func worker(log *slog.Logger, tftpAddr string, jobs <-chan string, results chan<- int) {
	for ip := range jobs {

		modelName, err := snmp.GetSNMPDescription(ip)
		if err != nil {
			switch {
			case errors.Is(err, snmp.ErrConnect):
				log.Debug("worker couldnt connect", sl.Err(err))
			case errors.Is(err, snmp.ErrGetOID):
				log.Error("worker error", sl.Err(err))
			}
			continue
		}

		tc, err := telnet.NewTelnetConnetor(modelName, ip+":23")
		if err != nil {
			log.Error("Failed obtain telnet connection", sl.Err(err))
			continue
		}
		err = tc.Backup(tftpAddr)
		if err != nil {
			log.Error("Failed backup", sl.Err(err))
			continue
		}
	}
}

func ConcurrentBackup(log *slog.Logger, networks []models.Network, tftpAddr string) {
	const numJobs = 5
	jobs := make(chan string, numJobs)
	results := make(chan int, numJobs)
	var wg sync.WaitGroup

	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func() {
			wg.Done()
			worker(log, tftpAddr, jobs, results)
		}()
	}
	for i := range networks {
		for _, ip := range networks[i].GetIPs() {
			jobs <- ip
		}
	}
	close(jobs)
	wg.Wait()

}
