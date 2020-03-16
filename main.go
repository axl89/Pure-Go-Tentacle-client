package main

func main() {
	/*
		sendStatus, sendError := SendFile("/tmp/gato.xml")
		sendStatus2, sendError2 := SendFile("/tmp/gata.xml")
		// Any error from tentacle?
		if !sendStatus || !sendStatus2 {

			printError(sendError.Error())
			printError(sendError2.Error())
		}

		close()
	*/

	var c, _ = Dial("127.0.01", "41121")
	logEnabled = false
	c.SendFile("/tmp/gato.xml")
	c.SendFile("/tmp/gata.xml")
	c.Close()
}
