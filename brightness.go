package main

import (
    "time"
    "fmt"
    "os"
    "io/ioutil"
)

func setBrightnessString() {
    volFilePath := os.Getenv("HOME") + "/.cache/brightness/percentage"

    for {
        /* Fetch brightness from cache file */
        brightnessPercentageFromFile, err := ioutil.ReadFile(volFilePath)
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }

        newBrightnessString := " " + string(brightnessPercentageFromFile) + "%"
        if newBrightnessString != brightnessString {
            brightnessString = newBrightnessString
            printBuffer()
        }
        time.Sleep(300 * time.Millisecond)
    }
}
