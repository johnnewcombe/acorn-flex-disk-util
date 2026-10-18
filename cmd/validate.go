package cmd

import (
	"fmt"
	"os"

	"github.com/johnnewcombe/System3-disk-util/disk"
	"github.com/spf13/cobra"
)

var validateStrict bool

var validateCmd = &cobra.Command{
	Use:     "validate <disk-image>",
	Aliases: []string{"verify", "check"},
	Short:   "Check the structure of a FLEX disk image",
	Long: `Checks a FLEX disk image for structural damage:

  - the image size and the geometry recorded in the system information record
  - the directory and every file's sector chain (links, record numbers,
    sector counts and end sectors against the directory)
  - the free chain against the system information record
  - sectors belonging to more than one chain, or to none
  - that .CMD and .BIN files contain valid FLEX binary records

Errors mean the disk is damaged and should not be written to. Warnings are
suspicious but not necessarily damage; use --strict to treat them as errors.
The exit status is non-zero if the disk fails.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			diskFilename string
			diskImage    []byte
			err          error
		)

		if diskFilename, err = GetDiskFileName(cmd); err != nil {
			return err
		}

		if diskImage, err = os.ReadFile(diskFilename); err != nil {
			return err
		}

		// from here on a failure is a finding about the disk, not a usage mistake
		cmd.SilenceUsage = true

		report := disk.ValidateImage(diskImage)
		out := cmd.OutOrStdout()

		fmt.Fprintf(out, "%s: %d tracks, %d files, %d free sectors\n",
			diskFilename, report.Tracks, report.Files, report.FreeSectors)
		for _, e := range report.Errors {
			fmt.Fprintf(out, "  ERROR    %s\n", e)
		}
		for _, w := range report.Warnings {
			fmt.Fprintf(out, "  WARNING  %s\n", w)
		}

		switch {
		case !report.OK():
			return fmt.Errorf("%s: %d error(s), %d warning(s)", diskFilename, len(report.Errors), len(report.Warnings))
		case validateStrict && len(report.Warnings) > 0:
			return fmt.Errorf("%s: %d warning(s) (--strict)", diskFilename, len(report.Warnings))
		}

		fmt.Fprintln(out, "OK")
		return nil
	},
}

func init() {
	validateCmd.Flags().BoolVar(&validateStrict, "strict", false, "treat warnings as errors")
	rootCmd.AddCommand(validateCmd)
}
