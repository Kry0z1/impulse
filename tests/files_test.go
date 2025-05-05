package tests

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/Kry0z1/impulse/app"
	"github.com/Kry0z1/impulse/config"
	"github.com/Kry0z1/impulse/fs"
	"github.com/Kry0z1/impulse/lib"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dirs = []string{"not-finished", "basic"}

func TestRun(t *testing.T) {
	for _, dir := range dirs {
		t.Run(dir, func(t *testing.T) {
			t.Parallel()

			in := fs.MustLoadInputFile(fmt.Sprintf("files/%s/input", dir))
			cfg := config.MustLoad(fmt.Sprintf("files/%s/config.json", dir))

			out, err := os.Open(fmt.Sprintf("files/%s/output", dir))
			require.NoError(t, err, "Failed to open output file")

			log, err := os.Open(fmt.Sprintf("files/%s/log", dir))
			require.NoError(t, err, "Failed to open log file")

			var appOut, appLog bytes.Buffer

			orch := lib.NewOrchestrator(cfg)

			application := app.New(in, bufio.NewWriter(&appOut), bufio.NewWriter(&appLog), cfg, orch)
			err = application.Run()
			require.NoError(t, err, "Application failed but shouldn't")

			outReader := bufio.NewReader(out)
			logReader := bufio.NewReader(log)

			appOutReader := bufio.NewReader(&appOut)
			appLogReader := bufio.NewReader(&appLog)

			for {
				appLine, appErr := appOutReader.ReadString('\n')
				line, err := outReader.ReadString('\n')

				require.Equal(t, err, appErr)
				if err != nil {
					break
				}

				assert.Equal(t, line, appLine)
			}

			for {
				appLine, appErr := appLogReader.ReadString('\n')
				line, err := logReader.ReadString('\n')

				require.Equal(t, err, appErr)
				if err != nil {
					break
				}

				assert.Equal(t, line, appLine)
			}
		})
	}
}
