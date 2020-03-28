package main

import (
    "net"
    "regexp"
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
        bspwmReadStatus = 2
        return
    }
    go readSocket(sock)

    msg := []byte("wm\x00--get-status\x00")
    _, err = sock.Write(msg)
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        bspwmReadStatus = 2
        return
    }

    for {
        if bspwmReadStatus != 0 {
            return
        }
    }
}

func setBspwmStatus() {
    reg := regexp.MustCompile("[^oOfF]*")
    for {
        setStatusFromSocket()
        if bspwmReadStatus == 2 {
            time.Sleep(150 * time.Millisecond)
            continue
        }

        socketStatusFormatted := reg.ReplaceAllString(string(bspwmStatusFromSocket), "")

        /* Don't continue if socketStatus hasn't changed */
        if socketStatusFormatted == lastBspwmStatus {
            time.Sleep(150 * time.Millisecond)
            continue
        }
        lastBspwmStatus = socketStatusFormatted

        newBspwmStatus := ""

        for i := 0; i < len(socketStatusFormatted); i++ {
            switch socketStatusFormatted[i] {
            case 'F':
                fallthrough
            case 'O':
                newBspwmStatus += " %{+u}  " + strconv.Itoa(i+1) + "  %{-u} |"
            case 'o':
                newBspwmStatus += "   " + strconv.Itoa(i+1) + "   |"
            }
        }

        /* Remove ending spacer */
        newBspwmStatus = newBspwmStatus[:len(newBspwmStatus)-1]

        bspwmStatus = newBspwmStatus
        printBuffer()
        time.Sleep(150 * time.Millisecond)
    }
}
