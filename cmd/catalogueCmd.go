package cmd

import (
	"fmt"
	"os"

	"github.com/johnnewcombe/System3-disk-util/disk"
	"github.com/spf13/cobra"
)

var catalogueCmd = &cobra.Command{
	Use:   "cat",
	Short: "Displays disk and catalogue information.",
	Long:  `Displays disk and catalogue information.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			disk         disk.Disk
			diskFilename string
			err          error
			catalogue    string
			diskImage    []byte
		)

		if diskFilename, err = GetDiskFileName(cmd); err != nil {
			return err
		}

		if diskImage, err = os.ReadFile(diskFilename); err != nil {
			return err
		}

		if err = disk.Load(diskImage); err != nil {
			return err
		}

		if catalogue, err = disk.Catalogue(); err != nil {
			return err
		}
		fmt.Print(catalogue)

		return nil

	},
}
