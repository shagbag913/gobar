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
    /* Map of module names and their corresponding goroutine */
    moduleMap := map[string]interface{} {
        "time": setTimeString,
        "charge": setChargeString,
        "bspwm": setBspwmStatus,
        "net": setNetStatus,
        "date": setDateString,
        "volume": setVolumeString,
        "brightness": setBrightnessString,
        "used_memory": setMemoryString,
    }

    enabledModules := getConfValue("main;enabled_modules")

    for moduleString, moduleFunction := range moduleMap {
        if enabledModules == "" || strings.Contains(enabledModules, moduleString) {
            go moduleFunction.(func())()
        }
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
