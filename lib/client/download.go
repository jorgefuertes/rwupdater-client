package client

import (
	"fmt"
	"os"
	"time"

	"git.martianoids.com/queru/retroupdater-client/lib/file"
	"github.com/cavaliercoder/grab"
	"github.com/dustin/go-humanize"
)

const API = "https://updater.retrowiki.es/api"

func Download(name, dst, url string) (uint64, error) {
	gc := grab.NewClient()
	req, _ := grab.NewRequest(
		dst,
		url,
	)

	fmt.Printf("  - Downloading %s...", name)
	res := gc.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("\r  - Downloading %s...%s...", name,
				humanize.Bytes(uint64(res.BytesComplete())))
		case <-res.Done:
			t.Stop()
			if res.Err() != nil {
				fmt.Printf("FAIL (%s)\n", res.Err())
				fmt.Println("  - Deleting failed donwload")
				os.Remove(dst)
				return 0, res.Err()
			} else {
				fmt.Println("OK")
			}
			break Loop
		}
	}

	return uint64(res.BytesComplete()), nil
}

func UpdateFile(arch string, path string, f *file.File) uint64 {
	name := f.Path + "/" + f.Name
	if len(name) > 40 {
		name = "…" + name[len(name)-40:]
	}
	bytes, _ := Download(name, path+"/"+f.Path+"/"+f.Name, API+"/files/download/"+arch+"/"+f.ID)
	return bytes
}
