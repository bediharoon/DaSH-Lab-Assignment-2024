#!/usr/bin/env bash

export DBUS_SESSION_BUS_ADDRESS=$(dbus-daemon --fork --session --print-address)
/opt/serv/main localhost:12345
