package main

import (
    "net"
    "time"
    "strconv"
    "fmt"
    "os"
    "io"
)

var bspwmStatusFromSocket string
var bspwmReadStatus int

func openBspwmSocket() (net.Conn, error) {
    sock, err := net.Dial("unix", "/tmp/bspwm_0_0-socket")
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return nil, err
    }
    return sock, err
}

func readSocket(reader io.Reader) {
    buf := make([]byte, 100)
    num, _ := reader.Read(buf)
    statusString := string(buf[:num])
    if statusString != bspwmStatusFromSocket {
        bspwmStatusFromSocket = statusString
        bspwmReadStatus = 1
    } else {
        bspwmReadStatus = 2
    }
}

func setStatusFromSocket() {
    bspwmReadStatus = 0
    sock, err := openBspwmSocket()
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        bspwmReadStatus = -1
        return
    }
    go readSocket(sock)

    msg := []byte("wm\x00--get-status\x00")
    _, err = sock.Write(msg)
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        bspwmReadStatus = -1
        return
    }

    for {
        if bspwmReadStatus != 0 {
            return
        }
    }
}

func setBspwmStatus() {
    for {
        setStatusFromSocket()
        if bspwmReadStatus == 2 {
            time.Sleep(150 * time.Millisecond)
            continue
        } else if bspwmReadStatus == -1 {
            return
        }

        newBspwmStatus := ""

        wsIndx := 1
        for _, e := range bspwmStatusFromSocket {
            switch e {
            case 'f':
                wsIndx++
            case 'F':
                fallthrough
            case 'O':
                newBspwmStatus += " %{+u}  " + strconv.Itoa(wsIndx) + "  %{-u} |"
                wsIndx++
            case 'o':
                newBspwmStatus += "   " + strconv.Itoa(wsIndx) + "   |"
                wsIndx++
            }
        }

        /* Remove ending spacer */
        newBspwmStatus = newBspwmStatus[:len(newBspwmStatus)-1]

        bspwmStatus = newBspwmStatus
        printBuffer()
        time.Sleep(150 * time.Millisecond)
    }
}
