package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.bug.st/serial"
)

// commandError is an error type to wrap os/exec.Command errors. This provides
// some more information regarding what went wrong while running a command.
type commandError struct {
	Msg  string
	File string
	Err  error
}

func (e *commandError) Error() string {
	return e.Msg + " " + e.File + ": " + e.Err.Error()
}

func flash(port, target, msdFirmwareName string) error {
	flashVolume := getFlashVolumeFromBuildTag(target)
	return _flash(port, flashVolume, msdFirmwareName)
}

func getFlashVolumeFromBuildTag(target string) string {
	ret := ``
	switch target {
	case `pyportal`:
		ret = `PORTALBOOT`
	case `feather-m4`:
		ret = `FEATHERBOOT`
	case `trinket-m0`:
		ret = `TRINKETBOOT`
	}
	return ret
}

func _flash(port, flashVolume, tmppath string) error {

	err := touchSerialPortAt1200bps(port)
	if err != nil {
		return err
	}
	time.Sleep(3 * time.Second)

	err = flashUF2UsingMSD(flashVolume, tmppath)
	if err != nil {
		return &commandError{"failed to flash", tmppath, err}
	}
	return nil
}

// copyFile copies the given file from src to dst. It copies first to a .tmp
// file which is then moved over a possibly already existing file at the
// destination.
func copyFile(src, dst string) error {
	inf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer inf.Close()
	outpath := dst + ".tmp"
	outf, err := os.Create(outpath)
	if err != nil {
		return err
	}

	_, err = io.Copy(outf, inf)
	if err != nil {
		os.Remove(outpath)
		return err
	}

	err = outf.Close()
	if err != nil {
		return err
	}

	return os.Rename(dst+".tmp", dst)
}

func touchSerialPortAt1200bps(port string) (err error) {
	retryCount := 3
	for i := 0; i < retryCount; i++ {
		// Open port
		p, e := serial.Open(port, &serial.Mode{BaudRate: 1200})
		if e != nil {
			//fmt.Println("DEBUG:", i, e.Error()) // <<-- debug print
			//for j := 0; j < retryCount; j++ {
			//	port, err := getDefaultPort()
			//	if err != nil {
			//		fmt.Println("  touch > get :", i, j, err.Error())
			//	} else {
			//		fmt.Println("  touch > get :", i, j, port)
			//		break
			//	}
			//	time.Sleep(1 * time.Second)
			//}
			//time.Sleep(1 * time.Second)
			//err = e
			//continue
			se, ok := e.(*serial.PortError)
			if ok && se.Code() == serial.InvalidSerialPort {
				// とりあえず OK とする
				return nil
			}

			time.Sleep(1 * time.Second)
			err = e
			fmt.Printf("err: %d : %s\n", i, err.Error())
			continue
		}
		defer p.Close()

		p.SetDTR(false)
		return nil
	}
	return fmt.Errorf("opening port: %s", err)
}

func flashUF2UsingMSD(volume, tmppath string) error {
	// find standard UF2 info path
	var infoPath string
	switch runtime.GOOS {
	case "linux", "freebsd":
		infoPath = "/media/*/" + volume + "/INFO_UF2.TXT"
	case "darwin":
		infoPath = "/Volumes/" + volume + "/INFO_UF2.TXT"
	case "windows":
		path, err := windowsFindUSBDrive(volume)
		if err != nil {
			return err
		}
		infoPath = path + "/INFO_UF2.TXT"
	}

	d, err := filepath.Glob(infoPath)
	if err != nil {
		return err
	}
	if d == nil {
		return errors.New("unable to locate UF2 device: " + volume)
	}

	return copyFile(tmppath, filepath.Dir(d[0])+"/flash.uf2")
}

func windowsFindUSBDrive(volume string) (string, error) {
	cmd := exec.Command("wmic",
		"PATH", "Win32_LogicalDisk", "WHERE", "VolumeName = '"+volume+"'",
		"get", "DeviceID,VolumeName,FileSystem,DriveType")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(out.String(), "\n") {
		words := strings.Fields(line)
		if len(words) >= 3 {
			if words[1] == "2" && words[2] == "FAT" {
				return words[0], nil
			}
		}
	}
	return "", errors.New("unable to locate a USB device to be flashed")
}

// getDefaultPort returns the default serial port depending on the operating system.
func getDefaultPort() (port string, err error) {
	var portPath string
	switch runtime.GOOS {
	case "darwin":
		portPath = "/dev/cu.usb*"
	case "linux":
		portPath = "/dev/ttyACM*"
	case "freebsd":
		portPath = "/dev/cuaU*"
	case "windows":
		cmd := exec.Command("wmic",
			"PATH", "Win32_SerialPort", "WHERE", "Caption LIKE 'USB % (COM%)'", "GET", "DeviceID")

		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return "", err
		}

		if out.String() == "No Instance(s) Available." {
			return "", errors.New("no serial ports available")
		}

		for _, line := range strings.Split(out.String(), "\n") {
			words := strings.Fields(line)
			if len(words) == 1 {
				if strings.Contains(words[0], "COM") {
					return words[0], nil
				}
			}
		}
		return "", errors.New("unable to locate a serial port")
	default:
		return "", errors.New("unable to search for a default USB device to be flashed on this OS")
	}

	d, err := filepath.Glob(portPath)
	if err != nil {
		return "", err
	}
	if d == nil {
		return "", errors.New("unable to locate a serial port")
	}

	return d[0], nil
}
