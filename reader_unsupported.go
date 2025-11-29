//go:build !linux && !windows && !darwin

package gosmbios

// readSMBIOS returns an error for unsupported operating systems
func readSMBIOS() (*SMBIOS, error) {
	return nil, ErrUnsupportedOS
}
