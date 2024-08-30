package main

import (
	"log"
	"os"
	"sync"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

const intro = `
<node>
    <interface name="com.github.bediharoon.ServLLM">
        <method name="RegisterHashes">
            <arg direction="in" type="ax"/>
            <arg direction="out" type="ax"/>
            <arg direction="out" type="x"/>
        </method>
        <method name="PromptRequest">
            <arg direction="in" type="ax"/>
        </method>
        <method name="WriteCheck">
            <arg direction="in" type="x"/>
        </method>
        <method name="WriteFail">
            <arg direction="in" type="x"/>
        </method>
        <signal name="NewResponse">
            <arg direction="out" type="(ssxxsxx)"/>
        </signal>
    </interface>` + introspect.IntrospectDataString + 
`</node> `

func main() {
    conn, err := dbus.ConnectSessionBus()
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    if os.Getenv("HUG_API") == "" {
        log.Fatal(`Missing API Key, Use:
        $ export HUG_API="<YourTokenHere>"`)
    }

    var staticDataMutex sync.Mutex
    var preServedHashes []uint32
    var clientSerials []int64
    var currentIndex int
    apikey:= os.Getenv("HUG_API")
    successChan := make(chan int64)
    errorChan := make(chan int64)
    staticData := serverData{
        mutex: &staticDataMutex,
        servedHashes: &preServedHashes,
        userSerial: &clientSerials,
        APIKey: &apikey,
        roundIndex: &currentIndex,
        connection: conn,
        writeChan: &successChan,
        errChan: &errorChan,
    }
    conn.Export(staticData, "/com/github/bediharoon/ServLLM", "com.github.bediharoon.ServLLM")
    conn.Export(introspect.Introspectable(intro), "com/github/bediharoon/ServLLM", "org.freedesktop.DBus.Introspectable")

    reply, err := conn.RequestName("com.github.bediharoon.ServLLM", dbus.NameFlagDoNotQueue)
    if err != nil {
        log.Fatal(err)
    }

    if reply != dbus.RequestNameReplyPrimaryOwner {
        log.Fatalln("Name already taken")
    }

    log.Println("Listening on com.github.bediharoon.ServLLM")
    select { }
}
