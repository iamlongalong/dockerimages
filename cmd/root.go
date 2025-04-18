package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	outputDir    string
	mergeImages  bool
	imagesFile   string
	composeFile  string
	interactive  bool
	platform     string
	gzipCompress bool
)

var rootCmd = &cobra.Command{
	Use:   "dockerimages",
	Short: "A tool for saving Docker images to tar files",
	Long: `dockerimages is a CLI tool that helps you save Docker images to tar files.
It supports multiple ways to specify images:
- Command line arguments
- Image list file
- docker-compose.yml file
- Interactive selection from docker images`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", ".", "Output directory for tar files")
	rootCmd.PersistentFlags().BoolVarP(&mergeImages, "merge", "m", false, "Merge all images into a single tar file")
	rootCmd.PersistentFlags().StringVarP(&imagesFile, "file", "f", "", "File containing image names (one per line)")
	rootCmd.PersistentFlags().StringVarP(&composeFile, "compose", "c", "", "Docker compose file path")
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode to select images")
	rootCmd.PersistentFlags().StringVarP(&platform, "platform", "p", "", "Target platform (e.g., linux/amd64, linux/arm64)")
	rootCmd.PersistentFlags().BoolVarP(&gzipCompress, "gzip", "z", false, "Compress the output tar file with gzip")
}
