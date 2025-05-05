package app

import (
	"bufio"
	"fmt"
	"github.com/Kry0z1/impulse/config"
	"github.com/Kry0z1/impulse/lib"
	"io"
)

type App struct {
	input  *bufio.Reader
	output *bufio.Writer
	log    *bufio.Writer
	cfg    *config.Config
	orch   *lib.Orchestrator
}

func (a *App) ParseLine() error {
	line, err := a.input.ReadString('\n')
	if err == io.EOF {
		return io.EOF
	}

	output, err := a.orch.ParseLine(line[:len(line)-1])
	if err != nil {
		return err
	}
	_, err = a.log.WriteString(output + "\n")

	if err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	_ = a.log.Flush()

	return nil
}

func (a *App) Run() error {
	for {
		err := a.ParseLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	_, err := a.output.WriteString(a.orch.Result())
	if err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}
	_ = a.output.Flush()

	return nil
}

func New(input *bufio.Reader, output *bufio.Writer, log *bufio.Writer, cfg *config.Config, orch *lib.Orchestrator) App {
	return App{
		input:  input,
		output: output,
		log:    log,
		cfg:    cfg,
		orch:   orch,
	}
}
