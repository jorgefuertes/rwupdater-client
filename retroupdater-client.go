package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"git.martianoids.com/queru/retroupdater-client/lib/build"
	"git.martianoids.com/queru/retroupdater-client/lib/client"
	"git.martianoids.com/queru/retroupdater-client/lib/file"
	"github.com/alecthomas/kong"
	"github.com/dustin/go-humanize"
)

// CLI
var CLI struct {
	Version  struct{} `cmd:"" help:"Display version and exit."`
	ArchList struct{} `cmd:"" help:"List of supported achitectures."`
	Update   struct {
		Arch string `arg:"" help:"FPGA Achitecture: use arch-list for supported architectures"`
		Root string `arg:"" name:"path" type:"existingDir" help:"SD card root path."`
	} `cmd:"" help:"Update architecture over given path. Try update --help."`
}

// Stats
type Stats struct {
	Download   uint
	UpToDate   uint
	Dir        uint
	Total      uint
	TotalBytes uint64
	Deleted    uint
	Begin      time.Time
}

// check for error
func check(err error) {
	if err == nil {
		return
	}

	panic(err)
	// fmt.Printf("\n> *** ERROR %s ***\n", err)
	// os.Exit(1)
}

// ok
func ok() {
	fmt.Println("OK")
}

// fail
func fail(err error) {
	fmt.Printf("FAIL (%s)\n", err.Error())
}

func main() {
	var stats Stats
	stats.Begin = time.Now()

	// command line
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "version":
		fmt.Print(build.Version())
		os.Exit(0)
	case "arch-list":
		fmt.Print("> Getting architecture list...")
		list, err := client.GetArchList()
		check(err)
		ok()
		for _, arch := range list {
			fmt.Printf("  - %s\n", arch)
		}
		os.Exit(0)
	}

	root := CLI.Update.Root
	arch := CLI.Update.Arch

	if !client.IsArch(arch) {
		fmt.Println("Unknown architecture. Please try 'arch-list'.")
		os.Exit(1)
	}

	if root == "/" || root == "\\" {
		fmt.Println("Refusing to work in your root directory")
		os.Exit(1)
	}

	if len(root) > 1 {
		root = strings.TrimSuffix(root, "/")
		root = strings.TrimSuffix(root, "\\")
	}

	// check if sdcard root exists
	_, err := os.Stat(root)
	check(err)

	// banner and version
	fmt.Print(build.Version())

	// catalog
	fmt.Printf("> Getting remote catalog...")
	rcat, err := client.GetCatalog(arch)
	check(err)
	ok()
	fmt.Printf("> Getting local catalog...")
	lcat, err := file.GetLocalCatalog(root, "", 0)
	check(err)
	ok()

	for _, r := range *rcat {
		stats.Total += 1

		// check and create dir
		if !lcat.PathExists(r.Path) {
			_, err := os.Stat(r.CompletePath(root))
			if err == os.ErrExist {
				fmt.Printf("> Creating dir %s...", r.Path)
				if err := os.MkdirAll(r.CompletePath(root), 0755); err != nil {
					fail(err)
					check(err)
				}
				ok()
				stats.Dir += 1
			}
		}

		// download if doesn't exists
		_, err := lcat.Find(r.ID)
		if err != nil {
			stats.TotalBytes += client.Download(arch, root, &r)
			stats.Download += 1
			// delete old versions if its a core
			for {
				l, err := lcat.FindByName(r.Path, r.Name)
				if err != nil {
					break
				}
				fmt.Printf("Deleting %s...", l.Abbr())
				if err := os.Remove(l.CompleteFileName(root)); err != nil {
					fail(err)
					continue
				}
				ok()
			}
		}
	}

	// stats
	fmt.Printf("> Completed in %.2f seconds\n", time.Since(stats.Begin).Seconds())
	fmt.Printf("  %v files checked, %v were up to date\n", stats.Total, stats.UpToDate)
	fmt.Printf("  %v dirs created, %v files download, %s\n",
		stats.Dir, stats.Download, humanize.Bytes(stats.TotalBytes))
}
