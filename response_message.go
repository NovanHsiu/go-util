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
	"400.4": {"40004", "parameter structure error", "參數結構錯誤"},                    // structure error.
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

var errorCodeTableEnDescription = map[string]string{
	"200.1": "Operation successful",
	"400.1": "Missing parameter @param",
	"400.2": "Parameter type error @param",
	"400.3": "Parameter format error @param",
	"400.4": "Parameter structure error",
	"400.5": "Incorrect password",
	"400.6": "Incorrect time interval parameter arrangement",
	"400.7": "Pararameter's resource not found: @param",
	"401.1": "Login authentication failed",
	"401.2": "Incorrect password",
	"401.3": "Exceeded maximum login attempts, temporarily banned from logging in for ten minutes",
	"401.4": "Your device is not authorized to use this service",
	"403.1": "Insufficient permissions",
	"404.1": "@param not found",
	"422.1": "Already defined @param",
	"422.2": "MAC already bound to a key",
	"500.1": "Internal server error",
	"500.2": "Database error",
	"500.3": "File upload processing error",
	"500.4": "Signing error",
	"500.5": "Error calling HIS API",
}

var errorCodeTableZhtDescription = map[string]string{
	"200.1": "操作成功",
	"400.1": "缺少參數 @param",
	"400.2": "參數型別錯誤 @param",
	"400.3": "參數格式錯誤 @param",
	"400.4": "參數結構錯誤",
	"400.5": "密碼錯誤",
	"400.6": "時間區間參數安排錯誤",
	"400.7": "找不到 @param",
	"401.1": "登入認證失敗",
	"401.2": "密碼錯誤",
	"401.3": "登入失敗太多次，暫時禁止登入十分鐘",
	"401.4": "您的裝置未被認證可使用此服務",
	"403.1": "權限不足",
	"404.1": "找不到 @param",
	"422.1": "已定義 @param",
	"422.2": "MAC已經綁定金鑰",
	"500.1": "伺服器內部錯誤",
	"500.2": "資料庫錯誤",
	"500.3": "檔案上傳處理錯誤",
	"500.4": "送簽錯誤",
	"500.5": "呼叫HIS API錯誤",
}

var errorCodeTableZhsDescription = map[string]string{
	"200.1": "操作成功",
	"400.1": "缺少参数 @param",
	"400.2": "参数类型错误 @param",
	"400.3": "参数格式错误 @param",
	"400.4": "参数结构错误",
	"400.5": "密码错误",
	"400.6": "时间区间参数安排错误",
	"400.7": "找不到 @param",
	"401.1": "登录认证失败",
	"401.2": "密码错误",
	"401.3": "登录失败次数过多，暂时禁止登录十分钟",
	"401.4": "您的设备未被认证可使用此服务",
	"403.1": "权限不足",
	"404.1": "找不到 @param",
	"422.1": "已定义 @param",
	"422.2": "MAC已经绑定金钥",
	"500.1": "服务器内部错误",
	"500.2": "数据库错误",
	"500.3": "文件上传处理错误",
	"500.4": "送签错误",
	"500.5": "调用HIS API错误",
}

func setErrorCodeTable(langCode string) {
	if langCode == "en" {
		for key := range ErrorCodeTable {
			ErrorCodeTable[key][2] = errorCodeTableEnDescription[key]
		}
	} else if langCode == "zh-Hant" {
		for key := range ErrorCodeTable {
			ErrorCodeTable[key][2] = errorCodeTableZhtDescription[key]
		}
	} else if langCode == "zh-Hans" {
		for key := range ErrorCodeTable {
			ErrorCodeTable[key][2] = errorCodeTableZhsDescription[key]
		}
	}
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
		key = "500.1"
		code = 500
		param = "code of ErrorCodeTable not integer, please fix it, utils/error-message.go"
		desc = "API回傳錯誤處理功能錯誤"
	}
	if strings.HasPrefix(key, "2") {
		return map[string]interface{}{
			"error_code": code,
		}
	} else {
		errMsg := setMessage(ErrorCodeTable[key][1], param)
		log.Println("error message: ", code, errMsg)
		if desc == "" {
			desc = setMessage(ErrorCodeTable[key][2], param)
		}
		return map[string]interface{}{
			"error_code":  code,
			"error_msg":   errMsg,
			"description": desc,
		}
	}
}
