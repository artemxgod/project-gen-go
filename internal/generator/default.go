package generator

import (
	"fmt"
	"os"
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

	directories := []string{"cmd/api", "pkg", "internal", "configs", "scripts"}

	files := map[string]string{
		"README.md":           getReadme(),
		"configs/config.yaml": getConfigyaml(),
		"configs/config.go":   getConfig(),
		"scripts/run.sh":      "# Add your scripts here",
		"cmd/api/main.go":     getMain(),
		".gitignore":          getGitignore(),
		".env":                "",
	}

	if err := genDirectories(directories); err != nil {
		return err
	}

	if err := createFile(files); err != nil {
		return err
	}

	return nil
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
	content := "package main\n\nfunc main() {\n\n}"

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

ServiceName: %s`, strings.ToUpper(getProjectName()))

	return content
}

func getConfig() string {
	content := `package configs

	import (
			"os"
	
			"gopkg.in/yaml.v3"
	)
	
	type Config struct {
			ServiceName string 
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
