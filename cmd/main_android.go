// +build android

package main

import "os"

func init() {
	_, ok := os.Lookup("ROOT_DIRECTORY")
	if !ok {
		os.Setenv("ROOT_DIRECTORY") = os.Getenv("FILESDIR")
	}
}
