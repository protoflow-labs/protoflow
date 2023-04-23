package pkg

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"os/exec"
	"time"
)

func startProcess(cmd *exec.Cmd) (cleanup func(), err error) {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Error().Msgf("Error getting standard output pipe: %v", err)
		return
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error().Msgf("Error getting standard output pipe: %v", err)
		return
	}
	cleanup = func() {
		stdout.Close()
		stderr.Close()
	}

	stderrScan := bufio.NewScanner(stderr)
	go func() {
		for stderrScan.Scan() {
			fmt.Println(stderrScan.Text())
		}
	}()

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	log.Debug().Msgf("Starting subprocess: %v", cmd.String())

	cmd.WaitDelay = 10 * time.Minute

	// Start the subprocess
	if err = cmd.Start(); err != nil {
		log.Error().Msgf("Error starting subprocess: %v", err)
		return
	}
	return
}
