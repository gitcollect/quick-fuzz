package main

func reconnectors(spawn int) {
	for i := 0; i < spawn; i++ {
		go reconnector()
	}
}

func reconnector() {
	for {
		qio := utilCreateClient()
		qio.Open()
		utilPause()
		qio.Close()
	}
}
