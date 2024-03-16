package Compare

// IsStringInStringArray test if str string is contained in strArray []string
// return true when contains
func IsStringInStringArray(str string, strArray []string) bool {
	for _, v := range strArray {
		if v == str {
			return true
		}
	}
	return false
}
