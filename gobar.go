package main

import (
    "fmt"
)

// CONFIGURATION FLAGS

/* Whether to use 12 or 24 hour time */
var use12HourTime bool = true

/* Whether to show seconds in time */
var showSeconds bool = false

/* Enable charging animation */
var animateChargeGlyphWhenCharging = true

// END CONFIGURATION FLAGS

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
    go setTimeString()
    go setChargeString()
    go setBspwmStatus()
    go setNetStatus()
    go setDateString()
    go setVolumeString()
    go setBrightnessString()
    go setMemoryString()

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
