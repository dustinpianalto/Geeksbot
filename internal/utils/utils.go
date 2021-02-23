package utils

func RemoveDuplicateStrings(s []string) []string {
	keys := make(map[string]bool)
	o := []string{}

	for _, e := range s {
		if _, v := keys[e]; !v {
			keys[e] = true
			o = append(o, e)
		}
	}
	return o
}

func PluralizeString(s string, i int) string {
	if i == 1 {
		return s
	}
	if s == "is" {
		return "are"
	}
	return s + "s"
}
