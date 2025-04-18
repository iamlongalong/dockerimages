# Docker Images Save Tool

A command-line tool for saving Docker images to tar files.

## Features

- Save multiple Docker images to tar files
- Multiple ways to specify images:
  - Command line arguments
  - Image list file (one image per line)
  - docker-compose.yml file
  - Interactive selection from available Docker images
- Option to merge multiple images into a single tar file or save them separately

## Installation

```bash
go install github.com/iamlongalong/dockerimages@latest
```

## Usage

### Basic Usage

Save a single image:
```bash
dockerimages save nginx:latest
```

Save multiple images:
```bash
dockerimages save nginx:latest redis:alpine mysql:8
```

### Using Image List File

Create a file with image names (one per line):
```
nginx:latest
redis:alpine
mysql:8
```

Then run:
```bash
dockerimages save -f images.txt
```

### Using docker-compose.yml

```bash
dockerimages save -c docker-compose.yml
```

### Interactive Mode

Select images interactively:
```bash
dockerimages save -i
```

### Merge Images

Save multiple images to a single tar file:
```bash
dockerimages save -m nginx:latest redis:alpine
```

### Specify Output Directory

```bash
dockerimages save -o /path/to/output nginx:latest
```

## Options

- `-o, --output`: Output directory for tar files (default: current directory)
- `-m, --merge`: Merge all images into a single tar file
- `-f, --file`: File containing image names (one per line)
- `-c, --compose`: Docker compose file path
- `-i, --interactive`: Interactive mode to select images 