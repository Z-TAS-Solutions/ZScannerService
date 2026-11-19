package zpi_camera

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
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

// Calls ZStreamer with given width, height and fps parameters
// ZStreamer starts waiting on ports 8888 and 8889 and supports only max 2 clients
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
				fmt.Sprintf("%d", width),
				fmt.Sprintf("%d", height),
				fmt.Sprintf("%d", fps),
			}

			p.cmd = exec.Command("/home/zischl/ZTAS/ZStreamer", args...)
			p.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
			p.cmd.Stderr = os.Stderr
			p.cmd.Stdout = os.Stdout

			fmt.Println("[Camera] ZTAS Streamer starting on ports 8888 & 8889...")

			_ = p.cmd.Run()

			fmt.Println("[Camera] ZTAS Streamer exited. Restarting...")
			time.Sleep(200 * time.Millisecond)
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

	_ = syscall.Kill(-p.cmd.Process.Pid, syscall.SIGTERM)
	time.Sleep(100 * time.Millisecond)
	_ = syscall.Kill(-p.cmd.Process.Pid, syscall.SIGKILL)
	return nil
}

func (p *CameraProcess) Kill() error {
	return p.Stop()
}
