package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
)

var conn, connErr = net.Dial("tcp", "127.0.0.1:41121")
var w = bufio.NewWriter(conn)
var r = bufio.NewReader(conn)

// Log messages, true enabled, false disabled
var logEnabled = true

//Do not output error messages, true enabled, false disabled
var quiet = false

func printLog(data string) {

	if logEnabled {
		fmt.Println("[log] " + data)
	}
}

func printError(data string) {
	if !quiet {
		fmt.Fprintln(os.Stderr, "[err]:", data)
	}
}

func send(message string) (string, error) {
	// Send message
	wInt, wErr := w.WriteString(message + "\n")
	w.Flush()

	printLog("We wrote: " + strconv.Itoa(wInt))

	if wErr != nil {
		printError(wErr.Error())
	}

	// wait for reply
	message, err := r.ReadString('\n')

	if err != nil {
		printError(err.Error())
	}

	// Return the message without the trailing newline
	return strings.TrimSuffix(message, "\n"), err
}

func sendByte(data []byte) (string, error) {
	// Send message
	wInt, wErr := w.Write(data)
	w.Flush()

	printLog("We wrote: " + strconv.Itoa(wInt))

	if wErr != nil {
		printError(wErr.Error())
	}

	// wait for reply
	message, err := r.ReadString('\n')

	if err != nil {
		printError(err.Error())
	}

	// Return the message without the trailing newline
	return strings.TrimSuffix(message, "\n"), err
}

func close() (bool, error) {
	// Tell Tentacle that we're leaving
	//fmt.Fprintf(conn, "QUITA\n")
	send("QUIT")
	error := conn.Close()

	return error != nil, error
}

// SendFile sends a file to a Tentacle server
func SendFile(filePath string) (bool, error) {
	fi, err := os.Stat(filePath)

	if err != nil || fi == nil {
		close()
		printError(err.Error())
		return false, err
	}
	// Get file size to string
	size := strconv.FormatInt(fi.Size(), 10)

	// Extract file name from the fill path
	fileName := path.Base(filePath)

	// Send request to Tentacle server
	message, sendErr := send("SEND <" + fileName + "> SIZE " + size)

	if message != "SEND OK" {
		printError(message)
		return false, sendErr
	}

	// Send data file
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		return false, err
	}

	rMsg, rErr := sendByte(content)

	return rMsg == "SEND OK", rErr
}

func main() {
	sendStatus, sendError := SendFile("/tmp/gato.xml")
	sendStatus2, sendError2 := SendFile("/tmp/gata.xml")
	// Any error from tentacle?
	if !sendStatus || !sendStatus2 {

		printError(sendError.Error())
		printError(sendError2.Error())
	}

	close()
}
