package respond

func Success(data interface{}, message string) Respond {
	return Respond{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

func MissingParams() Respond {
	return Respond{
		Code:    1001,
		Message: "Missing params",
		Data:    nil,
	}
}

func ErrorResponse(message string) Respond {
	return Respond{
		Code:    1010,
		Message: message,
		Data:    nil,
	}
}
