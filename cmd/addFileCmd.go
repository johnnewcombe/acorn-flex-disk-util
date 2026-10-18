package cmd

import (
	"fmt"
	"os"

	"github.com/johnnewcombe/System3-disk-util/disk"
	"github.com/johnnewcombe/System3-disk-util/globals"
	"github.com/spf13/cobra"
)

var addBinaryFileCmd = &cobra.Command{
	Use:   "add-binary-file",
	Short: "Adds a raw binary or Flex binary file to the specified Acorn Flex disk image.",
	Long: `Adds a binary file to the specified Acorn Flex single sided 35, 40 or 80 track disk image. 
raw binary files are automatically converted to Flex CMD format files.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			disk            disk.Disk
			data            []byte
			diskFilename    string
			err             error
			localFilename   string
			flexFilename    string
			catalogue       string
			diskImage       []byte
			fileData        []byte
			loadAddress     uint16
			transferAddress uint16
			hasTransfer     bool
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

		if loadAddress, err = cmd.Flags().GetUint16("load-address"); err != nil {
			return err
		}
		if transferAddress, err = cmd.Flags().GetUint16("transfer-address"); err != nil {
			return err
		}
		hasTransfer = cmd.Flags().Changed("transfer-address")

		// get diskimage as data from storage
		if diskImage, err = os.ReadFile(diskFilename); err != nil {
			return err
		}

		// load up the disk object
		if err = disk.Load(diskImage); err != nil {
			return err
		}

		// get the file data
		if fileData, err = os.ReadFile(localFilename); err != nil {
			return err
		}

		// add the file to the disk
		if err = disk.AddBinaryFile(flexFilename, fileData, loadAddress,
			transferAddress, hasTransfer); err != nil {
			return err
		}

		// The file FLEX9.SYS is a special case and requires that the boot sector is linked to it.
		if flexFilename == globals.BootSectorLink {
			if err = disk.LinkBootSector(flexFilename); err != nil {
				return fmt.Errorf("unable to link boot sector: %w", err)
			}
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
