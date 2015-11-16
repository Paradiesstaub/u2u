package main

import "strings"

// View object of main.
type View interface {
	SetDropdownItems(arr []string)
}

// SetDropdownItems sets a slice of items as QML ComboBox model.
func (b Bridge) SetDropdownItems(arr []string) {
	strArr := strings.Join(arr, ":")
	b.Root.Call("setDropdownItems", strArr)
}
