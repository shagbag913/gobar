package main

import (
    "time"
    "strconv"
)

func setDateString() {
    for {
        tNow := time.Now()
        year, month, day := tNow.Date()

        dateString = strconv.Itoa(int(month)) + "/" + strconv.Itoa(day) + "/" + strconv.Itoa(year)

        // Determine how long we can sleep until we have to refresh dateString
        sleepHours := time.Duration(23 - tNow.Hour()) * time.Hour
        sleepMinutes := time.Duration(59 - tNow.Minute()) * time.Minute
        sleepSeconds := time.Duration(60 - tNow.Second()) * time.Second
        totalSleep := sleepHours + sleepMinutes + sleepSeconds

        time.Sleep(totalSleep)
    }
}

func setTimeString() {
    for {
        tNow := time.Now()
        minute, hour, second := tNow.Minute(), tNow.Hour(), tNow.Second()

        /* 12 hour time */
        if getConfBool("time;12_hour_time") {
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

        showSeconds := getConfBool("time;show_seconds")
        if showSeconds {
            newTimeString += ":"
            if second < 10 {
                newTimeString += "0"
            }
            newTimeString += strconv.Itoa(second)
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
