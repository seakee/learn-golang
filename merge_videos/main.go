package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	envVarDeleteSuccess = "DELETE_SUCCESS" // 删除成功的环境变量
)

var (
	// 使用 flag 包定义命令行参数
	path      = flag.String("path", "./", "视频的保存目录")
	deleteSrc = flag.Bool("deleted", false, "删除已经合并的视频文件夹")
)

func main() {
	flag.Parse()

	fmt.Println(
		"小米摄像头视频合并工具启动成功\n",
		"视频储存的主目录：", *path,
		"合并后是否删除源文件夹：", getDeleteSrc(),
	)

	mergeHour(*path)
	fmt.Println("合并完成")
}

// getDeleteSrc 获取是否删除源文件夹的设置
func getDeleteSrc() bool {
	if os.Getenv(envVarDeleteSuccess) == "true" {
		return true
	}
	return *deleteSrc
}

// mergeHour 遍历视频文件的父目录，获取所有的小时子文件夹并进行合并处理
func mergeHour(pathDir string) {
	entries, err := os.ReadDir(pathDir)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历所有的子文件夹
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) == 10 {
			workDir := filepath.Join(pathDir, entry.Name())
			fmt.Println(workDir)

			// 如果成功合并并且设置了删除源文件夹，则删除源文件夹
			if mergeMp4ToMovByHour(workDir) && getDeleteSrc() {
				os.RemoveAll(workDir)
			}
		}
	}
}

// makeDirectory 在指定路径创建目录
func makeDirectory(pathOut string) bool {
	if err := os.MkdirAll(pathOut, 0777); err != nil {
		log.Println("创建目录失败:", pathOut, err.Error())
		return false
	}
	return true
}

// mergeMp4ToMovByHour 处理一个小时的视频文件，将其合并成一个 .mov 文件
func mergeMp4ToMovByHour(dirPath string) bool {
	var fileContent string
	_, pathName := filepath.Split(dirPath)
	year, month, day := pathName[0:4], pathName[4:6], pathName[6:8]

	outputPath := filepath.Join(*path, year, month, day)
	if ok := makeDirectory(outputPath); !ok {
		return false
	}

	// 将文件列表写入 fileContent
	writeFiles(dirPath, &fileContent)

	filesTxtPath := filepath.Join(dirPath, "files.txt")

	// 保存 fileContent 到文件
	SaveFileList(filesTxtPath, fileContent)
	defer os.Remove(filesTxtPath)

	mergeFileName := pathName[8:10] + ".mov"
	// 使用 ffmpeg 命令合并 .mp4 文件到 .mov 文件
	if ExecCommand(dirPath, "ffmpeg", "-f", "concat", "-i", "files.txt", "-c", "copy", mergeFileName) {
		// 将合并文件移到输出目录
		os.Rename(filepath.Join(dirPath, mergeFileName), filepath.Join(outputPath, mergeFileName))
		fmt.Println("保存到文件", filepath.Join(outputPath, mergeFileName))
		return true
	}
	return false
}

// writeFiles 将指定目录下的 .mp4 文件路径写入到 fileContent
func writeFiles(dirName string, fileContent *string) {
	filepath.WalkDir(dirName, func(path string, entry os.DirEntry, err error) error {
		if filepath.Ext(path) == ".mp4" {
			_, fileName := filepath.Split(path)
			*fileContent += "file " + fileName + "\n"
		}
		return nil
	})
}

// SaveFileList 保存文件列表到文件
func SaveFileList(filePath, content string) {
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		log.Println("写入文件失败:", filePath, err.Error())
	}
}

// ExecCommand 执行指定命令
func ExecCommand(workPath, name string, arg ...string) bool {
	cmd := exec.Command(name, arg...)
	cmd.Dir = workPath
	log.Printf("执行命令:%+v\n", cmd)

	// 获取命令的输出
	if output, err := cmd.Output(); err != nil {
		log.Println("指令执行错误:\n", err)
	} else {
		log.Println("指令执行结果输出:\n", output)
		return true
	}

	return false
}
