package main

import (
    "fmt"
    "bufio"
    "os"
    "time"
    "strconv"
)

func setMemoryString() {
    file, err := os.Open("/proc/meminfo")
    if logFatal(err) {
        return
    }
    defer file.Close()


    for {
        memTotal := make([]byte, 0)
        memAvailable := make([]byte, 0)
        _, err = file.Seek(0, 0)
        if logFatal(err) {
            break
        }
        reader := bufio.NewReader(file)

        /* Get MemTotal and MemAvailable values from /proc/meminfo */
        for {
            str, err := reader.ReadString(':')
            if err != nil {
                break
            }

            if str == "MemTotal:" {
                for {
                    rdrByte, _ := reader.ReadByte()
                    if rdrByte == 'k' {
                        break
                    }

                    if rdrByte != ' ' {
                        memTotal = append(memTotal, rdrByte)
                    }
                }
            }

            if str == "MemAvailable:" {
                for {
                    rdrByte, _ := reader.ReadByte()
                    if rdrByte == 'k' {
                        break
                    }

                    if rdrByte != ' ' {
                        memAvailable = append(memAvailable, rdrByte)
                    }
                }
            }

            /* Read to end of line */
            str, err = reader.ReadString('\n')
        }

        memTotalInt, _ := strconv.Atoi(string(memTotal))
        memAvailableInt, _ := strconv.Atoi(string(memAvailable))
        dec := 1 - float32(memAvailableInt) / float32(memTotalInt)
        newMemoryString := fmt.Sprintf("%.2f", dec)
        newMemoryString = " " + newMemoryString[2:] + "%"
        if newMemoryString != memoryString {
            memoryString = newMemoryString
            printBuffer()
        }
        time.Sleep(3 * time.Second)
    }
}
