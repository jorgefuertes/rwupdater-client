package client

import (
	"fmt"
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
			fmt.Println("OK")
			break Loop
		}
	}

	if err := res.Err(); err != nil {
		fmt.Println("FAIL")
	}

	return uint64(res.BytesComplete())
}
