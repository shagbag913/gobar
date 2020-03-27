package main

import (
    "time"
    "strconv"
)

func setTimeString() {
    for {
        tNow := time.Now()
        minute, hour, second := tNow.Minute(), tNow.Hour(), tNow.Second()

        /* 12 hour time */
        if use12HourTime == true {
            if hour > 12 {
                hour -= 12
            } else if hour == 0 {
                hour = 12
            }
        }

        newTimeString := strconv.Itoa(hour) + ":"
        if minute < 10 {
            newTimeString += "0"
        }
        newTimeString += strconv.Itoa(minute)

        if showSeconds {
            newTimeString += ":" + strconv.Itoa(second)
        }

        if newTimeString != timeString {
            timeString = newTimeString
            printBuffer()
        }

        if showSeconds {
            time.Sleep(1 * time.Second)
        } else {
            time.Sleep(time.Duration(60 - second) * time.Second)
        }
    }
}
