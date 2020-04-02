package main

import (
    "os"
    "io/ioutil"
    "fmt"
    "time"
)

func setNetStatus() {
    var newNetStatus string
    for {
        baseDir := "/sys/class/net/"
        netDirs, err := ioutil.ReadDir(baseDir)
        if err != nil {
            fmt.Fprintln(os.Stderr, err.Error())
            break
        }

        newNetStatus = ""

        for _, netDir := range netDirs {
            state, err := ioutil.ReadFile(baseDir + netDir.Name() + "/operstate")
            if err != nil {
                fmt.Fprintln(os.Stderr, err.Error())
                continue
            }

            if len(state) == 3 {
                if _, err := os.Stat(baseDir + netDir.Name() + "/wireless"); os.IsNotExist(err) {
                    newNetStatus += "   "
                } else {
                    newNetStatus += "   "
                }
            }
        }

        if newNetStatus != "" {
            newNetStatus = newNetStatus[:len(newNetStatus)-3]
        }

        if newNetStatus != netStatus {
            netStatus = newNetStatus
            printBuffer()
        }

        time.Sleep(time.Second * 5)
    }
}
