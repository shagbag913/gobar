package main

import (
    "strconv"
    "time"
    "fmt"
    "os"
)

func setBrightnessString() {
    intelBacklight := "/sys/class/backlight/intel_backlight/"
    brightnessFile, err := os.Open(intelBacklight + "brightness")
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }
    defer brightnessFile.Close()

    maxBrightnessFile, err := os.Open(intelBacklight + "max_brightness")
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }

    maxBrightnessRaw := make([]byte, 4)
    var num int
    num, err = maxBrightnessFile.Read(maxBrightnessRaw)
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }
    maxBrightnessFile.Close()

    var maxBrightness int
    maxBrightness, err = strconv.Atoi(string(maxBrightnessRaw[:num-1]))
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }

    for {
        _, err = brightnessFile.Seek(0, 0)
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }

        brightnessRaw := make([]byte, 4)
        num, err = brightnessFile.Read(brightnessRaw)
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }
        var brightnessInt int
        brightnessInt, err = strconv.Atoi(string(brightnessRaw[:num-1]))
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }
        brightnessDiv := float32(brightnessInt) / float32(maxBrightness)
        brightnessPercentage := ""
        if brightnessDiv != 1 {
            brightnessPercentage = fmt.Sprintf("%.2f", brightnessDiv)[2:]
        } else {
            brightnessPercentage = "100"
        }
        /* Remove beginning 0 from percentage string */
        if brightnessPercentage[0] == '0' {
            brightnessPercentage = brightnessPercentage[1:]
        }
        newBrightnessString := " " + brightnessPercentage + "%"
        if newBrightnessString != brightnessString {
            brightnessString = newBrightnessString
            printBuffer()
            time.Sleep(100 * time.Millisecond)
        } else {
            time.Sleep(300 * time.Millisecond)
        }
    }
}
