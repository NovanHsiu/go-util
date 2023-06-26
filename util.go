package goutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// TimeLayout used to time.parse
const TimeLayout = "2006-01-02 15:04:05"
const randomText = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var ImageExtName = []string{".jpg", ".jpeg", ".png", ".bmp", ".gif"}

var moduleLanguage = "zh-Hant"
var moduleLanguageCodes = []string{"en", "zh-Hant", "zh-Hans"}

// SQLTimeFormatToString turn time to sql format time in string
func SQLTimeFormatToString(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

// SQLTimeStringToTime turn sql fromat time string to time type
func SQLTimeStringToTime(t string) (time.Time, error) {
	return time.Parse(TimeLayout, t)
}

// SQLTimeStringToTimeUTCZero turn sql fromat time string to time type and time zone is UTC+0
func SQLTimeStringToTimeUTCZero(t string) (time.Time, error) {
	_, offset := time.Now().Zone()
	offsetSecond, _ := time.ParseDuration(fmt.Sprintf("%ds", offset))
	ctime, err := time.Parse(TimeLayout, t)
	if err != nil {
		return ctime, err
	}
	return ctime.Add(-offsetSecond), nil
}

// TimeStamptoTime convert timestamp string to time
func TimeStamptoTime(timeStampStr string) (time.Time, error) {
	timeStampInt, err := strconv.ParseInt(timeStampStr, 10, 64)
	return time.Unix(timeStampInt/1000, (timeStampInt%1000)*1000000), err
}

// GetFormatTimeString set time to fomrat string
// example:
// - input: (time.Now(),"YYYY-MM-DD hh:mm:ss")
// - output: "2019-10-17 15:30:46"
func GetFormatTimeString(t time.Time, format string) string {
	timestring := format
	timestring = strings.ReplaceAll(timestring, "YYYY", fmt.Sprintf("%04d", t.Year()))
	timestring = strings.ReplaceAll(timestring, "MM", fmt.Sprintf("%02d", t.Month()))
	timestring = strings.ReplaceAll(timestring, "DD", fmt.Sprintf("%02d", t.Day()))
	timestring = strings.ReplaceAll(timestring, "hh", fmt.Sprintf("%02d", t.Hour()))
	timestring = strings.ReplaceAll(timestring, "mm", fmt.Sprintf("%02d", t.Minute()))
	timestring = strings.ReplaceAll(timestring, "ss", fmt.Sprintf("%02d", t.Second()))
	return timestring
}

// ContainsString Contains find element slice contains the element or not
func ContainsString(sl []string, v string) bool {
	return ArrayIncludeString(sl, v, true)
}

// HasString find the same string in array
func HasString(sl []string, v string) bool {
	return ArrayIncludeString(sl, v, false)
}

func ArrayIncludeString(sl []string, v string, islike bool) bool {
	for _, vv := range sl {
		if islike {
			if strings.Contains(v, vv) {
				return true
			}
		} else {
			if v == vv {
				return true
			}
		}
	}
	return false
}

// SetRequestBodyParams turn c.Request.Body ReaderCloser to params map[string]interface{}
func SetRequestBodyParams(body io.ReadCloser) (map[string]interface{}, error) {
	bodyBytes, _ := ioutil.ReadAll(body)
	params := make(map[string]interface{})
	err := json.Unmarshal(bodyBytes, &params)
	if err != nil {
		fmt.Println("error:", err)
	}
	return params, err
}

// GetMimeType get mimetype of file from *gin.Context.ContentType()
// example format of contentType: Content-Type:[text/plain]] multipart/form-data
func GetMimeType(contentType string) string {
	if !strings.Contains(contentType, "[") {
		return ""
	}
	return contentType[strings.Index(contentType, "[")+1 : strings.Index(contentType, "]")]
}

func ExecutionDir() string {
	exdir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return exdir
}

// GetRandomString 取得指定長度的隨機字串，長度必需大於0小於等於100，否則會回傳空字串 ""
func GetRandomString(length int) (rndtxt string) {
	rand.Seed(time.Now().UnixNano())
	if length < 1 || length > 100 {
		return rndtxt
	}
	for i := 0; i < length; i++ {
		rndtxt += string(randomText[rand.Intn(len(randomText))])
	}
	return rndtxt
}

// GetRandomTempDirName 取得指定長度且不重複的暫存資料夾路徑
func GetRandomTempDirName(basePath string, length int) (rndtxt string) {
	rndtxt = GetRandomString(length)
	for j := 0; j < 10; j++ {
		if _, err := os.Stat(basePath + "/" + rndtxt); os.IsNotExist(err) {
			break
		}
		time.Sleep(10 * time.Millisecond)
		rndtxt = GetRandomString(length)
	}
	return basePath + "/" + rndtxt
}

// PadLeft add pad string to the left of main string
func PadLeft(str, pad string, length int) string {
	for i := 0; i < length; i++ {
		str = pad + str
	}
	return str
}

// CopyFile copy a file from source to destination path
func CopyFile(src, dst string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()
	to, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	if err != nil {
		return err
	}
	return nil
}

func RemoveSuffixVersion(name string) string {
	eindex := strings.LastIndex(name, "v")
	if eindex <= 0 {
		return name
	}
	_, err := strconv.Atoi(name[eindex+1:])
	if err != nil {
		return name
	}
	if name[eindex-1] == ' ' || name[eindex-1] == '-' || name[eindex-1] == '_' {
		eindex--
	}
	return name[:eindex]
}

// MatchDatePattern check date string is match pattern "YYYY-MM-DD"，ex:"2015-11-26"
func MatchDatePattern(date string) bool {
	var validDate = regexp.MustCompile(`^[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])$`)
	return validDate.MatchString(date)
}

// ToRocDate change date to roc year format, ex: 1951-12-11 -> 401211
func ToRocDate(date string) string {
	dateSplit := strings.Split(date, "-")
	adyear, _ := strconv.Atoi(dateSplit[0])
	dateSplit[0] = strconv.Itoa(adyear - 1911)
	return strings.Join(dateSplit, "")
}

// RocToADYear ROC YYYMMDD(or YYMMDD) format to AD YYYY-MM-DD format
func RocToADYear(rocDate string) string {
	var yeindex int
	if len(rocDate) == 6 {
		yeindex = 2
	} else if len(rocDate) == 7 {
		yeindex = 3
	} else {
		return "0000-00-00"
	}
	rocyint, _ := strconv.Atoi(rocDate[:yeindex])
	year := strconv.Itoa(rocyint + 1911)
	month := rocDate[yeindex : yeindex+2]
	day := rocDate[yeindex+2 : yeindex+4]
	return fmt.Sprintf("%s-%s-%s", year, month, day)
}

// ConvertChineseGender convert "男", "女" to "M", "F"
func ConvertChineseGender(gender string) string {
	if gender == "男" {
		return "M"
	} else if gender == "女" {
		return "F"
	}
	return "Unknown"
}

// Execute execute terminal's command
// retrun stdout, stderr & error
func Execute(cmdstr string, args ...string) (string, string, error) {
	cmd := exec.Command(cmdstr, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err

}

// ExecuteBackground execute terminal's command in the background
// retrun pid, error
func ExecuteBackground(cmdstr string, args ...string) (int, error) {
	cmd := exec.Command(cmdstr, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		return 0, err
	}
	go func() {
		err = cmd.Wait()
		if err != nil {
			log.Println("ExecuteBackground error: " + cmdstr + fmt.Sprintf(" %v"+err.Error()))
			log.Printf("stdout: %q, stderr: %q\n", stdout.String(), stderr.String())
		}
	}()
	return cmd.Process.Pid, nil
}

// GetMacAddress get mac address of this device
func GetMacAddress() string {
	// for windows
	if runtime.GOOS == "windows" {
		stdout, _, _ := Execute("getmac")
		var macAddress string
		lineSplit := strings.Split(stdout, "\n")
		for _, line := range lineSplit {
			if strings.Contains(line, "Device") {
				macAddress = strings.Split(line, " ")[0]
				break
			}
		}
		return macAddress
	}
	// for liniux
	stdout, _, _ := Execute("cat", "/sys/class/net/eth0/address")
	return strings.Replace(stdout, "\n", "", 1)
}

// IsImage check filename is image
func IsImage(filename string) bool {
	lfn := strings.ToLower(filename)
	for i := range ImageExtName {
		if strings.HasSuffix(lfn, ImageExtName[i]) {
			return true
		}
	}
	return false
}

func ReadJSONConfig(path string) (map[string]interface{}, error) {
	dir := ExecutionDir()
	data, err := ioutil.ReadFile(dir + "/" + path)
	if err != nil {
		return nil, err
	}
	jsonMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(string(data)), &jsonMap)
	if err != nil {
		return nil, err
	}
	return jsonMap, nil
}

func WriteJSONConfig(path string, jsonMap map[string]interface{}) error {
	jsonBytes, err := json.MarshalIndent(jsonMap, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, jsonBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func RemoveDuplicateString(slice []string) []string {
	allKeys := make(map[string]bool)
	newSlice := []string{}
	for _, item := range slice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			newSlice = append(newSlice, item)
		}
	}
	return newSlice
}

// GetSystemLanguage get system's language setting
//
// Only support for detecting the following languages: en, zh-Hant, zh-Hans. (The returned names refer to the List of ISO 639-1 codes)
// All other languages will be treated as English.
func GetSystemLanguage() string {
	switch os := runtime.GOOS; os {
	case "windows":
		cmd := exec.Command("powershell", "-c", "Get-WinSystemLocale")
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("error:", err)
			return "en"
		}
		stdout := strings.ToLower(string(output))
		if strings.Contains(stdout, "zh-") {
			if strings.Contains(stdout, "zh-cn") || strings.Contains(stdout, "zh-sg") || strings.Contains(stdout, "zh-hans") {
				return "zh-Hans"
			} else {
				return "zh-Hant"
			}
		} else if strings.Index(stdout, "en-") == 0 {
			return "en"
		} else {
			return "en"
		}
	case "linux":
		return getLangFromUnixSystem()
	case "darwin":
		return getLangFromUnixSystem()
	default:
		fmt.Println("Others OS")
		return "en"
	}
}

func getLangFromUnixSystem() string {
	cmd := exec.Command("locale")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("error:", err)
		return "en"
	}
	stdout := strings.ToLower(string(output))
	stdout = strings.ReplaceAll(stdout, `"`, "")
	if strings.Index(stdout, "lang=zh") > 0 {
		if strings.Index(stdout, "lang=zh_cn") > 0 {
			return "zh-Hans"
		} else {
			return "zh-Hant"
		}
	} else if strings.Index(stdout, "lang=en") > 0 {
		return "en"
	} else {
		return "en"
	}
}

// SetModuleLanguage set module's language
//
// Only support the following languages: en, zh-Hant, zh-Hans.
func SetModuleLanguage(langCode string) error {
	if ArrayIncludeString(moduleLanguageCodes, langCode, false) {
		moduleLanguage = langCode
		setErrorCodeTable(moduleLanguage)
		return nil
	} else {
		return fmt.Errorf("language code error, must be %v", moduleLanguageCodes)
	}
}

func GetModuleLanguage() string {
	return moduleLanguage
}
