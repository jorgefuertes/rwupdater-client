package client

import (
	"fmt"
	"os"
	"time"

	"git.martianoids.com/queru/retroupdater-client/lib/file"
	"github.com/cavaliercoder/grab"
	"github.com/dustin/go-humanize"
)

const API = "https://core.abadiaretro.com/api"

func Download(arch string, path string, f *file.File) uint64 {
	gc := grab.NewClient()
	req, _ := grab.NewRequest(
		path+"/"+f.Path+"/"+f.Name,
		API+"/files/download/"+arch+"/"+f.ID,
	)

	fmt.Printf("Downloading %s...", f.Name)
	res := gc.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("\rDownloading %s...%s...", f.Name,
				humanize.Bytes(uint64(res.BytesComplete())))
		case <-res.Done:
			t.Stop()
			if res.Err() != nil {
				fmt.Printf("FAIL (%s)\n", res.Err())
				fmt.Println("Deleting failed donwload")
				os.Remove(path + "/" + f.Path + "/" + f.Name)
			} else {
				fmt.Println("OK")
			}
			break Loop
		}
	}

	return uint64(res.BytesComplete())
}
