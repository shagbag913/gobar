package main

import (
    "time"
    "fmt"
    "os"
)

func setBrightnessString() {
    for {
        /* Fetch brightness from cache file */
        file, err := os.Open(os.Getenv("HOME") + "/.cache/brightness/percentage")
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }

        brightnessPercentageFromFile := make([]byte, 3)
        var num int
        num, err = file.Read(brightnessPercentageFromFile)
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }

        newBrightnessString := " " + string(brightnessPercentageFromFile[:num]) + "%"
        if newBrightnessString != brightnessString {
            brightnessString = newBrightnessString
            printBuffer()
        }
        time.Sleep(300 * time.Millisecond)
    }
}
