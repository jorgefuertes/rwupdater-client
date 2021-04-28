package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"git.martianoids.com/queru/retroupdater-client/lib/build"
	"git.martianoids.com/queru/retroupdater-client/lib/client"
	"github.com/alecthomas/kong"
	"github.com/dustin/go-humanize"
)

// CLI
var CLI struct {
	Version struct {
	} `cmd:"" help:"Display version and exit."`
	Update struct {
		Arch string `arg:"" enum:"mister" help:"FPGA Achitecture: mister."`
		Root string `arg:"" name:"path" type:"existingDir" help:"SD card root path."`
	} `cmd:"" help:"Update architecture over given path."`
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
	ctx := kong.Parse(&CLI)
	if ctx.Command() == "version" {
		fmt.Print(build.Version())
		os.Exit(0)
	}

	root := CLI.Update.Root
	arch := CLI.Update.Arch

	if len(root) > 1 {
		root = strings.TrimSuffix(root, "/")
	}

	// check if sdcard root exists
	_, err := os.Stat(root)
	check(err)

	fmt.Printf("Getting catalog...")
	files, err := client.GetCatalog(arch)
	check(err)
	fmt.Println("OK")

	for _, file := range *files {
		stats.Total += 1

		// check and create dir
		_, err := os.Stat(root + "/" + file.Path)
		if err != nil {
			fmt.Printf("Creating dir %s...", file.Path)
			if err := os.MkdirAll(root+"/"+file.Path, 0755); err != nil {
				fmt.Println("FAIL")
				check(err)
			}
			fmt.Println("OK")
			stats.Dir += 1
		}

		// check file
		// download if doesn't exists
		_, err = os.Stat(root + "/" + file.Path + "/" + file.Name)
		if err != nil {
			stats.TotalBytes += client.Download(arch, root, &file)
			stats.Download += 1
			// delete old version if core
			if len(file.Version) > 0 {
				dir, err := ioutil.ReadDir(root + "/" + file.Path)
				if err != nil {
					fmt.Println("ERROR:", err)
				} else {
					for _, entry := range dir {
						if entry.Name() == file.Name {
							continue
						}
						if strings.HasPrefix(entry.Name(), file.Core) {
							fmt.Printf("Deleting %s...", entry.Name())
							if err := os.Remove(root + "/" + file.Path + "/" + entry.Name()); err != nil {
								fmt.Printf("FAIL (%s)", err.Error())
							} else {
								fmt.Println("OK")
							}
						}
					}
				}
			}
			continue
		}
	}

	// stats
	fmt.Printf("Completed in %.2f seconds\n", time.Since(stats.Begin).Seconds())
	fmt.Printf("%v files checked, %v were up to date\n", stats.Total, stats.UpToDate)
	fmt.Printf("%v dirs created, %v files download, %s\n",
		stats.Dir, stats.Download, humanize.Bytes(stats.TotalBytes))
}
