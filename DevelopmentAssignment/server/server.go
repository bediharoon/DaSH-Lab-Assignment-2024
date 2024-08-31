package main

import (
    "bytes"
    "encoding/json"
    "log"
	"sync"
	"time"

	"github.com/bediharoon/dashassign/apireq"
	"github.com/bediharoon/dashassign/datamodel"

	"github.com/godbus/dbus/v5"
)

type serverData struct {
    mutex *sync.Mutex
    servedHashes *[]uint32
    userSerial *[]int64
    APIKey *string
    roundIndex *int
    connection *dbus.Conn
    writeChan *chan int64
    errChan *chan int64
}

func setDifference(a []uint32, b []uint32) ([]uint32){
    x := make(map[uint32]bool)
    var diff []uint32

    for _, e := range b {
        x[e] = true
    }

    for _, e := range a {
        if _, ok := x[e]; !ok {
            diff = append(diff, e)
        }
    }

    return diff
}

func (global serverData) RegisterHashes(promptHashes []uint32) ([]uint32, int64, *dbus.Error) {
    global.mutex.Lock()

    promptHashes = setDifference(promptHashes, *global.servedHashes)
    *global.servedHashes = append(*global.servedHashes, promptHashes...)

    uid := time.Now().UnixMilli() // No need of user privacy since, all data is public
    *global.userSerial = append(*global.userSerial, uid)

    global.mutex.Unlock()

    return promptHashes, uid, nil
}

func (global serverData) PromptRequest(promptData []datamodel.RequestData) (*dbus.Error) {
    for _, e := range promptData {
        format := make(map[string]string)
        format["inputs"] = e.Prompt

        payload, err := json.Marshal(format)
        if err != nil {
            log.Println("API Request Failed: ", err)
            e.Message = "API REQUEST FAILED" + err.Error()
        }

        e.Message, err = apireq.LLMRequest(bytes.NewReader(payload), *global.APIKey)
        if err != nil {
            log.Println("API Request Failed: ", err)
            e.Message = "API REQUEST FAILED" + err.Error()
        }

        e.Source = "gemma"

        AssignmentBlock:
        for _ = range len(*global.userSerial) {
            global.mutex.Lock()

            if len(*global.userSerial) <= 0 {
                log.Println("No Connected Clients")
                return nil
            } else if *global.roundIndex >= len(*global.userSerial) {
                *global.roundIndex = 0
            }

            e.ClientID = (*global.userSerial)[*global.roundIndex]
            *global.roundIndex += 1


            err = global.connection.Emit("/com/github/bediharoon/ServLLM", "com.github.bediharoon.ServLLM.NewResponse", e)
            if err != nil {
                global.mutex.Unlock()
                log.Println("Error Scheduling Response Write: ", err)
                continue
            }

            global.mutex.Unlock()

            ChannelResponseBlock:
            for {
                select {
                case x := <- *global.writeChan:
                    if x == e.ClientID {
                        break AssignmentBlock
                    }
                case x := <- *global.errChan:
                    if x == e.ClientID {
                        break ChannelResponseBlock
                    }
                case <- time.After(15 * time.Second):
                    *global.userSerial = append((*global.userSerial)[:*global.roundIndex], (*global.userSerial)[*global.roundIndex+1:]...)
                    break ChannelResponseBlock
                }
            }
        }
    }

    return nil
}

func (global serverData) WriteCheck(uid int64) (*dbus.Error) {
    *global.writeChan <- uid
    return nil
}

func (global serverData) WriteFail(uid int64) (*dbus.Error) {
    *global.errChan <- uid
    return nil

}
