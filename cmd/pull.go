package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull [images...]",
	Short: "Pull Docker images",
	Long: `Pull Docker images from registry.
It supports multiple ways to specify images:
- Command line arguments
- Image list file
- docker-compose.yml file
- Interactive selection from docker images`,
	RunE: runPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) error {
	var images []string

	// 1. Command line arguments
	if len(args) > 0 {
		images = append(images, args...)
	}

	// 2. Image list file
	if imagesFile != "" {
		fileImages, err := readImageListFile(imagesFile)
		if err != nil {
			return err
		}
		images = append(images, fileImages...)
	}

	// 3. Docker compose file
	if composeFile != "" {
		composeImages, err := readComposeFile(composeFile)
		if err != nil {
			return err
		}
		images = append(images, composeImages...)
	}

	if len(images) == 0 {
		return fmt.Errorf("no images specified")
	}

	// Remove duplicates
	images = removeDuplicates(images)

	// Pull all images
	for _, image := range images {
		fmt.Printf("Pulling image %s...\n", image)
		if err := pullImageIfNotExists(image); err != nil {
			return fmt.Errorf("failed to pull image %s: %v", image, err)
		}
		fmt.Printf("Successfully pulled %s\n", image)
	}

	return nil
}
