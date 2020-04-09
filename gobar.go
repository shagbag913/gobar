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

/* Left, center, and right print buffers */
var buffers = [3]string{"%{l}", "%{c}", "%{r}"}

/* Buffer element seperator */
var elementSeperator = "   |   "

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

func addToBuffer(element string, bufferIndex int) {
    if element != "" {
        buffers[bufferIndex] += element + elementSeperator
    }
}

func printBuffer() {
    printBuffer := ""

    /* clear buffers */
    for i := range buffers {
        buffers[i] = buffers[i][:4]
    }

    addToBuffer(bspwmStatus, 0)
    addToBuffer(timeString, 1)
    addToBuffer(dateString, 1)
    addToBuffer(memoryString, 2)
    addToBuffer(brightnessString, 2)
    addToBuffer(volumeString, 2)
    addToBuffer(netStatus, 2)
    addToBuffer(chargeString, 2)

    if buffers[0] != "%{l}" {
        printBuffer += buffers[0][:len(buffers[0])-len(elementSeperator)]
    }

    if buffers[1] != "%{c}" {
        printBuffer += buffers[1][:len(buffers[1])-len(elementSeperator)]
    }

    if buffers[2] != "%{r}" {
        printBuffer += buffers[2][:len(buffers[2])-len(elementSeperator)]
    }

    if printBuffer != "" {
        fmt.Println(printBuffer)
    }
}
