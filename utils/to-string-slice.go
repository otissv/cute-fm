package utils

func ToStringSlice[T ~string](from []T) []string {
	slice := []string{}
	for _, col := range from {
		slice = append(slice, string(col))
	}

	return slice
}
