package main

func checkError(err error) {
	if err != nil {
		panic(err)
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
