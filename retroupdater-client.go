package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"git.martianoids.com/queru/retroupdater-client/lib/client"
	"github.com/alecthomas/kong"
)

// Config
type Config struct {
	Arch string `arg:"" enum:"mist,mister" help:"FPGA Achitecture: mist or mister."`
	Root string `arg:"" type:"existingDir" help:"SD card root path."`
}

// Stats
type Stats struct {
	Updated    int
	Created    int
	UpToDate   int
	DirCreated int
	Total      int
	Begin      time.Time
	End        time.Time
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

	files, err := client.GetCatalog(cfg.Arch)
	check(err)

	for _, file := range *files {
		stats.Total += 1
		_, err := os.Stat(cfg.Root + "/" + file.Path)
		if err != nil {
			fmt.Printf("Creating dir %s...", file.Path)
			if err := os.MkdirAll(cfg.Root+"/"+file.Path, 0755); err != nil {
				fmt.Println("FAIL")
				check(err)
			}
			fmt.Println("OK")
			stats.DirCreated += 1
		}
	}
}
