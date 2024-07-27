package utils

func GetShortId(value string) string {
	return ShortString(value, 10)
}

func ShortString(value string, length int) string {
	shortValue := value
	if len(value) > length {
		shortValue = value[:length]
	}
	return shortValue
}
