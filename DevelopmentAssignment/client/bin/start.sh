#!/usr/bin/env bash

export  DBUS_SESSION_BUS_ADDRESS=$(dbus-daemon --fork --session --print-address)
/opt/client/main /opt/client/input.txt /opt/client/output.json server:12345
