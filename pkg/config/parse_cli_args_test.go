package config

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"

	"github.com/markhork8s/markhor/pkg"
	"github.com/spf13/viper"
)

func TestExitCodeForHelpFlags(t *testing.T) {
	t.Parallel()
	flags := []string{"-h", "--help"}
	for _, flag := range flags {
		f := flag
		t.Run("flag "+f, func(t *testing.T) {
			t.Parallel()
			cmd := exec.Command("go", "run", "../../main.go", f)
			outputb, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Expected program to exit normally for flag %s, but it didn't", f)
			}

			var exitCode int
			done := make(chan struct{})
			go func() {
				defer func() {
					if r := recover(); r != nil {
						exitCode = 0
					}
					close(done)
				}()
				os.Exit(0)
			}()

			<-done

			if exitCode != 0 {
				t.Fatalf("Expected os.Exit(0) to be called for flag %s, but it wasn't", f)
			}
			output := string(outputb)
			if !strings.Contains(output, versionHelpMessage) {
				t.Fatal("The help message for the version is absent from the --help output: ", output)
			}
			if !strings.Contains(output, pkg.DEFAULT_CONFIG_PATH) {
				t.Fatal("The default config path is absent from the --help output: ", output)
			}
		})
	}
}

func TestExitCodeForVersionFlags(t *testing.T) {
	t.Parallel()
	flags := []string{"-v", "--version"}

	for _, flag := range flags {
		f := flag
		t.Run("flag "+f, func(t *testing.T) {
			t.Parallel()
			cmd := exec.Command("go", "run", "../../main.go", f)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := cmd.Run()
				if err != nil {
					t.Error("Expected program to exit normally, but it didn't:", err)
				}
			}()

			wg.Wait()

			output := stdout.String()
			if output != "v"+pkg.VERSION+"\n" { // Update with expected version string
				t.Fatalf("Expected output to contain 'v%s', got: %s", pkg.VERSION, output)
			}

			exitCode := cmd.ProcessState.ExitCode()
			if exitCode != 0 {
				t.Fatalf("Expected exit code 0, got %d", exitCode)
			}
		})
	}
}

// When an unknown argument is passed, the program should exit with code 1
func TestExitCodeForUnknownFlags(t *testing.T) {
	t.Parallel()
	flags := []string{"-ø", "--anUnknownFlag", "--", "foobar"}
	for _, flag := range flags {
		f := flag
		t.Run("flag "+f, func(t *testing.T) {
			t.Parallel()
			cmd := exec.Command("go", "run", "../../main.go", f)
			err := cmd.Run()
			if err == nil {
				t.Fatalf("Expected program to exit with error for flag %s, but it didn't", f)
			}

			exitCode := cmd.ProcessState.ExitCode()
			if exitCode != 1 {
				t.Fatalf("Expected os.Exit(1) to be called for flag %s, but it wasn't", f)
			}
		})
	}
}

// When the config flag is specified, the option should be read
func TestParseCliArgsForConfig(t *testing.T) {
	reset()
	flags := []string{"-c", "--config"}
	for _, flag := range flags {
		t.Run("flag "+flag, func(t *testing.T) {
			reset()
			xConfigFilePath := "my/custom/path"
			os.Args = []string{"", flag, xConfigFilePath}
			parseCliArgs()
			configFilePath := viper.GetString("config")
			if configFilePath != xConfigFilePath {
				t.Fatalf("The parsed config file path '%s' does not correspond to the expected one '%s'", configFilePath, xConfigFilePath)
			}
		})
	}
}
