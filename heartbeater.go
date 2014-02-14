package main

func heartbeaters(spawn int) {
	for i := 0; i < spawn; i++ {
		go heartbeater()
	}
}

func heartbeater() {
	qio := createClient()
	qio.Open()
	select {}
}
