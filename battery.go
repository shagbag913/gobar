package main

import (
    "strconv"
    "time"
    "os"
)

func getBatteryPercentGlyphIndex(batteryPercentage, overrideIndex int) int {
    switch {
    case overrideIndex >= 0:
        return overrideIndex
    case batteryPercentage >= 90:
        return 4
    case batteryPercentage >= 75:
        return 3
    case batteryPercentage >= 50:
        return 2
    case batteryPercentage >= 25:
        return 1
    default:
        return 0
    }
}

func getBatteryPercentWithGlyph(batteryPercentage, overrideIndex int, charging bool) string {
    batteryGlyphs := []string{"", "", "", "", ""}
    glyphIndex := getBatteryPercentGlyphIndex(batteryPercentage, overrideIndex)
    chargingString := batteryGlyphs[glyphIndex] + " " + strconv.Itoa(batteryPercentage) + "%"
    if charging {
        chargingString += "+"
    }
    return chargingString
}

func isCharging(file *os.File) bool {
    status := make([]byte, 12)
    var num int
    num, err := file.Read(status)
    if logFatal(err) {
        return false
    }
    if string(status[:num-1]) == "Discharging" {
        return false
    }

    return true
}

func setChargeString() {
    chargingIndexCounter := -1

    /* Open files */
    statusFile, err := os.Open("/sys/class/power_supply/BAT0/status")
    if logFatal(err) {
        return
    }
    defer statusFile.Close()

    var capacityFile *os.File
    capacityFile, err = os.Open("/sys/class/power_supply/BAT0/capacity")
    if logFatal(err) {
        return
    }
    defer capacityFile.Close()

    for {
        _, err = statusFile.Seek(0, 0)
        if logFatal(err) {
            break
        }
        _, err = capacityFile.Seek(0, 0)
        if logFatal(err) {
            break
        }
        charge := make([]byte, 4)
        var num int
        num, err = capacityFile.Read(charge)
        chargeInt, err := strconv.Atoi(string(charge[:num-1]))
        if logFatal(err) {
            break
        }

        isCharging := isCharging(statusFile)

        if !(isCharging && getConfBool("battery;animate_glyph_when_charging", false)) {
            /* Reset index counter */
            chargingIndexCounter = -1
        } else {
            if chargingIndexCounter == 4 || chargingIndexCounter < 0 {
                chargingIndexCounter = getBatteryPercentGlyphIndex(chargeInt, -1)
            } else {
                chargingIndexCounter++
            }
        }

        newChargeString := getBatteryPercentWithGlyph(chargeInt, chargingIndexCounter, isCharging)

        if newChargeString != chargeString {
            chargeString = newChargeString
            printBuffer()
        }

        time.Sleep(2 * time.Second)
    }
}
