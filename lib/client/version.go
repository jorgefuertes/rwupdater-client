package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"git.martianoids.com/queru/retroupdater-client/lib/build"
)

const dlURL = "https://updater.retrowiki.es/download/bin"

type Map map[string]interface{}

func checkErr(err error) error {
	if err != nil {
		fmt.Printf("FAIL (%s)\n", err.Error())
	}

	return err
}

func CheckUpdate() {
	// checking api for new version
	fmt.Printf("> Cheking for new version...%s...", build.VersionShort())
	resp, err := http.Get(API + "/version/client")
	if checkErr(err) != nil {
		return
	}
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	data := make(Map)
	err = d.Decode(&data)
	if checkErr(err) != nil {
		return
	}
	fmt.Printf("%s...", data["latest"])
	fmt.Printf("")
	if data["latest"] == build.VersionShort() {
		// up to date
		fmt.Println("OK")
		return
	}
	// need to update
	fmt.Println("UPDATE")

	// get my exe name
	exe, err := os.Executable()
	if checkErr(err) != nil {
		return
	}

	// updating
	fmt.Printf("> Updating myself to %s/%s %s\n", runtime.GOOS, runtime.GOARCH, data["latest"])

	// temp download to same path
	_, err = Download(
		data["latest"].(string),
		fmt.Sprintf("%s.tmp", exe),
		fmt.Sprintf("%s/%s/%s", dlURL, runtime.GOOS, runtime.GOARCH),
	)

	if checkErr(err) != nil {
		return
	}
	// x perm
	os.Chmod(exe+".tmp", os.FileMode(0755))
	// move over
	if err := os.Rename(exe+".tmp", exe); err != nil {
		fmt.Println("  Update ERROR:", err.Error())
		fmt.Println("  Please move by yourself", exe+".tmp", "to", exe)
		return
	}

	// re-execute
	args := []string{"-n"}
	args = append(args, os.Args[1:]...)
	cmd := exec.Command(exe, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
