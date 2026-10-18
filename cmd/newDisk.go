package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/johnnewcombe/System3-disk-util/disk"
	"github.com/spf13/cobra"
)

var newDiskCmd = &cobra.Command{
	Use:   "new-disk",
	Short: "Creates a blank Acorn Flex single sided 35, 40 or 80 track disk image.",
	Long:  `Creates a blank Acorn Flex single sided 35, 40 or 40 or 80 track disk image.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			dsk               *disk.Disk
			diskImageFilename string
			volumeNumber      int
			err               error
			volumeLabel       string
			tracks            int
			data              []byte
			catalogue         string
		)

		// get parameters
		if diskImageFilename, err = cmd.Flags().GetString("disk-image-filename"); err != nil {
			return err
		}

		if volumeLabel, err = cmd.Flags().GetString("volume-label"); err != nil {
			return err
		}

		if volumeNumber, err = cmd.Flags().GetInt("volume-number"); err != nil {
			return err
		}

		if tracks, err = cmd.Flags().GetInt("tracks"); err != nil {
			return err
		}

		// get parameters
		if diskImageFilename, err = cmd.Flags().GetString("disk-image-filename"); err != nil {
			return err
		}

		if dsk, err = disk.NewDisk(tracks, volumeLabel, volumeNumber); err != nil {
			return err
		}

		// write the file to dsk
		if path.Ext(diskImageFilename) != ".img" {
			diskImageFilename += ".img"
		}

		if data, err = dsk.ToByte(); err != nil {
			return err
		}

		if err = os.WriteFile(diskImageFilename, data, 0644); err != nil {
			return err
		}

		if catalogue, err = dsk.Catalogue(); err != nil {
			return err
		}
		fmt.Print(catalogue)

		return nil

	},
}
