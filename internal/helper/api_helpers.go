package helper

import "database/sql"

type NullableString struct {
	sql.NullString
}

func LoopThroughFlags(givenMap map[string]any) bool {
	for _, flag := range givenMap {
		if flagsMap, ok := flag.(map[string]any); ok {
			if !LoopThroughFlags(flagsMap) {
				return false
			}
		} else if !flag.(bool) {
			givenMap["success"] = false
			return false
		}
	}
	givenMap["success"] = true
	return true
}
