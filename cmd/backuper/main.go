package main

import (
	"context"
	"entelekom/backuper/internal/config"
	"entelekom/backuper/internal/fileread"
	"entelekom/backuper/internal/sl"
	"entelekom/backuper/internal/workers"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	config := flagParse()
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	networks, err := fileread.ReadNetworks(config.NETWORKS_FILE)
	if err != nil {
		return err
	}
	for i := range networks {
		for _, ips := range networks[i].GetIPs() {
			fmt.Println(ips)
		}
	}

	if config.Test {
		fmt.Println("Test OK")
		return nil
	}

	backups := workers.ConcurrentBackup(log, networks, config.TFTP)
	time.Sleep(30 * time.Second)
	counter := 0
	for _, filename := range backups {
		filepath := filepath.Join(config.TFTPData, filename+".cfg")
		fileInfo, err := os.Stat(filepath)
		fmt.Println(filepath)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				log.Error("Бэкап отстутствует", "имя", filename)
				continue
			}
			log.Error("Ошибка проверки бэкапов", sl.Err(err))
			continue
		}
		if fileInfo.Size() == 0 {
			log.Error("Ошибка проверки бэкапов: нулевой размер", "имя", filename)
			continue
		}
		counter += 1
	}
	if counter == len(backups) {
		log.Info("Успешно завершенный процесс бэкапирования")
	} else {
		log.Error("Были бэкапированы не все коммутаторы")
	}
	// IPAdresses := []string{"192.168.47.55", "192.168.47.56"}

	// // TODO: разделить запуск и созадание
	// for i := range IPAdresses {

	// 	modelName, err := snmp.GetSNMPDescription(IPAdresses[i])
	// 	if err != nil {
	// 		log.Error("error getting model name", sl.Err(err))
	// 		return err
	// 	}

	// 	fmt.Println(modelName)
	// 	fmt.Println(i)
	// 	tc, err := telnet.NewTelnetConnetor(modelName, IPAdresses[i]+":23")
	// 	if err != nil {
	// 		log.Error("Failed obtain telnet connection")
	// 		continue
	// 	}
	// 	err = tc.Backup(config.selfAddr)
	// 	if err != nil {
	// 		log.Error("Failed backup", sl.Err(err))
	// 	}
	// }

	// go func() {
	// 	log.Info(fmt.Sprintf("listening...on %s", httpServer.Addr))
	// 	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		fmt.Fprintf(os.Stderr, "error listening and serving: %s.\n", err)
	// 	}
	// }()

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	<-ctx.Done()
	// 	if err := httpServer.Shutdown(ctx); err != nil {
	// 		fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
	// 	}
	// }()
	// wg.Wait()
	return nil
}

func flagParse() config.Config {

	progname := filepath.Base(os.Args[0])
	filename := flag.String("file", "", "Файл с сетями для сканирования")
	test := flag.Bool("test", false, "Запустить программу в холостую для проверки конфигурации")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`%s запускает бэкапирование на коммутаторах:
 %s [Flags]

 Flags:
 `,
			progname, progname)
		flag.PrintDefaults()
	}

	config := config.MustLoad()
	flag.Parse()
	config.Test = *test
	if *filename != "" {
		config.NETWORKS_FILE = *filename
	}
	return *config
}
