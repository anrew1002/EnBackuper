package main

import (
	"context"
	"entelekom/backuper/internal/sl"
	"entelekom/backuper/internal/snmp"
	"entelekom/backuper/internal/telnet"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
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

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	IPAdresses := []string{"192.168.47.55", "192.168.47.56"}

	modelName, err := snmp.GetSNMPDescription(IPAdresses[0])
	if err != nil {
		log.Error("error getting model name", sl.Err(err))
		return err
	}

	// TODO: разделить запуска и созадание
	tc, err := telnet.NewTelnetConnetor(modelName, IPAdresses[0]+":23")
	if err != nil {
		return err
	}
	tc.Backup("10.10.1.2")
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
