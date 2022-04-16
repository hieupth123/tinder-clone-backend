package util

import (
	"encoding/json"
	"fmt"
	"github.com/phamtrunghieu/tinder-clone-backend/constant"
	"strconv"
	"strings"
	"time"
)

func LogPrint(jsonData interface{}) {
	prettyJSON, _ := json.MarshalIndent(jsonData, "", "")
	fmt.Printf("%s\n", strings.ReplaceAll(string(prettyJSON), "\n", ""))
}

func LogInfo(jsonData ...interface{}) {
	cond := make(map[string]interface{})
	cond["message"] = "[INFO]"
	for index, log := range jsonData {
		cond[strconv.Itoa(index)] = log
	}
	LogPrint(cond)
}

func LogError(jsonData ...interface{}) {
	cond := make(map[string]interface{})
	cond["message"] = "[ERROR]"
	for index, log := range jsonData {
		cond[strconv.Itoa(index)] = log
	}
	LogPrint(cond)
}
func GetNowUTC() time.Time {
	loc, _ := time.LoadLocation("UTC")
	currentTime := time.Now().In(loc)
	return currentTime
}

func ConvertGender(gender string) string {
	if gender == "male" {
		return "mr"
	} else if gender == "female" {
		return "ms"
	} else {
		return ""
	}
}

func GetAge(dob string) int {
	dateOfBirth, errParse := time.Parse(constant.DATETIME_FORMAT_DB, dob)
	if errParse != nil {
		LogError(errParse)
		return 0
	}
	now := time.Now()
	age := int(now.Sub(dateOfBirth).Hours() / 24 / 365)
	return age
}
