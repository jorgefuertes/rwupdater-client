package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"git.martianoids.com/queru/retroupdater-client/lib/client"
	"github.com/alecthomas/kong"
	"github.com/dustin/go-humanize"
)

// Config
type Config struct {
	Arch string `arg:"" enum:"mist,mister" help:"FPGA Achitecture: mist or mister."`
	Root string `arg:"" type:"existingDir" help:"SD card root path."`
}

// Stats
type Stats struct {
	Download   int
	UpToDate   int
	Dir        int
	Total      int
	TotalBytes uint64
	Begin      time.Time
}

// check for error
func check(err error) {
	if err == nil {
		return
	}

	fmt.Println("ERROR:", err)
	os.Exit(1)
}

func main() {
	var stats Stats
	stats.Begin = time.Now()

	// command line
	cfg := new(Config)
	kong.Parse(cfg)

	if len(cfg.Root) > 1 {
		cfg.Root = strings.TrimSuffix(cfg.Root, "/")
	}

	// check if sdcard root exists
	_, err := os.Stat(cfg.Root)
	check(err)

	fmt.Printf("Getting catalog...")
	files, err := client.GetCatalog(cfg.Arch)
	check(err)
	fmt.Println("OK")

	for _, file := range *files {
		stats.Total += 1

		// check and create dir
		_, err := os.Stat(cfg.Root + "/" + file.Path)
		if err != nil {
			fmt.Printf("Creating dir %s...", file.Path)
			if err := os.MkdirAll(cfg.Root+"/"+file.Path, 0755); err != nil {
				fmt.Println("FAIL")
				check(err)
			}
			fmt.Println("OK")
			stats.Dir += 1
		}

		// check file, download if doesn't exists
		_, err = os.Stat(cfg.Root + "/" + file.Path + "/" + file.Name)
		if err != nil {
			stats.TotalBytes += client.Download(cfg.Arch, cfg.Root, &file)
			stats.Download += 1
		}
	}

	// stats
	fmt.Printf("Completed in %.2f seconds\n", time.Since(stats.Begin).Seconds())
	fmt.Printf("%v files checked, %v were up to date\n", stats.Total, stats.UpToDate)
	fmt.Printf("%v dirs created, %v files download, %s\n",
		stats.Dir, stats.Download, humanize.Bytes(stats.TotalBytes))
}
