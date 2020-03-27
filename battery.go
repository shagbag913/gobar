package main

import (
    "io/ioutil"
    "strconv"
    "time"
)

func getBatteryPercentGlyphIndex(batteryPercentage, overrideIndex int) int {
    if overrideIndex >= 0 {
        return overrideIndex
    }

    if batteryPercentage >= 90 {
        return 4
    } else if batteryPercentage >= 75 {
        return 3
    } else if batteryPercentage >= 50 {
        return 2
    } else if batteryPercentage >= 25 {
        return 1
    } else {
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

func isCharging() bool {
    status, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/status")
    if err != nil {
        return false
    }

    if string(status[:len(status) - 1]) == "Discharging" {
        return false
    }

    return true
}

func setChargeString() {
    chargingIndexCounter := -1
    for {
        charge, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/capacity")
        if err != nil {
            break
        }
        chargeInt, _ := strconv.Atoi(string(charge[:len(charge)-1]))

        isCharging := isCharging()

        if !(isCharging && animateChargeGlyphWhenCharging) {
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
