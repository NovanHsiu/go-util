package goutil

import (
	"log"
	"strconv"
	"strings"
)

var atParameter = "@param"

// ErrorCodeTable record errror code and message for response
var ErrorCodeTable = map[string][]string{
	// 200
	"200.1": {"20001", "successful operation", "操作成功"}, // successful operation
	// 400
	"400.1": {"40001", "missing parameter: @param", "缺少參數 @param"},               // missing parameter
	"400.2": {"40002", "parameter type error: @param", "參數型別錯誤 @param"},          // Type error
	"400.3": {"40003", "parameter format error: @param", "參數格式錯誤 @param"},        // format error.
	"400.4": {"40004", "parameter structure error", "參數格式錯誤"},                    // structure error.
	"400.5": {"40005", "password and confirmed password dose not match", "密碼錯誤"}, // password not matched error.
	"400.6": {"40006", atParameter, "時間區間參數安排錯誤"},                                // time error: start_time > end_time
	"400.7": {"40007", "pararameter's resource not found: @param", "找不到 @param"}, // resource defined by parameters not found
	// 401
	"401.1": {"40101", "authentication error@param", "登入認證失敗"},                                           // Authentication error!
	"401.2": {"40102", "wrong password", "密碼錯誤"},                                                         // Wrong password.
	"401.3": {"40103", "login failed too many times, login prohibited ten minutes", "登入失敗太多次，暫時禁止登入十分鐘"}, // login failed too many times
	"401.4": {"40104", "your device not certified to using this service", "您的裝置未被認證可使用此服務"},              // device no certified
	// 403
	"403.1": {"40301", "permission denied@param", "權限不足"}, // Permission denied!
	// 404
	"404.1": {"40401", "@param not found", "找不到 @param"}, // data not found.
	// 422
	"422.1": {"42201", "@param exists", "已定義 @param"}, // data exists.
	"422.2": {"42202", atParameter, "MAC已經綁定金鑰"},      // mac is already bound to key
	// 500
	"500.1": {"50001", atParameter, "伺服器內部錯誤"},     // Internal server error
	"500.2": {"50002", atParameter, "資料庫錯誤"},       // db error
	"500.3": {"50003", atParameter, "檔案上傳處理錯誤"},    // file upload error
	"500.4": {"50004", atParameter, "送簽錯誤"},        // confirm error
	"500.5": {"50005", atParameter, "呼叫HIS API錯誤"}, // call cch api error
}

func setMessage(msg string, param string) string {
	return strings.ReplaceAll(msg, atParameter, param)
}

// CreateResponseParam create response of web api with parameter
func CreateResponseParam(key string, param string) map[string]interface{} {
	return CreateResponseDesc(key, param, "")
}

// CreateResponse create response of web api
func CreateResponse(key string) map[string]interface{} {
	return CreateResponseParam(key, "")
}

// CreateResponseDesc create response of web api with parameter & description
func CreateResponseDesc(key, param, desc string) map[string]interface{} {
	code, err := strconv.Atoi(ErrorCodeTable[key][0])
	if err != nil {
		return map[string]interface{}{
			"error_code":  500,
			"error_msg":   "code of ErrorCodeTable not integer, please fix it, utils/error-message.go",
			"description": "API回傳錯誤處理功能錯誤",
		}
	}
	errMsg := setMessage(ErrorCodeTable[key][1], param)
	if !strings.HasPrefix(key, "200") {
		log.Println("error message: ", code, errMsg)
	}
	if desc == "" {
		desc = setMessage(ErrorCodeTable[key][2], param)
	}
	return map[string]interface{}{
		"error_code":  code,
		"error_msg":   errMsg,
		"description": desc,
	}
}
