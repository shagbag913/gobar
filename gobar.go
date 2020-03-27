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

/* Other */
var lastBspwmStatus string

func main() {
    /* Initialize goroutines */
    go setTimeString()
    go setChargeString()
    go setBspwmStatus()
    go setNetStatus()

    /* Block main thread and let goroutines do everything */
    select { }
}

func printBuffer() {
    printBuffer := ""
    rightBuffer := "%{r}"

    if bspwmStatus != "" {
        printBuffer += "%{l}" + bspwmStatus
    }

    if timeString != "" {
        printBuffer += "%{c}" + timeString
    }

    if netStatus != "" {
        rightBuffer += netStatus + "   |   "
    }

    if chargeString != "" {
        rightBuffer += chargeString + "   |   "
    }

    if rightBuffer != "%{r}" {
        printBuffer += rightBuffer[:len(rightBuffer)-4]
    }

    if printBuffer != "" {
        printBuffer = printBuffer
        fmt.Println(printBuffer)
    }
}
