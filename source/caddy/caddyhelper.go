package caddy

import (
	"net/http"
	"os"
	"os/exec"
)

func StartCaddy() error {
	cmd := exec.Command("caddy", "run", "--config", "caddy.json")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func CheckCaddyAPI() (bool, error) {
	resp, err := http.Get("http://localhost:2025/config")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK, err
}
