package main

import (
    "time"
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

/* List of modules to enable for left, center, and right */
var enabledModules [3]string

/* List of default modules to enable for left, center, and right */
var defaultEnabledModules = [3]string{
    "bspwm",
    "time,date",
    "brightness,volume,used_memory,net,charge",
}

/* Left and right padding */
var leftPadding int
var rightPadding int

/* Buffer element seperator */
var elementSeperator = "|"

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

    updateConfigValues()

    for moduleString, moduleFunction := range moduleMap {
        for _, moduleSide := range enabledModules {
            if strings.Contains(moduleSide, moduleString) {
                go moduleFunction.(func())()
                break
            }
        }
    }

    /* Start config updated checker */
    go checkConfigUpdate()

    /* Block main thread and let goroutines do everything */
    select { }
}

func addToBuffer(element string, bufferIndex int) {
    if element != "" {
        buffers[bufferIndex] += element + "   " + elementSeperator + "   "
    }
}

func printBuffer() {
    moduleNameMap := map[string]string {
        "time": timeString,
        "charge": chargeString,
        "bspwm": bspwmStatus,
        "net": netStatus,
        "date": dateString,
        "volume": volumeString,
        "brightness": brightnessString,
        "used_memory": memoryString,
    }

    printBuffer := ""

    /* Reset buffers */
    for i := range buffers {
        buffers[i] = buffers[i][:4]
        if i == 0 {
            buffers[i] = buffers[i] + strings.Repeat(" ", leftPadding)
        }
    }

    for side, enabledModulesSide := range enabledModules {
        for _, module := range strings.Split(enabledModulesSide, ",") {
            addToBuffer(moduleNameMap[module], side)
        }
    }

    /*
     * A space is left between buffers to prevent status strings containing
     * percentage signs from conflicting from the side identifier
     */
    if strings.ReplaceAll(buffers[0], " ", "") != "%{l}" {
        printBuffer += buffers[0][:len(buffers[0])-(len(elementSeperator)+6)+1]
    }

    if buffers[1] != "%{c}" {
        printBuffer += buffers[1][:len(buffers[1])-(len(elementSeperator)+6)+1]
    }

    if buffers[2] != "%{r}" {
        printBuffer += buffers[2][:len(buffers[2])-(len(elementSeperator)+6)+1]
        printBuffer += strings.Repeat(" ", rightPadding)
    }

    if printBuffer != "" {
        fmt.Println(printBuffer)
    }
}

/* Returns true if a value was changed */
func updateConfigValues() bool {
    leftModules := getConfValue("main;modules_left", "")
    centerModules := getConfValue("main;modules_center", "")
    rightModules := getConfValue("main;modules_right", "")

    if leftModules == "" && centerModules == "" && rightModules == "" {
        leftModules = defaultEnabledModules[0]
        centerModules = defaultEnabledModules[1]
        rightModules = defaultEnabledModules[2]
    }

    paddingLeft := getConfInt("main;left_padding", 0)
    paddingRight := getConfInt("main;right_padding", 0)

    seperator := getConfValue("main;item_seperator", "|")

    if leftModules != enabledModules[0] || centerModules != enabledModules[1] ||
            rightModules != enabledModules[2] {
        enabledModules[0] = leftModules
        enabledModules[1] = centerModules
        enabledModules[2] = rightModules
        return true
    }

    if paddingLeft != leftPadding || paddingRight != rightPadding {
        leftPadding = paddingLeft
        rightPadding = paddingRight
        return true
    }

    if seperator != elementSeperator {
        elementSeperator = seperator
        return true
    }

    return false
}

func checkConfigUpdate() {
    for {
        if updateConfigValues() {
            printBuffer()
        }

        time.Sleep(time.Second * 1)
    }
}

