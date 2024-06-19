package main

import (
	"context"
	"entelekom/backuper/internal/fileread"
	"entelekom/backuper/internal/sl"
	"entelekom/backuper/internal/workers"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

type Config struct {
	selfAddr string
	net_file string
	test     bool
}

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

	networks, err := fileread.ReadNetworks(config.net_file)
	if err != nil {
		return err
	}
	for i := range networks {
		for _, ips := range networks[i].GetIPs() {
			fmt.Println(ips)
		}
	}

	if config.test {
		fmt.Println("Test OK")
		return nil
	}

	backups := workers.ConcurrentBackup(log, networks, config.selfAddr)
	time.Sleep(30 * time.Second)
	counter := 0
	for _, filename := range backups {
		filepath := filepath.Join(`E:\EN+\tftpd\`, filename+".cfg")
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

func flagParse() Config {
	progname := filepath.Base(os.Args[0])
	selfAddr := flag.String("selfaddr", "", "Адрес сети с которого запускается backuper")
	filename := flag.String("file", "networks.txt", "Файл с сетями для сканирования")
	test := flag.Bool("test", false, "Запустить программу в холостую для проверки конфигурации")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			`%s запускает бэкапирование на коммутаторах:
 %s [Flags]

Flags:
`, progname, progname)
		flag.PrintDefaults()
	}
	flag.Parse()
	if *selfAddr == "" {
		log.Fatal("Selfaddr флаг должен быть указан")
	}
	return Config{selfAddr: *selfAddr, net_file: *filename, test: *test}
}
