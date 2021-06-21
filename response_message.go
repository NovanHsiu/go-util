package util

import (
	"log"
	"strconv"
	"strings"
)

var atParameter = "@param"

// ErrorCodeTable record errror code and message for response
var ErrorCodeTable = map[string][]string{
	// 200
	"200.1": {"20001", "successful operation"}, // successful operation
	// 400
	"400.1": {"40001", "missing parameter: @param"},                      // missing parameter
	"400.2": {"40002", "parameter type error: @param"},                   // Type error
	"400.3": {"40003", "parameter format error: @param"},                 // format error.
	"400.4": {"40004", "parameter structure error"},                      // structure error.
	"400.5": {"40005", "password and confirmed password dose not match"}, // password not matched error.
	"400.6": {"40006", atParameter},                                      // time error: start_time > end_time
	// 401
	"401.1": {"40101", "authentication error@param"},                                // Authentication error!
	"401.2": {"40102", "wrong password"},                                            // Wrong password.
	"401.3": {"40103", "login failed too many times, login prohibited ten minutes"}, // login failed too many times
	"401.4": {"40104", "your device not certified to using this service"},           // device no certified
	// 403
	"403.1": {"40301", "permission denied@param"}, // Permission denied!
	// 404
	"404.1": {"40401", "@param not found"}, // data not found.
	// 422
	"422.1": {"42201", "@param exists"}, // data exists.
	"422.2": {"42202", atParameter},     // mac is already bound to key
	// 500
	"500.1": {"50001", atParameter}, // Internal server error
	"500.2": {"50002", atParameter}, // db error
	"500.3": {"50003", atParameter}, // file upload error
	"500.4": {"50004", atParameter}, // confirm error
	"500.5": {"50005", atParameter}, // call cch api error
}

func setMessage(msg string, param string) string {
	return strings.ReplaceAll(msg, atParameter, param)
}

// CreateResponseParam create response of web api with parameter
func CreateResponseParam(key string, param string) map[string]interface{} {
	code, err := strconv.Atoi(ErrorCodeTable[key][0])
	if err != nil {
		return map[string]interface{}{
			"error_code": 500,
			"error_msg":  "code of ErrorCodeTable not integer, please fix it, util/error-message.go",
		}
	}
	errMsg := setMessage(ErrorCodeTable[key][1], param)
	if !strings.HasPrefix(key, "200") {
		log.Println("error message: ", code, errMsg)
	}
	return map[string]interface{}{
		"error_code": code,
		"error_msg":  errMsg,
	}
}

// CreateResponse create response of web api
func CreateResponse(key string) map[string]interface{} {
	return CreateResponseParam(key, "")
}
