package utils

import (
    "reflect"
)

// HasValue returns if the given interface has data beneath
func HasValue(value interface{}) bool {
    return value != nil && !reflect.ValueOf(value).IsNil() 
}