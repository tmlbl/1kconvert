package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const driveLetterScript = `
Get-PhysicalDisk | ForEach-Object {
    $physicalDisk = $_
    $physicalDisk |
        Get-Disk |
        Get-Partition |
        Where-Object DriveLetter |
        Select-Object DriveLetter, @{n='MediaType';e={ $physicalDisk.FriendlyName }}
}
`

func getDiskPathWindows() (string, error) {
	out, err := exec.Command("powershell", driveLetterScript).Output()
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(out), "\r\n")
	for _, l := range lines {
		if strings.Contains(l, "MPC1000") {
			return fmt.Sprintf("%s:\\", string(strings.TrimSpace(l)[0])), nil
		}
	}
	return "", fmt.Errorf("MPC drive not found")
}

func getFileList(path string) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if len(info.Name()) > 16 {
				shortened := info.Name()[len(info.Name())-16:]
				newPath := filepath.Join(filepath.Dir(path), shortened)
				fmt.Println("Renaming", path, "to", newPath)
				return os.Rename(path, newPath)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func main() {
	path, err := getDiskPathWindows()
	if err != nil {
		panic(err)
	}
	getFileList(path)
}
