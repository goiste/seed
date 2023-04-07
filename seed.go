package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	seed "github.com/goiste/seed/lib"

	"github.com/spf13/pflag"
)

var (
	dsnString = pflag.StringP("dsn", "d", "", "DSN string")
	seedFile  = pflag.StringP("file", "f", "./seed.json", "Data for seeding. Accepted filetypes: .sql, .json, .yaml, .yml")
	seedMode  = pflag.StringP("mode", "m", "file", "Seed mode. Accepted values: file, sql")
	randSeed  = pflag.Int64P("seed", "s", 0, "Seed value for pseudo-random source")
	batchSize = pflag.IntP("batch-size", "b", 0, "Batch size (for DB)")
)

func init() {
	pflag.Parse()
}

func main() {
	start := time.Now()

	if *dsnString == "" {
		log.Fatal("DSN string is required")
	}

	if *randSeed == 0 {
		*randSeed = time.Now().UnixNano()
	}

	mode, ok := seed.TypeFromString(*seedMode)
	if !ok {
		mode = ""
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	seeder, err := seed.NewSeeder(ctx, seed.Config{
		File:      *seedFile,
		DSN:       *dsnString,
		SeedType:  mode,
		RandSeed:  *randSeed,
		BatchSize: *batchSize,
	})
	if err != nil {
		log.Fatalf("create new seeder error: %v", err)
	}

	done := make(chan struct{})
	defer close(done)

	errors := make(chan error)
	defer close(errors)

	go func() {
		err = seeder.Seed(ctx)
		if err != nil {
			errors <- fmt.Errorf("seeding error: %w", err)
			return
		}

		done <- struct{}{}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signals:
		log.Println("interrupted, no data was written to the database")
	case e := <-errors:
		log.Fatal(e)
	case <-done:
		log.Printf("success, execution time: %v\n", time.Since(start))
	}
}
