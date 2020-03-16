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

// Log messages, true enabled, false disabled
var logEnabled = true

//Do not output error messages, true enabled, false disabled
var quiet = false

//Client abstracts a Tentacle client
type Client struct {
	conn        *net.Conn
	writeBuffer *bufio.Writer
	readBuffer  *bufio.Reader
}

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

//Dial a Tentacle server and connects to it
func Dial(serverAddress string, serverPort string) (*Client, error) {
	conn, err := net.Dial("tcp", serverAddress+":"+serverPort)
	if err != nil {
		printError(err.Error())
	}
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)
	c := Client{&conn, w, r}

	return &c, err
}

func (c Client) send(message string) (string, error) {
	// Send message
	wInt, wErr := c.writeBuffer.WriteString(message + "\n")
	c.writeBuffer.Flush()

	printLog("Written " + strconv.Itoa(wInt) + " bytes")

	if wErr != nil {
		printError(wErr.Error())
	}

	// wait for reply
	message, err := c.readBuffer.ReadString('\n')

	if err != nil {
		printError(err.Error())
	}

	// Return the message without the trailing newline
	return strings.TrimSuffix(message, "\n"), err
}

func (c Client) sendByte(data []byte) (string, error) {

	// Send message
	wInt, wErr := c.writeBuffer.Write(data)
	c.writeBuffer.Flush()

	printLog("Written " + strconv.Itoa(wInt) + " bytes")

	if wErr != nil {
		printError(wErr.Error())
	}

	// wait for reply
	message, err := c.readBuffer.ReadString('\n')

	if err != nil {
		printError(err.Error())
	}

	// Return the message without the trailing newline
	return strings.TrimSuffix(message, "\n"), err
}

//Close the connection to the Tentacle server
func (c Client) Close() (bool, error) {
	// Tell Tentacle that we're leaving
	//fmt.Fprintf(conn, "QUITA\n")

	c.send("QUIT")

	conn := *c.conn
	error := conn.Close()

	return error != nil, error
}

// SendFile sends a file to a Tentacle server
func (c Client) SendFile(filePath string) (bool, error) {
	fi, err := os.Stat(filePath)

	if err != nil || fi == nil {
		c.Close()
		printError(err.Error())
		return false, err
	}
	// Get file size to string
	size := strconv.FormatInt(fi.Size(), 10)

	// Extract file name from the fill path
	fileName := path.Base(filePath)

	// Send request to Tentacle server
	message, sendErr := c.send("SEND <" + fileName + "> SIZE " + size)

	if message != "SEND OK" {
		printError(message)
		return false, sendErr
	}

	// Send data file
	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		return false, err
	}

	rMsg, rErr := c.sendByte(content)

	return rMsg == "SEND OK", rErr
}
