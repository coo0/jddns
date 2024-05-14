package install

import (
	"errors"
	"fmt"
	"io"
	"jddns/lib/common"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func copyStaticFile(srcPath, bin string) string {
	binPath, _ := filepath.Abs(os.Args[0])
	if !common.IsWindows() {
		if _, err := copyFile(filepath.Join(srcPath, bin), "/usr/bin/"+bin); err != nil {
			if _, err := copyFile(filepath.Join(srcPath, bin), "/usr/local/bin/"+bin); err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		copyFile(filepath.Join(srcPath, bin+".exe"), filepath.Join(common.GetAppPath(), bin+".exe"))
	}
	chMod(binPath, 0755)
	return binPath
}
func Install() string {
	path := common.GetInstallPath()
	if common.FileExists(path) {

	} else {
		MkidrDirAll(path, "conf")
		// not copy config if the config file is exist
		if err := CopyDir(filepath.Join(common.GetAppPath(), "conf"), filepath.Join(path, "conf")); err != nil {
			log.Fatalln(err)
		}
		chMod(filepath.Join(path, "conf"), 0766)
	}
	binPath := copyStaticFile(common.GetAppPath(), "jddns")
	log.Println("install ok!")
	log.Println("Static files and configuration files in the current directory will be useless")
	log.Println("The new configuration file is located in", path, "you can edit them")
	if !common.IsWindows() {
		log.Println(`You can start with: jddns install|start|stop|uninstall anywhere!`)
	} else {
		log.Println(`You can start with: jddns install|start|stop|uninstall now!`)
	}
	chMod(common.GetLogPath(), 0777)
	return binPath
}
func Uninstall(bin string) {
	removeFile("/usr/bin/" + bin)
	removeFile("/usr/local/bin/" + bin)
	removeFile("/etc/" + bin)
}
func MkidrDirAll(path string, v ...string) {
	for _, item := range v {
		if err := os.MkdirAll(filepath.Join(path, item), 0755); err != nil {
			log.Fatalf("Failed to create directory %s error:%s", path, err.Error())
		}
	}
}

func CopyDir(srcPath string, destPath string) error {
	//检测目录正确性
	if srcInfo, err := os.Stat(srcPath); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if !srcInfo.IsDir() {
			e := errors.New("SrcPath is not the right directory!")
			return e
		}
	}
	if destInfo, err := os.Stat(destPath); err != nil {
		return err
	} else {
		if !destInfo.IsDir() {
			e := errors.New("DestInfo is not the right directory!")
			return e
		}
	}
	err := filepath.Walk(srcPath, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			destNewPath := strings.Replace(path, srcPath, destPath, -1)
			log.Println("copy file ::" + path + " to " + destNewPath)
			copyFile(path, destNewPath)
			if !common.IsWindows() {
				chMod(destNewPath, 0766)
			}
		}
		return nil
	})
	return err
}

// 生成目录并拷贝文件
func copyFile(src, dest string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFile.Close()
	//分割path目录
	destSplitPathDirs := strings.Split(dest, string(filepath.Separator))

	//检测时候存在目录
	destSplitPath := ""
	for index, dir := range destSplitPathDirs {
		if index < len(destSplitPathDirs)-1 {
			destSplitPath = destSplitPath + dir + string(filepath.Separator)
			b, _ := pathExists(destSplitPath)
			if b == false {
				log.Println("mkdir:" + destSplitPath)
				//创建目录
				err := os.Mkdir(destSplitPath, os.ModePerm)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}
	dstFile, err := os.Create(dest)
	if err != nil {
		return
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func removeFile(dest string) {
	err := os.RemoveAll(dest)
	if err != nil {
		fmt.Println("删除失败", err)
	} else {
		fmt.Println("删除成功 ", dest)
	}
}

// 检测文件夹路径时候存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func chMod(name string, mode os.FileMode) {
	if !common.IsWindows() {
		os.Chmod(name, mode)
	}
}
