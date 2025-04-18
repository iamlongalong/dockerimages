# Docker Images Save Tool

A command-line tool for saving Docker images to tar files and pulling images from registry.

## Features

- Save multiple Docker images to tar files
- Pull Docker images from registry
- Multiple ways to specify images:
  - Command line arguments
  - Image list file (one image per line)
  - docker-compose.yml file
  - Interactive selection from available Docker images
- Option to merge multiple images into a single tar file or save them separately
- Support for gzip compression of tar files

## Installation

```bash
go install github.com/iamlongalong/dockerimages@latest
```

## Usage

### Save Images

Save a single image:
```bash
dockerimages save nginx:latest
```

Save multiple images:
```bash
dockerimages save nginx:latest redis:alpine mysql:8
```

Save with gzip compression:
```bash
dockerimages save -z nginx:latest
```

### Pull Images

Pull a single image:
```bash
dockerimages pull nginx:latest
```

Pull multiple images:
```bash
dockerimages pull nginx:latest redis:alpine mysql:8
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
# To save images
dockerimages save -f images.txt

# To pull images
dockerimages pull -f images.txt
```

### Using docker-compose.yml

```bash
# To save images
dockerimages save -c docker-compose.yml

# To pull images
dockerimages pull -c docker-compose.yml
```

### Interactive Mode

Select images interactively:
```bash
# To save images
dockerimages save -i
```

### Save with Merge Option

Save multiple images to a single tar file:
```bash
dockerimages save -m nginx:latest redis:alpine
```

Save multiple images to a single compressed tar file:
```bash
dockerimages save -m -z nginx:latest redis:alpine
```

### Specify Output Directory (save only)

```bash
dockerimages save -o /path/to/output nginx:latest
```

## Options

### Global Options
- `-f, --file`: File containing image names (one per line)
- `-c, --compose`: Docker compose file path
- `-i, --interactive`: Interactive mode to select images
- `-p, --platform`: Target platform (e.g., linux/amd64, linux/arm64)

### Save Command Options
- `-o, --output`: Output directory for tar files (default: current directory)
- `-m, --merge`: Merge all images into a single tar file
- `-z, --gzip`: Compress the output tar file with gzip 