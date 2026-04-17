package zpi_camera

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

type CameraProcess struct {
	cmd        *exec.Cmd
	camMutex   sync.Mutex
	feedActive bool
}

func NewCameraProcess() *CameraProcess {
	return &CameraProcess{}
}

func (p *CameraProcess) Start(width, height, fps uint32, output string) {
	p.feedActive = true

	go func() {
		for {
			p.camMutex.Lock()
			active := p.feedActive
			p.camMutex.Unlock()
			if !active {
				return
			}

			args := []string{
				"-t", "0", "--width", fmt.Sprintf("%d", width),
				"--height", fmt.Sprintf("%d", height),
				"--framerate", fmt.Sprintf("%d", fps),
				"--config", "~/ZPiCam.txt", "-o", output, "--listen", "--inline",
			}
			p.cmd = exec.Command("rpicam-vid", args...)

			fmt.Println("Camera starting (Waiting for connection...)")

			_ = p.cmd.Run()

			fmt.Println("Client disconnected. Restarting listener...")
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (p *CameraProcess) Stop() error {
	p.camMutex.Lock()
	p.feedActive = false
	p.camMutex.Unlock()

	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	if err := p.cmd.Process.Signal(os.Interrupt); err != nil {
		return p.cmd.Process.Kill()
	}

	return nil
}

func (p *CameraProcess) Kill() error {
	p.camMutex.Lock()
	p.feedActive = false
	p.camMutex.Unlock()

	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	return p.cmd.Process.Kill()
}
