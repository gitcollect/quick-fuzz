package main

func reconnectors(spawn int) {
	for i := 0; i < spawn; i++ {
		go reconnector()
	}
}

func reconnector() {
	for {
		qio := createClient()
		qio.Open()
		util_pause()
		qio.Close()
	}
}
