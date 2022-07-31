//go:build linux
// +build linux

package winapi

//Set windows console title
func SetConsoleTitle(title string) (int, error) {
	return 0, nil
}
