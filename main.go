package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
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

func main() {
	path, err := getDiskPathWindows()
	if err != nil {
		panic(err)
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	fmt.Println(files)
}
