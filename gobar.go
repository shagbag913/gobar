package main

import (
    "strings"
    "fmt"
)

/* Printed bar strings */
var timeString string
var chargeString string
var bspwmStatus string
var netStatus string
var dateString string
var volumeString string
var brightnessString string
var memoryString string

func main() {
    /* Initialize goroutines */
    enabledModules := getConfValue("main;enabled_modules")

    if strings.Contains(enabledModules, "time") {
        go setTimeString()
    }
    if strings.Contains(enabledModules, "battery") {
        go setChargeString()
    }
    if strings.Contains(enabledModules, "bspwm") {
        go setBspwmStatus()
    }
    if strings.Contains(enabledModules, "net") {
        go setNetStatus()
    }
    if strings.Contains(enabledModules, "date") {
        go setDateString()
    }
    if strings.Contains(enabledModules, "volume") {
        go setVolumeString()
    }
    if strings.Contains(enabledModules, "brightness") {
        go setBrightnessString()
    }
    if strings.Contains(enabledModules, "used_memory") {
        go setMemoryString()
    }

    /* Block main thread and let goroutines do everything */
    select { }
}

func printBuffer() {
    printBuffer := ""
    rightBuffer := "%{r}"
    centerBuffer := "%{c}"

    if bspwmStatus != "" {
        printBuffer += "%{l}" + bspwmStatus
    }

    if timeString != "" {
        centerBuffer += timeString + "   |   "
    }

    if dateString != "" {
        centerBuffer += dateString + "   |   "
    }

    if memoryString != "" {
        rightBuffer += memoryString + "   |   "
    }

    if brightnessString != "" {
        rightBuffer += brightnessString + "   |   "
    }

    if volumeString != "" {
        rightBuffer += volumeString + "   |   "
    }

    if netStatus != "" {
        rightBuffer += netStatus + "   |   "
    }

    if chargeString != "" {
        rightBuffer += chargeString + "   |   "
    }

    if centerBuffer != "%{c}" {
        printBuffer += centerBuffer[:len(centerBuffer)-4]
    }

    if rightBuffer != "%{r}" {
        printBuffer += rightBuffer[:len(rightBuffer)-4]
    }

    if printBuffer != "" {
        printBuffer = printBuffer
        fmt.Println(printBuffer)
    }
}
