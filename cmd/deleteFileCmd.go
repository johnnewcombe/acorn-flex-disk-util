package cmd

import (
	"fmt"
	"os"

	"github.com/johnnewcombe/System3-disk-util/disk"
	"github.com/spf13/cobra"
)

var deleteFileCmd = &cobra.Command{
	Use:   "delete-file",
	Short: "Deletes a file from the specified Acorn Flex single sided 35, 40 or 80 track disk image.",
	Long:  `Deletes a file from the specified Acorn Flex single sided 35, 40 or 80 track disk image.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			disk          disk.Disk
			data          []byte
			diskFilename  string
			localFilename string
			flexFilename  string
			err           error
			catalogue     string
			diskImage     []byte
		)

		if diskFilename, err = GetDiskFileName(cmd); err != nil {
			return err
		}

		if localFilename, err = GetFilename(cmd); err != nil {
			return err
		}

		if flexFilename, err = getFlexFilename(localFilename); err != nil {
			return err
		}

		if diskImage, err = os.ReadFile(diskFilename); err != nil {
			return err
		}

		if err = disk.Load(diskImage); err != nil {
			return err
		}

		if err = disk.DeleteFile(flexFilename); err != nil {
			return err
		}

		if data, err = disk.ToByte(); err != nil {
			return err
		}

		if err = os.WriteFile(diskFilename, data, 0644); err != nil {
			return err
		}

		if catalogue, err = disk.Catalogue(); err != nil {
			return err
		}
		fmt.Print(catalogue)

		return nil

	},
}
