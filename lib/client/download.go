package client

import (
	"fmt"
	"os"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/dustin/go-humanize"
)

func Download(arch string, path string, file *File) uint64 {
	gc := grab.NewClient()
	req, _ := grab.NewRequest(
		path+"/"+file.Path+"/"+file.Name,
		API+"/files/download/"+arch+"/"+file.ID,
	)

	fmt.Printf("Downloading %s...", file.Name)
	res := gc.Do(req)

	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("\rDownloading %s...%s...", file.Name,
				humanize.Bytes(uint64(res.BytesComplete())))
		case <-res.Done:
			t.Stop()
			if res.Err() != nil {
				fmt.Printf("FAIL (%s)\n", res.Err())
				fmt.Println("Deleting failed donwload")
				os.Remove(path + "/" + file.Path + "/" + file.Name)
			} else {
				fmt.Println("OK")
			}
			break Loop
		}
	}

	return uint64(res.BytesComplete())
}
