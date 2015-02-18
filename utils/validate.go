package utils

import (
	"errors"
	"reflect"
	"strings"
)

func Validate(obj interface{}) *[]error {
	faults := []error{}
	ValidateStruct(&faults, obj)
	return &faults
}

func ValidateStruct(faults *[]error, obj interface{}) {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Skip ignored and unexported fields in the struct
		if !val.Field(i).CanInterface() {
			continue
		}

		fieldValue := val.Field(i).Interface()
		zero := reflect.Zero(field.Type).Interface()

		// Validate nested and embedded structs (if pointer, only do so if not nil)
		if field.Type.Kind() == reflect.Struct ||
			(field.Type.Kind() == reflect.Ptr && !reflect.DeepEqual(zero, fieldValue) &&
				field.Type.Elem().Kind() == reflect.Struct) {
			ValidateStruct(faults, fieldValue)
		}

		if strings.Index(field.Tag.Get("required"), "true") > -1 {
			if reflect.DeepEqual(zero, fieldValue) {

				name := field.Name
				if j := field.Tag.Get("json"); j != "" {
					name = j
				}

				*faults = append(*faults, errors.New(name+" is a required field"))
			}
		}
	}
}
