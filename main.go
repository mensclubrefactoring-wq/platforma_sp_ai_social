package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func startService(name, path string) *exec.Cmd {
	cmd := exec.Command("go", "run", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start service %s: %v", name, err)
	}
	fmt.Printf("Service %s started as PID %d\n", name, cmd.Process.Pid)
	return cmd
}

func main() {
	// Start microservices
	authCmd := startService("Auth", "./services/auth/main.go")
	tasksCmd := startService("Tasks", "./services/tasks/main.go")
	aiCmd := startService("AI", "./services/ai/main.go")
	gatewayCmd := startService("Gateway", "./services/gateway/main.go")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nShutting down services...")
	authCmd.Process.Kill()
	tasksCmd.Process.Kill()
	aiCmd.Process.Kill()
	gatewayCmd.Process.Kill()
}
