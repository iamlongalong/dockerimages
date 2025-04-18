package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type ComposeConfig struct {
	Services map[string]struct {
		Image string `yaml:"image"`
	} `yaml:"services"`
}

var saveCmd = &cobra.Command{
	Use:   "save [images...]",
	Short: "Save Docker images to tar files",
	RunE:  runSave,
}

func init() {
	rootCmd.AddCommand(saveCmd)
}

func runSave(cmd *cobra.Command, args []string) error {
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

	// 4. Interactive selection
	if interactive {
		selectedImages, err := selectImagesInteractively()
		if err != nil {
			return err
		}
		images = append(images, selectedImages...)
	}

	if len(images) == 0 {
		return fmt.Errorf("no images specified")
	}

	// Remove duplicates
	images = removeDuplicates(images)

	if mergeImages {
		return saveMergedImages(images)
	}
	return saveIndividualImages(images)
}

func readImageListFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var images []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if image := strings.TrimSpace(scanner.Text()); image != "" {
			images = append(images, image)
		}
	}
	return images, scanner.Err()
}

func readComposeFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config ComposeConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	var images []string
	for _, service := range config.Services {
		if service.Image != "" {
			images = append(images, service.Image)
		}
	}
	return images, nil
}

func selectImagesInteractively() ([]string, error) {
	output, err := exec.Command("docker", "images", "--format", "{{.Repository}}:{{.Tag}}").Output()
	if err != nil {
		return nil, err
	}

	allImages := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(allImages) == 0 {
		return nil, fmt.Errorf("no Docker images found")
	}

	var selected []string
	prompt := &survey.MultiSelect{
		Message:  "Select images to save:",
		Options:  allImages,
		Help:     "Use arrow keys to move, Space to select/unselect, Enter to confirm",
		PageSize: 15,
	}

	// Remove validator since MultiSelect already ensures at least one selection
	if err = survey.AskOne(prompt, &selected); err != nil {
		if err == terminal.InterruptErr {
			return nil, fmt.Errorf("operation cancelled")
		}
		return nil, fmt.Errorf("selection error: %v", err)
	}

	if len(selected) == 0 {
		return nil, fmt.Errorf("no images selected")
	}

	return selected, nil
}

// pullImageIfNotExists pulls the image if it doesn't exist locally
func pullImageIfNotExists(image string) error {
	// Check if image exists
	cmd := exec.Command("docker", "image", "inspect", image)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Pulling image %s...\n", image)
		args := []string{"pull"}
		if platform != "" {
			args = append(args, "--platform", platform)
		}
		args = append(args, image)
		pullCmd := exec.Command("docker", args...)
		pullCmd.Stdout = os.Stdout
		pullCmd.Stderr = os.Stderr
		return pullCmd.Run()
	}
	return nil
}

func saveMergedImages(images []string) error {
	// Pull images if they don't exist
	for _, image := range images {
		if err := pullImageIfNotExists(image); err != nil {
			return fmt.Errorf("failed to pull image %s: %v", image, err)
		}
	}

	outputFile := filepath.Join(outputDir, "images.tar")
	args := append([]string{"save", "-o", outputFile}, images...)
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func saveIndividualImages(images []string) error {
	for _, image := range images {
		// Pull image if it doesn't exist
		if err := pullImageIfNotExists(image); err != nil {
			return fmt.Errorf("failed to pull image %s: %v", image, err)
		}

		safeName := strings.ReplaceAll(strings.ReplaceAll(image, "/", "_"), ":", "_")
		outputFile := filepath.Join(outputDir, safeName+".tar")
		cmd := exec.Command("docker", "save", "-o", outputFile, image)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to save image %s: %v", image, err)
		}
	}
	return nil
}

func removeDuplicates(images []string) []string {
	seen := make(map[string]bool)
	result := []string{}
	for _, image := range images {
		if !seen[image] {
			seen[image] = true
			result = append(result, image)
		}
	}
	return result
}
