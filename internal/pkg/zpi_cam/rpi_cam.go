package zpi_camera

import (
	"fmt"
	"os/exec"
)

type CameraProcess struct {
	cmd *exec.Cmd
}

func NewCameraProcess() *CameraProcess {
	return &CameraProcess{}
}

func (p *CameraProcess) Start(width, height, fps uint32, output string) error {
	args := []string{
		"-t", "0",
		"--width", fmt.Sprintf("%d", width),
		"--height", fmt.Sprintf("%d", height),
		"--framerate", fmt.Sprintf("%d", fps),
		"-o", output,
		"--listen",
	}

	p.cmd = exec.Command("rpicam-vid", args...)

	return p.cmd.Start()
}

func (p *CameraProcess) Kill() error {
	if p.cmd != nil && p.cmd.Process != nil {
		return p.cmd.Process.Kill()
	}
	return nil
}
