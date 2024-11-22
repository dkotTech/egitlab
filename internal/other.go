package internal

import "fmt"

func EncodeHyperlink(href string, text string) string {
	return fmt.Sprintf("\u001B]8;;%s\u001B\\%s\u001B]8;;\u001B\\", href, text)
}
