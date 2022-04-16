package respond

func ErrorCommonNotFound(message string) interface{} {
	return Respond{
		Message: message,
	}
}
