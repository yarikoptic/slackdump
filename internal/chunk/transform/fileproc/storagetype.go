package fileproc

import (
	"fmt"
	"strings"
)

type StorageType uint8

//go:generate stringer -type=StorageType -trimprefix=ST
const (
	STNone StorageType = iota
	STStandard
	STMattermost
)

// Set translates the string value into the ExportType, satisfies flag.Value
// interface.  It is based on the declarations generated by stringer.
func (e *StorageType) Set(v string) error {
	v = strings.ToLower(v)
	for i := 0; i < len(_StorageType_index)-1; i++ {
		if strings.ToLower(_StorageType_name[_StorageType_index[i]:_StorageType_index[i+1]]) == v {
			*e = StorageType(i)
			return nil
		}
	}
	return fmt.Errorf("unknown format: %s", v)
}
