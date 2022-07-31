//go:build darwin
// +build darwin

package winapi

//Set windows console title
func SetConsoleTitle(title string) (int, error) {
	return 0, nil
}
