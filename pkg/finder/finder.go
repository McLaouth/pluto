package finder

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/klog"

	"github.com/fairwindsops/api-version-finder/pkg/api"
)

// App is the finder application
type App struct {
	RootPath *string
	FileList []*string
	APIFiles []*File
}

// File is a file that has an apiVersion in it
type File struct {
	Name       string
	APIVersion *api.Version
}

// NewFinder returns a new struct with config portions complete.
func NewFinder(path string) *App {
	cfg := &App{}
	if path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			klog.Fatal(err)
		}
		cfg.RootPath = &cwd
	} else {
		cfg.RootPath = &path
	}
	return cfg
}

// Execute runs the finder
func (app *App) Execute() error {
	err := app.findFiles()
	if err != nil {
		return err
	}
	err = app.scanFiles()
	if err != nil {
		return err
	}
	if app.APIFiles != nil {
		for _, file := range app.APIFiles {
			klog.Infof("%t - %s - %s", file.APIVersion.Deprecated, file.APIVersion.Name, file.Name)
		}
	}
	return nil
}

// FindFiles gets a list of all the files in the directory.
func (app *App) findFiles() error {
	var files []*string

	err := filepath.Walk(*app.RootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, &path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	app.FileList = files
	return nil
}

func (app *App) scanFiles() error {
	if app.FileList == nil {
		return fmt.Errorf("files list is nil")
	}

	for _, file := range app.FileList {
		klog.V(8).Infof("processing file: %s", *file)
		apiFile, err := checkForAPIVersion(*file)
		if err != nil {
			klog.Errorf("error scanning file %s: %s", file, err.Error())
		}
		if apiFile != nil {
			app.APIFiles = append(app.APIFiles, apiFile)
		}
	}
	return nil
}

func checkForAPIVersion(file string) (*File, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		version := api.StringContainsVersion(scanner.Text())
		if version != nil {
			apiFile := &File{
				Name:       file,
				APIVersion: version,
			}
			return apiFile, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}
