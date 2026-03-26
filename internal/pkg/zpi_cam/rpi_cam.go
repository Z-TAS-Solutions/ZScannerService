package zpi_camera

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type CameraProcess struct {
	cmd        *exec.Cmd
	camMutex   sync.Mutex
	feedActive bool
}

func NewCameraProcess() *CameraProcess {
	return &CameraProcess{}
}

func (p *CameraProcess) Start(width, height, fps uint32, output string) error {
	p.camMutex.Lock()
	defer p.camMutex.Unlock()

	if p.feedActive {
		return errors.New("camera already running")
	}

	args := []string{
		"-t", "0",
		"--width", fmt.Sprintf("%d", width),
		"--height", fmt.Sprintf("%d", height),
		"--framerate", fmt.Sprintf("%d", fps),
		"-o", output,
		"--listen",
	}

	cmd := exec.Command("rpicam-vid", args...)

	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	p.cmd = cmd
	p.feedActive = true

	return nil
}

func (p *CameraProcess) Stop() error {
	p.camMutex.Lock()
	defer p.camMutex.Unlock()

	if !p.feedActive || p.cmd == nil || p.cmd.Process == nil {
		return errors.New("camera not running")
	}

	if err := p.cmd.Process.Signal(os.Interrupt); err != nil {
		return p.cmd.Process.Kill()
	}

	return nil
}

func (p *CameraProcess) Kill() error {
	p.camMutex.Lock()
	defer p.camMutex.Unlock()

	if !p.feedActive || p.cmd == nil || p.cmd.Process == nil {
		return errors.New("camera not running")
	}

	return p.cmd.Process.Kill()
}

func (p *CameraProcess) IsRunning() bool {
	p.camMutex.Lock()
	defer p.camMutex.Unlock()
	return p.feedActive
}
