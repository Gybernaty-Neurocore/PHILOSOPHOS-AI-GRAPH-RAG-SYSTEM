package services

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func GetEmbedding(text string) ([]float32, error) {
	fmt.Printf("Calling GetEmbedding for text: %q\n", text)

	if _, err := os.Stat("embed.py"); os.IsNotExist(err) {
		return nil, fmt.Errorf("embed.py not found in current directory")
	}

	cmd := exec.Command("python", "embed.py", text)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start Python script: %v", err)
	}

	stdoutBytes, _ := io.ReadAll(stdout)
	stderrBytes, _ := io.ReadAll(stderr)

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Python script failed: %v, stderr: %s\n", err, string(stderrBytes))
		return nil, fmt.Errorf("failed to run Python script: %v, stderr: %s", err, string(stderrBytes))
	}

	var result []float32
	if err := json.Unmarshal(stdoutBytes, &result); err != nil {
		fmt.Printf("Error parsing Python output: %v, stdout: %s\n", err, string(stdoutBytes))
		return nil, fmt.Errorf("failed to parse Python output: %v", err)
	}

	fmt.Printf("Embedding retrieved successfully, length: %d\n", len(result))
	return result, nil
}
