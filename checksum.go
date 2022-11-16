package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"math"
	"os"
	"path/filepath"
)

const HashCount = 3

var hashNames = []string{"md5", "sha1", "sha256"}

func main() {
	if len(os.Args) <= 1 || os.Args[1] == "" {
		println("use like \"checksum.exe F:/a.zip\"")
		return
	}
	path := os.Args[1]
	println("要处理的文件:", path)
	hashArr := [HashCount]hash.Hash{md5.New(), sha1.New(), sha256.New()}
	result, err := GetfileHash(path, &hashArr)
	if err != nil {
		return
	}

	writeToFile(path, result)
}

//goland:noinspection GoUnhandledErrorResult
func writeToFile(path string, hashValue *[HashCount]string) {
	_, file := filepath.Split(path)
	signFile := path + ".checksum"
	fp, err := os.Create(signFile)
	if err != nil {
		println("create file fail")
		return
	}
	defer fp.Close()
	fp.WriteString(fmt.Sprintf("%s\n", file))
	for index, hashName := range hashNames {
		s := fmt.Sprintf("%s=%s\n", hashName, hashValue[index])
		println(s)
		fp.WriteString(s)
	}
}

//goland:noinspection GoUnhandledErrorResult
func GetfileHash(path string, hashArr *[HashCount]hash.Hash) (*[HashCount]string, error) {
	fp, err := os.Open(path)
	if err != nil {
		fmt.Println("open file fail", err)
		return nil, err
	}
	defer fp.Close()
	fileStat, _ := fp.Stat()
	fmt.Printf("文件大小为 %d bytes\n", fileStat.Size())
	const BuffSize = 1024 ^ 2
	var processedSize int64 // 已处理的文件大小
	var nextAlter float64 = 0.05
	buf := make([]byte, BuffSize)
	for {
		n, _ := fp.Read(buf)
		if n == 0 {
			break
		}
		for _, hashUnit := range hashArr {
			if n == BuffSize {
				// full data
				hashUnit.Write(buf)
			} else {
				hashUnit.Write(buf[:n])
			}
		}
		processedSize += int64(n)
		if math.Abs(float64(processedSize) / float64(fileStat.Size()) - nextAlter) < 0.000001 {
			fmt.Printf("%d%%\n", int8(nextAlter*100))
			nextAlter += 0.05
		}
	}
	var result [HashCount]string
	for i := 0; i < HashCount; i++ {
		result[i] = hex.EncodeToString(hashArr[i].Sum(nil))
	}
	return &result, nil
}
