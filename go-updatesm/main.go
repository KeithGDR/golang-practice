package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"drixevel.dev/updatesm/internal/config"
	"drixevel.dev/updatesm/internal/sourcemod"
)

const (
	dirPrefix    = "D:/sourcemod/servers"     //Windows Only
	dirPrefixWSL = "/mnt/d/sourcemod/servers" //WSL Linux Only
	dirPostfix   = "addons/sourcemod/"
)

type Config struct {
	ServerDirectories  []string `json:"serverDirectories"`
	SourcemodVersion   string   `json:"smVersion"`
	ReplaceDirectories []string `json:"replaceDirectories"`
}

func copyDirRecursively(src, dst string) error {
	// If Windows, use the "copy" command.
	if runtime.GOOS == "windows" {
		return exec.Command("xcopy", src, dst+"\\", "/E").Run()
	} else if runtime.GOOS == "linux" {
		return exec.Command("/usr/bin/cp", "-r", src, dst).Run()
	}

	return fmt.Errorf("unsupported OS")
}

func getPath(path string) string {
	//Change the prefix of the path depending on if it's Windows or Linux.
	var prefix string
	if runtime.GOOS == "linux" {
		prefix = dirPrefixWSL
	} else {
		prefix = dirPrefix
	}

	//Format a valid path.
	return filepath.Join(prefix, path, path, dirPostfix)
}

func unzipSource(source, destination string) error {
	// 1. Open the zip file
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. Get the absolute destination path
	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	// 3. Iterate over zip files inside the archive and unzip each of them
	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
	}

	return nil
}

func unzipFile(f *zip.File, destination string) error {
	// 4. Check if file paths are not vulnerable to Zip Slip
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	// 5. Create directory tree
	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func ReplaceSourcemodFiles(config *Config, path string) error {
	var err error

	for _, i := range config.ReplaceDirectories {
		origpath := filepath.Join(sourcemod.LatestWindowsSMVersion, dirPostfix, i, "/")
		// replacepath := filepath.Join(path, i, "/")
		replacepath := filepath.Join(path)

		fmt.Println("Replacing " + origpath + " with " + replacepath)

		if err = copyDirRecursively(origpath, replacepath); err != nil {
			return err
		}
	}

	return err
}

func main() {
	//Get the config file.
	var cfg Config
	err := config.ParseConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	latestSmVersion, err := sourcemod.GetLatestSourceModVersion(cfg.SourcemodVersion)
	if err != nil {
		log.Fatal(err)
	}

	downloadUrl := fmt.Sprintf("%s/%s/%s", sourcemod.BaseSMDropURL, cfg.SourcemodVersion, latestSmVersion)

	fmt.Println("Download Sourcemod...")
	if err = DownloadFile(latestSmVersion, downloadUrl); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Downloaded " + downloadUrl)

	if err = unzipSource(latestSmVersion, sourcemod.LatestWindowsSMVersion); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Extracted " + latestSmVersion)

	wg := &sync.WaitGroup{}

	for _, i := range cfg.ServerDirectories {
		wg.Add(1)
		go func(dir string) {
			err := ReplaceSourcemodFiles(&cfg, getPath(dir))
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	//Clean up.
	if err := os.RemoveAll(latestSmVersion); err != nil {
		log.Fatal(err)
	}
	if err := os.RemoveAll(sourcemod.LatestWindowsSMVersion); err != nil {
		log.Fatal(err)
	}
}
