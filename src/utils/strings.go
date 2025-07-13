package utils

func IsStrEmpty(target string) bool {
	return target == ""
}

func IsNotStrEmpty(target string) bool {
	return !IsStrEmpty(target)
}
