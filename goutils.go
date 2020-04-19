package goutils

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

// CopyFilePath used
func CopyFilePath(srcDir string, dstDir string, fileName string) {
	srcFilePath := filepath.Join(srcDir, fileName)
	destFilePath := filepath.Join(dstDir, fileName)
	err := CopyFile(srcFilePath, destFilePath)
	Check(err)
}

// CopyDirPath used
func CopyDirPath(srcDir string, destDir string, dirName string) {
	srcDirPath := filepath.Join(srcDir, dirName)
	destDirPath := filepath.Join(destDir, dirName)
	err := CopyDir(srcDirPath, destDirPath)
	Check(err)
}

// RunBashCommand used
func RunBashCommand(args ...string) {
	name := args[0]
	cmd := exec.Command(name, args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	Check(err)
}

// Check used
func Check(err error) {
	if err != nil {
		log.Println("Error : ", err.Error())
		os.Exit(1)
	}
}

// CopyFile used
// File copies a single file from src to dst
// https://blog.depado.eu/post/copy-files-and-directories-in-go
func CopyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}

// CopyDir copies a whole directory recursively
// https://blog.depado.eu/post/copy-files-and-directories-in-go
func CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcinfo os.FileInfo
	// checking if src exist
	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}
	// create directory
	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}
	// read source dir return list of file dirs
	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	// for each file dir, get the full path of the source and destination
	for _, fd := range fds {
		srcfp := path.Join(src, fd.Name())
		dstfp := path.Join(dst, fd.Name())
		// recurse into directories
		if fd.IsDir() {
			if err = CopyDir(srcfp, dstfp); err != nil {
				log.Println(err)
			}
		} else {
			if err = CopyFile(srcfp, dstfp); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}
