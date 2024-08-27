package main

import (
    "bytes"
    "log"

    "github.com/bediharoon/dashassign/apireq"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

const intro = `
<node>
	<interface name="com.github.bediharoon.servllm">
		<method name="servreq">
			<arg direction="out" type="s"/>
		</method>
	</interface>` + introspect.IntrospectDataString +
`</node> `

type BaseNode int 

func (base BaseNode) GetResp(inp []byte, reqToken string) (string, *dbus.Error) {

    retOut, err := apireq.LLMRequest(bytes.NewReader(inp), reqToken)
    if (err != nil) {
        return "", dbus.MakeFailedError(err)
    }

    return retOut, nil
}

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
        log.Fatal(err)
	}
	defer conn.Close()

    var base BaseNode

	conn.Export(base, "/com/github/bediharoon/servllm", "com.github.bediharoon.ServLLM")
	conn.Export(introspect.Introspectable(intro), "/com/github/bediharoon/servllm", "org.freedesktop.DBus.Introspectable")

	reply, err := conn.RequestName("com.github.bediharoon.ServLLM", dbus.NameFlagReplaceExisting)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
        log.Panicln("Name Already Taken")
	}

	log.Println("Listening on com.github.bediharoon.ServLLM")
	select {}
}
