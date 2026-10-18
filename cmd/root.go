package cmd

import (
	_ "embed"
	"errors"
	"os"

	"github.com/johnnewcombe/System3-disk-util/globals"
	"github.com/spf13/cobra"
)

const (
	//TDOD: Update the help text here and within the commands
	k_diskImageName    = "Name of the disk image file (*.img) to be used."
	k_tracks           = "Number of tracks for the single sided disk (35, 40 or 80)."
	k_volumeLabel      = "Volume Label of the disk as displayed in the catalogue."
	k_volumeNumber     = "Volume Number of the disk as displayed in the catalogue."
	k_filenameLocalAdd = "The name of the local raw binary or Flex binary file to be added."
	k_filenameFlexDel  = "The name of the flex filename to deleted from the disk."
	k_filenameFlexGet  = "The name of the flex filename to be retrieved from the disk."
	k_loadAddress      = "The address where the file should be loaded into memory."
	k_transferAddress  = "The transfer (execute) address."
)

func init() {

	//rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(newDiskCmd)
	rootCmd.AddCommand(addBinaryFileCmd)
	rootCmd.AddCommand(catalogueCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(deleteFileCmd)
	rootCmd.AddCommand(getFileCmd)
	rootCmd.AddCommand(validateCmd)

	// new-disk
	newDiskCmd.PersistentFlags().StringP("disk-image-filename", "d", "mydisk.img", k_diskImageName)
	newDiskCmd.PersistentFlags().IntP("tracks", "t", 80, k_tracks)
	newDiskCmd.PersistentFlags().StringP("volume-label", "v", "ACORNFLEX", k_volumeLabel)
	newDiskCmd.PersistentFlags().IntP("volume-number", "n", 1, k_volumeNumber)
	//newDiskCmd.PersistentFlags().IntP("boot-option", "b", 0, k_bootOption)

	// add-binary-file
	addBinaryFileCmd.Flags().StringP("disk-image-filename", "d", "", k_diskImageName)
	addBinaryFileCmd.Flags().StringP("filename", "f", "", k_filenameLocalAdd)
	addBinaryFileCmd.Flags().Uint16P("load-address", "l", 0x100, k_loadAddress)
	// NOTE the transfer address needs to be explicitly specified if required
	addBinaryFileCmd.Flags().Uint16P("transfer-address", "t", 0xFFFF, k_transferAddress)

	// cat
	catalogueCmd.Flags().StringP("disk-image-filename", "d", "", k_diskImageName)

	// delete-file
	deleteFileCmd.Flags().StringP("disk-image-filename", "d", "", k_diskImageName)
	deleteFileCmd.Flags().StringP("filename", "f", "", k_filenameFlexDel)

	// get-file
	getFileCmd.Flags().StringP("disk-image-filename", "d", "", k_diskImageName)
	getFileCmd.Flags().StringP("filename", "f", "", k_filenameFlexGet)

	// validate
	validateCmd.Flags().StringP("disk-image-filename", "d", "", k_diskImageName)

}

func Execute() {
	os.Exit(1)
}

var rootCmd = &cobra.Command{
	Use:   "disk-util",
	Short: "Acorn System DOS disk Utility version " + globals.Version,
	Long:  `Acorn System DOS disk Utility for Acorn System DOS. Note that this utility only supports single sided (.ssd) disk images.`,
}

func GetFilename(cmd *cobra.Command) (string, error) {

	var (
		filename string
		err      error
	)

	if filename, err = cmd.Flags().GetString("filename"); err != nil {
		return "", err
	}

	if len(filename) == 0 {
		return "", errors.New("filename not specified")
	}

	return filename, nil
}

func GetDiskFileName(cmd *cobra.Command) (string, error) {

	var (
		diskFilename string
		err          error
	)
	if diskFilename, err = cmd.Flags().GetString("disk-image-filename"); err != nil {
		return "", err
	}

	if len(diskFilename) == 0 {
		return "", errors.New("disk image not specified")
	}
	return diskFilename, nil
}
