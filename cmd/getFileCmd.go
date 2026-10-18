package cmd

import (
	"os"

	"github.com/johnnewcombe/System3-disk-util/disk"
	"github.com/spf13/cobra"
)

var getFileCmd = &cobra.Command{
	Use:   "get-file",
	Short: "Gets a file from the specified Acorn Flex single sided 35, 40 or 80 track disk image.",
	Long:  `Deletes a file from the specified Acorn Flex single sided 35, 40 or 80 track disk image.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var (
			disk          disk.Disk
			data          []byte
			diskFilename  string
			localFilename string
			flexFilename  string
			err           error
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

		if data, err = disk.GetFile(flexFilename); err != nil {
			return err
		}

		if err = os.WriteFile(localFilename+".bin", data, 0644); err != nil {
			return err
		}

		/*
			// ===================================================================================
			// TODO temp code - remove
			// ===================================================================================


			tmpFileName := "BLANK_FORMATTED_35.img"
			tmpSize := 0x15e00

				var tmpData = []byte{}

				data, _ = os.ReadFile(tmpFileName)
				for i, b := range data {
					if i >= tmpSize {
						break
					}
					tmpData = append(tmpData, b)
				}
				if err = os.WriteFile("NEW_"+tmpFileName, tmpData, 0644); err != nil {
					return err
				}
				// ===================================================================================
		*/
		return nil

	},
}
