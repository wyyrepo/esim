package file_dir

import (
	"os"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"github.com/jukylin/esim/log"
)


func CreateFile(file string) (bool, error) {
	f, err := os.Create(file)
	if err != nil{
		return false, err
	}else{
		f.Close()
		return true, nil
	}
}


func IsExistsDir(dir string)  (bool, error) {
	_, err := os.Stat(dir)
	if err == nil{
		return true, err
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func CreateDir(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil{
		return err
	}

	return nil
}


func IsExistsFile(file string) (bool, error) {
	if _, err := os.Stat(file); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func IsEmptyDir(dir string) (bool, error) {

	exists, err := IsExistsDir(dir)
	if err != nil{
		return false, err
	}

	if exists == false{
		return false, errors.New("目录不存在")
	}


	dirs, err := ioutil.ReadDir(dir)

	if err != nil{
		return false, err
	}

	if len(dirs) == 0{
		return true, nil
	}else {
		return false, nil
	}
}

func GetParDir() string {
	wd,err := os.Getwd()
	if err != nil {
		panic(err)
	}
	path := filepath.Dir(wd) + "/"
	parDir := strings.Replace(wd, path, "", -1)

	return parDir
}


func GetCurrentDir() string {
	wd,err := os.Getwd()
	if err != nil {
		panic(err)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == ""{
		panic(errors.New("not set GOPATH"))
	}

	srcPath := gopath + "/src/"

	parDir := strings.Replace(wd + "/", srcPath, "", -1)

	return parDir
}

//BackUpFile backup files to os.Getenv("GOPATH") + "/pkg/esim/"
// Overwrite as soon as the file exists
func BackUpFile(backFile string) error {

	if backFile == ""{
		return errors.New("没有文件")
	}

	dir := filepath.Dir(backFile)
	backUpPath := os.Getenv("GOPATH") + "/pkg/esim/"

	targetPath := backUpPath + dir

	exists, err := IsExistsDir(targetPath)
	if err != nil{
		return err
	}

	if exists == false {
		err = os.Mkdir(targetPath, 0777)
		if err != nil {
			return err
		}
	}

	_, err = CreateFile(backUpPath + backFile)
	if err != nil{
		return err
	}

	input, err := ioutil.ReadFile(backFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(backUpPath + backFile, input, 0644)
	if err != nil {
		return err
	}

	log.Log.Infof("%s backup to %s", backFile, backUpPath)

	return nil
}