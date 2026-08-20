package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"industrial-key-rotation/internal/audit"
	"industrial-key-rotation/internal/persistence"
	"industrial-key-rotation/internal/report"
)

func main() {
	database := flag.String("db", filepath.Join(os.TempDir(), "industrial-key-rotation.db"), "bbolt database path")
	command := flag.String("command", "status", "status or report")
	sensorID := flag.String("sensor", "", "sensor identifier")
	formatName := flag.String("format", "text", "text, json, or csv")
	flag.Parse()
	store, err := persistence.Open(*database)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()
	ledger, err := audit.NewLedger(store)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch *command {
	case "status":
		counts, err := store.SnapshotCounts()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf("sensors=%d envelopes=%d rotations=%d audits=%d\n", counts["sensor"], counts["envelope"], counts["rotation"], counts["audit"])
	case "report":
		entries, err := ledger.Timeline(*sensorID)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		format, err := report.ParseFormat(*formatName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		output, err := report.Render(format, entries)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Print(string(output))
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", *command)
		os.Exit(2)
	}
}
