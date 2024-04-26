package config

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSetupLogging_Valid_Config(t *testing.T) {
	t.Run("Additional log files", func(t *testing.T) {
		f, err := createTempFile()
		if err != nil {
			t.Fatalf("Error creating temporary file for logging: %v", err)
		}
		prevMessage := "Don't you dare forgetting me"
		f.WriteString(prevMessage + "\n")
		err = f.Sync()
		if err != nil {
			t.Fatal(fmt.Sprint("Failed to sync the file to disk", err))
		}
		defer removeTempFile(f)
		config := Config{Logging: loggingConfig{
			Level:              DefaultLoggingLevel,
			Style:              DefaultLoggingStyle,
			LogToStdout:        false,
			AdditionalLogFiles: []string{f.Name()},
		}}
		err = SetupLogging(config)
		if err != nil {
			t.Fatal(fmt.Sprint("No error was expected here, but we got one:", err))
		}
		message := "A random message for those who will read it"
		slog.Info(message)
		time.Sleep(time.Millisecond * 300)
		found := false
		foundOld := false
		//Read the log file from the start
		logFile, err := os.Open(f.Name())
		if err != nil {
			t.Fatal("Could not open file", err)
		}
		scanner := bufio.NewScanner(logFile)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			// The if statements are in this order because prevMessage should
			// appear before message in the log file
			if strings.Contains(line, prevMessage) {
				foundOld = true
			}
			if strings.Contains(line, message) {
				found = true
				break
			}
		}
		if err := scanner.Err(); err != nil {
			t.Fatal("Could not read additional log file", err)
		}
		if !foundOld {
			t.Fatal("The previous log in the log file was deleted")
		}
		if !found {
			t.Fatal("The log message was not written to the additional log file")
		}
	})
	t.Run("No stdout", func(t *testing.T) {
		config := Config{Logging: loggingConfig{
			Level:              DefaultLoggingLevel,
			Style:              DefaultLoggingStyle,
			LogToStdout:        false,
			AdditionalLogFiles: DefaultLoggingAdditionalLogFiles,
		}}
		err := SetupLogging(config)
		if err != nil {
			t.Fatal(fmt.Sprint("No error was expected here, but we got one:", err))
		}
	})
}

func TestSetupLogging_Invalid(t *testing.T) {
	t.Run("Log level", func(t *testing.T) {
		level := "incorrect"
		config := Config{Logging: loggingConfig{
			Level:              level,
			Style:              DefaultLoggingStyle,
			LogToStdout:        false,
			AdditionalLogFiles: DefaultLoggingAdditionalLogFiles,
		}}
		err := SetupLogging(config)
		if err == nil {
			t.Fatal("This operation should have failed, but it did not")
		}
	})
	t.Run("Log format", func(t *testing.T) {
		style := "incorrect"
		config := Config{Logging: loggingConfig{
			Level:              DefaultLoggingLevel,
			Style:              style,
			LogToStdout:        false,
			AdditionalLogFiles: DefaultLoggingAdditionalLogFiles,
		}}
		err := SetupLogging(config)
		if err == nil {
			t.Fatal("This operation should have failed, but it did not")
		}
	})
}
