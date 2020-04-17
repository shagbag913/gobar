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
    seperator := getConfValue("bspwm;workspace_seperator", "|")
    for i, e := range statusFromSocket {
        if statusFromSocket[i:i+2] == "LT" || statusFromSocket[i:i+2] == "LM" {
            break
        }

        if !strings.Contains("FOfo", string(e)) {
            continue
        }

        wsIndxString := strconv.Itoa(wsIndx)

        if e == 'O' || e == 'F' {
            bspwmStatus += "  %{+u} " + strconv.Itoa(wsIndx) + " %{-u}  " + seperator
        } else if e == 'o' {
            bspwmStatus += "%{A:bspc desktop -f ^" + wsIndxString + ":}   " + wsIndxString
            bspwmStatus += "   %{A}" + seperator
        }
        wsIndx++
    }

    /* Remove ending spacer */
    bspwmStatus = bspwmStatus[:len(bspwmStatus)-len(seperator)]
    printBuffer()
}
