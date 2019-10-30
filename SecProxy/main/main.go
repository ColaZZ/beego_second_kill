package main

import _ "SecondKill/SecProxy/router"

func main() {
	err := initConfig()
	if err != nil {
		panic(err)
		return
	}

	err = initSec()
	if err != nil {
		panic(err)
		return
	}
}
