package shift

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultTimeFormat = "20060102150405"
	defaultTimezone   = "UTC"
)

var (
	errInvalidSequenceWidth     = errors.New("digits must be positive")
	errIncompatibleSeqAndFormat = errors.New("the seq and format options are mutually exclusive")
	errInvalidTimeFormat        = errors.New("time format may not be empty")
)

var createCommand = &cobra.Command{
	Use:   "create [migration name]",
	Short: "Create new database migrations",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if name == "" {
			log.Fatal("please spesify database migration name")
		}

		timezone, err := time.LoadLocation(defaultTimezone)
		if err != nil {
			log.Fatal(err)
		}

		if err := createCmd(dir, time.Now().In(timezone), format, name, ext, seq, digit); err != nil {
			log.Print(err)
		}
	},
}

func init() {
	createCommand.PersistentFlags().StringVarP(&ext, "ext", "e", ".sql", "file extension")
	createCommand.PersistentFlags().IntVarP(&digit, "digits", "n", 6, "the number of digits to use in sequences")
	createCommand.PersistentFlags().BoolVarP(&seq, "sequence", "s", true, "use this option to generate sequential up/down migrations with N digits.")
	createCommand.PersistentFlags().StringVarP(&dir, "dir", "d", "migrations", "directory to place file")
	createCommand.PersistentFlags().StringVarP(&format, "format", "f", defaultTimeFormat, "use this option to specify a Go time format string.")
}

func createCmd(dir string, now time.Time, format string, name string, ext string, seq bool, seqDigits int) error {
	if seq && format != defaultTimeFormat {
		return errIncompatibleSeqAndFormat
	}

	var version string
	var err error

	dir = filepath.Clean(dir)
	ext = "." + strings.TrimPrefix(ext, ".")

	if seq {
		matches, err := filepath.Glob(filepath.Join(dir, "*"+ext))

		if err != nil {
			return err
		}

		version, err = nextSeqVersion(matches, seqDigits)

		if err != nil {
			return err
		}
	} else {
		version, err = timeVersion(now, format)

		if err != nil {
			return err
		}
	}

	versionGlob := filepath.Join(dir, version+"_*"+ext)
	matches, err := filepath.Glob(versionGlob)

	if err != nil {
		return err
	}

	if len(matches) > 0 {
		return fmt.Errorf("duplicate migration version: %s", version)
	}

	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	for _, direction := range []string{"up", "down"} {
		basename := fmt.Sprintf("%s_%s.%s%s", version, name, direction, ext)
		filename := filepath.Join(dir, basename)

		if err = createFile(filename); err != nil {
			return err
		}

	}

	return nil
}

func timeVersion(startTime time.Time, format string) (version string, err error) {
	switch format {
	case "":
		err = errInvalidTimeFormat
	case "unix":
		version = strconv.FormatInt(startTime.Unix(), 10)
	case "unixNano":
		version = strconv.FormatInt(startTime.UnixNano(), 10)
	default:
		version = startTime.Format(format)
	}

	return
}

func nextSeqVersion(matches []string, seqDigits int) (string, error) {
	if seqDigits <= 0 {
		return "", errInvalidSequenceWidth
	}

	nextSeq := uint64(1)

	if len(matches) > 0 {
		filename := matches[len(matches)-1]
		matchSeqStr := filepath.Base(filename)
		idx := strings.Index(matchSeqStr, "_")

		if idx < 1 { // Using 1 instead of 0 since there should be at least 1 digit
			return "", fmt.Errorf("malformed migration filename: %s", filename)
		}

		var err error
		matchSeqStr = matchSeqStr[0:idx]
		nextSeq, err = strconv.ParseUint(matchSeqStr, 10, 64)

		if err != nil {
			return "", err
		}

		nextSeq++
	}

	version := fmt.Sprintf("%0[2]*[1]d", nextSeq, seqDigits)

	if len(version) > seqDigits {
		return "", fmt.Errorf("next sequence number %s too large. At most %d digits are allowed", version, seqDigits)
	}

	return version, nil
}

func createFile(filename string) error {
	// create exclusive (fails if file already exists)
	// os.Create() specifies 0666 as the FileMode, so we're doing the same
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)

	if err != nil {
		return err
	}

	return f.Close()
}
