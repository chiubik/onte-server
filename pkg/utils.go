package pkg

import (
	"bufio"
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func CreateFile(filename string) (*os.File, error) {
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0777)
	}

	dst, err := os.Create(filepath.Join("uploads", filename))
	if err != nil {
		return nil, err
	}

	return dst, err
}

func GenerateIdentifier() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder

	for i := 0; i < 3; i++ {
		randomnIndex := rand.IntN(len(charset)) //randomn string from the constant
		sb.WriteByte(charset[randomnIndex])
	}

	return sb.String()
}

func CloudtunnelRun(ctx context.Context) (string, error) {
	re := regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)
	cmd := exec.CommandContext(ctx, "cloudflared", "tunnel", "--url", "http://localhost:3333")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Print(err)
		return "error: ", err
	}

	err = cmd.Start()
	if err != nil {
		fmt.Print(err)
		return "error: ", err
	}

	c := make(chan string, 1)

	go func() {
		scanner := bufio.NewScanner(stderr)
		found := false
		for scanner.Scan() {
			if found {
				continue
			}
			if url := re.FindString(scanner.Text()); url != "" {
				c <- url
				found = true
			}
		}
		cmd.Wait()
		close(c)
	}()

	select {
	case u, ok := <-c:
		if !ok {
			return "", fmt.Errorf("Cloudflared exited without printing a tunnel URL")
		}
		return u, nil
	case <-time.After(30 * time.Second):
		cmd.Process.Kill()
		return "", fmt.Errorf("Timed out waiting for new tunnel URL")
	}
}
