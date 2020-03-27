package main

import (
    "regexp"
    "os/exec"
    "time"
    "strconv"
    "fmt"
    "os"
)

func setBspwmStatus() {
    reg := regexp.MustCompile("[^oOfF]*")

    for {
        /* TODO: talk directly to socket */
        socketStatus, err := exec.Command("bspc", "wm", "--get-status").Output()
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }
        socketStatusFormatted := reg.ReplaceAllString(string(socketStatus), "")

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
