package main

import "log"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func logError(err error) {
	if err != nil {
		log.Print(err)
	}
}

func Find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
