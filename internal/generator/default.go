package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:generate ifacemaker -f *.go -o igenerator.go -i Generator -s ApiGenerator -p generator
type ApiGenerator struct {
}

func New() Generator {
	return &ApiGenerator{}
}

func (ag ApiGenerator) Generate() error {

	directories := []string{"pkg", "internal", "configs", "scripts"}

	files := map[string]string{
		"README.md":           getReadme(),
		"configs/config.yaml": getConfigyaml(),
		"configs/config.go":   getConfig(),
		"scripts/run.sh":      "# Add your scripts here",
		"main.go":             getMain(),
		".gitignore":          getGitignore(),
		".env":                "# Add environment variables here",
		"Dockerfile":          generateDockerfile(),
		"docker-compose.yml":  generateDockerCompose(),
		".dockerignore":       "config/config.yaml",
	}

	if err := genDirectories(directories); err != nil {
		return err
	}

	if err := createFile(files); err != nil {
		return err
	}

	if err := generateModule(); err != nil {
		return err
	}

	return tidy()
}

func generateModule() error {
	_, err := os.Open("go.mod")

	if !os.IsNotExist(err) {
		return nil
	}

	cmd := exec.Command("go", "mod", "init", getModulePath())

	return cmd.Run()
}

func getModulePath() string {
	workDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return githubIndex(workDir)
}

func githubIndex(path string) string {
	index := strings.Index(path, "github.com")
	if index == -1 {
		fmt.Println("String 'github.com' not found in path")
		return ""
	}

	return path[index:]
}

func tidy() error {
	cmd := exec.Command("go", "mod", "tidy")

	return cmd.Run()
}

func genDirectories(dirs []string) error {
	for _, dir := range dirs {
		err := os.MkdirAll("./"+dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// create files and write necessary text in it
func createFile(filesMap map[string]string) error {
	// Iterate over each key-value pair in the map
	for filename, content := range filesMap {
		// Attempt to open the file for reading
		_, err := os.Open(filename)
		// If the file doesn't exist
		if os.IsNotExist(err) {
			// Create the file
			file, err := os.Create(filename)
			if err != nil {
				// If an error occurs during file creation, return the error
				return err
			}
			// Defer the closing of the file until the function returns
			defer file.Close()

			// Write the content to the file
			_, err = file.WriteString(content)
			if err != nil {
				// If an error occurs during writing, return the error
				return err
			}
		} else if err != nil {
			// If there's an error other than "file doesn't exist", return it
			return err
		}
	}

	// If no errors occurred, return nil
	return nil
}

// getProjectName extracts project name from the current working directory
func getProjectName() string {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	// Get the last element of the directory path, which is usually the project name
	projectName := filepath.Base(cwd)

	return projectName
}

func getReadme() string {
	projectName := getProjectName()

	content := "# " + projectName + "\n\nWrite your project description here."

	return content
}

func getMain() string {
	configPath := getModulePath() + "/configs"
	content := fmt.Sprintf(`package main

import (
	"log"
	"flag"

	"%s"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config-path", "./configs/config.yaml", "Path to config file")
	flag.Parse()

	config, err := configs.ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(config.ServiceName)
}`, configPath)

	return content
}

func getGitignore() string {
	content := `# Configs
/configs/config.yaml
	
# Environment 
/.env`

	return content
}

func getConfigyaml() string {
	content := fmt.Sprintf(`# Put your configuration files here

serviceName: %s`, strings.ToUpper(getProjectName()))

	return content
}

func getConfig() string {
	content := `package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServiceName string ` + "`yaml:\"serviceName\"`" + `
}

func ReadConfig(path string) (*Config, error) {
	configFile, err := os.ReadFile(path)
	if err != nil {
			return nil, err
	}

	return readConfigFromFile(configFile)
}

func readConfigFromFile(fileBytes []byte) (*Config, error) {
	cfg := new(Config)
	if err := yaml.Unmarshal(fileBytes, cfg); err != nil {
			return nil, err
	}

	return cfg, nil
}`

	return content
}

func generateDockerCompose() string {
	content := fmt.Sprintf(`version: "3.9"

services:
    %s:
        container_name: api
        platform: linux/amd64
        build:
            context: .
            dockerfile: Dockerfile
        restart: unless-stopped
        volumes:
            - ./configs/config.yaml:/app/configs/config.yaml
        ports:
            - "8080:8080"
`, getProjectName())
	return content
}

func generateDockerfile() string {
	content := `# Use specific versions for base images
FROM golang:1.22.0-alpine3.19 AS builder

WORKDIR /app

# Copy only necessary files for module downloading
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application with optimized flags
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/main main.go

# Use a smaller base image for the final stage
FROM alpine:3.19

WORKDIR /app

# Copy built binary and configuration file
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Set the entry point with necessary parameters
ENTRYPOINT ["./main"]
`
	return content
}
