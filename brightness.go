package main

import (
    "time"
    "fmt"
    "os"
)

func setBrightnessString() {
    volFilePath := os.Getenv("HOME") + "/.cache/brightness/percentage"

    for {
        /* Fetch brightness from cache file */
        file, err := os.Open(volFilePath)
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }
        defer file.Close()

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
