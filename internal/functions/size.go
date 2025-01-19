package functions

import "fmt"

const (
	kiloByte float64 = 1 << 10 // kilobyte in bytes
	megaByte float64 = 1 << 20 // megabyte in bytes
	gigaByte float64 = 1 << 30 // gigabyte in bytes
)

// FileSizeToString converts file size in bytes to a user friendly string.
func FileSizeToString(n int64) string {
	num := float64(n)

	switch {
	case num > gigaByte:
		return fmt.Sprintf("%.2f GB", num/gigaByte)
	case num > megaByte:
		return fmt.Sprintf("%.2f MB", num/megaByte)
	case num > kiloByte:
		return fmt.Sprintf("%.2f KB", num/kiloByte)
	default:
		return fmt.Sprintf("%.0f B", num)
	}
}
