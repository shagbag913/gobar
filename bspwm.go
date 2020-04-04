package main

import (
    "net"
    "strconv"
    "fmt"
    "os"
    "strings"
)

var statusFromSocket string
var socketPath string

func setBspwmStatus() {
    /* Set socket path */
    setSocketPath()

    /* Start socket listener */
    readBspwmSocket()
}

func setSocketPath() {
    socketPath = os.Getenv("BSPWM_SOCKET")
    if socketPath == "" {
        socketPath = "/tmp/bspwm_0_0-socket"
    }
}

func openBspwmSocket() (net.Conn, error) {
    sock, err := net.Dial("unix", socketPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return nil, err
    }
    return sock, err
}

func readBspwmSocket() {
    sock, err := openBspwmSocket()
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }

    /* Subscribe to wm events */
    _, err = sock.Write([]byte("subscribe\x00"))
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }

    for {
        status := make([]byte, 100)
        num, err := sock.Read(status)
        if err != nil {
            continue
        }
        statusString := strings.Split(string(status[:num]), "\n")[0]
        if statusFromSocket != statusString {
            statusFromSocket = statusString
            printNewBspwmStatusBuffer()
        }
    }
}

func printNewBspwmStatusBuffer() {
    wsIndx := 1
    bspwmStatus = ""
    for _, e := range statusFromSocket {
        switch e {
        case 'f':
            wsIndx++
        case 'F':
            fallthrough
        case 'O':
            bspwmStatus += " %{+u}  " + strconv.Itoa(wsIndx) + "  %{-u} |"
            wsIndx++
        case 'o':
            bspwmStatus += "   " + strconv.Itoa(wsIndx) + "   |"
            wsIndx++
        }
    }

    /* Remove ending spacer */
    bspwmStatus = bspwmStatus[:len(bspwmStatus)-1]
    printBuffer()
}
