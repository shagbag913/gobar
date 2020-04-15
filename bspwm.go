package main

import (
    "net"
    "strconv"
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
    if logFatal(err) {
        return nil, err
    }
    return sock, err
}

func readBspwmSocket() {
    sock, err := openBspwmSocket()
    if logFatal(err) {
        return
    }

    /* Subscribe to wm events */
    _, err = sock.Write([]byte("subscribe\x00"))
    if logFatal(err) {
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
        if !strings.Contains("FOfo", string(e)) {
            continue
        }

        wsIndxString := strconv.Itoa(wsIndx)

        if e == 'O' || e == 'F' {
            bspwmStatus += " %{+u}  " + strconv.Itoa(wsIndx) + "  %{-u} |"
        } else if e == 'o' {
            bspwmStatus += "%{A:bspc desktop -f ^" + wsIndxString + ":}   " + wsIndxString + "   %{A}|"
        }
        wsIndx++
    }

    /* Remove ending spacer */
    bspwmStatus = bspwmStatus[:len(bspwmStatus)-1]
    printBuffer()
}
