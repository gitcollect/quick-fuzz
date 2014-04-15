package main

func reconnectors(spawn int) {
	for i := 0; i < spawn; i++ {
		go reconnector()
	}
}

func reconnector() {
	ws := newWebSocket()

	for {
		ws.open()
		utilPause()
		ws.close()
	}
}
