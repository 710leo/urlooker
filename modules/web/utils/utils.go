package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/toolkits/str"

	"github.com/710leo/urlooker/modules/web/g"
)

func Getkey(idc string, sid int) string {
	keys := g.Config.MonitorMap[idc]
	count := len(keys)
	return keys[sid%count]
}

func IsIP(ip string) bool {
	if ip != "" {
		isOk, _ := regexp.MatchString(`^(\d{1,3}\.){3}\d{1,3}$`, ip)
		if isOk {
			return isOk
		}
	}
	return false
}

func ParseUrl(target string) (schema, host, port, path string) {
	targetArr := strings.Split(target, "//")

	schema = targetArr[0]
	url := strings.Split(targetArr[1], "/")
	addrArr := strings.Split(url[0], ":")
	if len(addrArr) == 2 {
		host = addrArr[0]
		port = addrArr[1]
	} else {
		host = url[0]
	}

	for _, seg := range url[1:] {
		path += ("/" + seg)
	}
	return
}

func KeysOfMap(m map[string]string) []string {
	keys := make([]string, len(m))
	i := 0
	for key, _ := range m {
		keys[i] = key
		i++
	}

	return keys
}

func EncryptPassword(raw string) string {
	return str.Md5Encode(g.Config.Salt + raw)
}

func CheckUrl(url string) error {
	if !strings.Contains(url, "https://") && !strings.Contains(url, "http://") {
		return fmt.Errorf("http or https is necessary")
	}
	if len(url) > 1024 {
		return fmt.Errorf("url is too long over 1024")
	}
	return nil
}

func TimeFormat(ts int64) string {
	t := time.Unix(ts, 0).Format("2006-01-02 15:04:05")
	arr := strings.Split(t, " ")
	t = arr[1]
	arr = strings.Split(t, ":")

	return fmt.Sprintf("%s:%s", arr[0], arr[1])
}

func ReadLastLine(filename string) (string, error) {
	var previousOffset int64 = 0

	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	defer f.Close()

	reader := bufio.NewReader(f)

	// we need to calculate the size of the last line for file.ReadAt(offset) to work

	// NOTE : not a very effective solution as we need to read
	// the entire file at least for 1 pass :(

	lastLineSize := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		lastLineSize = len(line)
	}

	fileInfo, err := os.Stat(filename)
	if err != nil {
		return "", err
	}

	// make a buffer size according to the lastLineSize
	// +1 to compensate for the initial 0 byte of the line
	// otherwise, the initial character of the line will be missing
	// instead of reading the whole file into memory, we just read from certain offset

	l := int64(lastLineSize) * 30

	if fileInfo.Size() < l {
		l = fileInfo.Size() - 1
	}
	buffer := make([]byte, l)

	offset := fileInfo.Size() - int64(l+1)
	numRead, err := f.ReadAt(buffer, offset)
	if err != nil && err != io.EOF {
		return "", err
	}

	if previousOffset != offset {

		// print out last line content
		buffer = buffer[:numRead]
		fmt.Printf("%s \n", buffer)

		previousOffset = offset
	}
	//res := strings.Split(string(buffer), "\n")
	return string(buffer), nil
}
