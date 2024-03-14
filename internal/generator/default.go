package generator

import (
	"os"
	"path/filepath"
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
		"configs/config.yaml": "# Put your configuration files here",
		"scripts/run.sh":      "# Add your scripts here",
		"cmd/api/main.go":     getMain(),
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
	for filename, content := range filesMap {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.WriteString(content)
		if err != nil {
			return err
		}
	}

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
	content := "package main\n\nfunc main(){\n\n}"

	return content
}
