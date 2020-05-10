// +build ios

package main

import "os"

func init() {
	_, ok := os.Lookup("ROOT_DIRECTORY")
	if !ok {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			os.Setenv("ROOT_DIRECTORY") = homeDir
		}
	}
}
