package main

import "satori/go.uuid"

func isEmpty(str string) bool {
	return len(str) == 0
}

func isUUID(str string) bool {
	_, err := uuid.FromString(str)

	return err == nil
}
