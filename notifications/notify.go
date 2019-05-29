package notifications

import (
	"bytes"
	"log"
	"os/exec"
)

func Message(message string) {
	cmd := exec.Command("notify-send", "Deluge", message, "--icon=delauncher")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	if err := cmd.Run(); err != nil {
		log.Fatalf("unable to send notification: %s", err.Error())
	}
}
