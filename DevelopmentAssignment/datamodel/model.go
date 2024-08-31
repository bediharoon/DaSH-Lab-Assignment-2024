package datamodel

type RequestData struct {
    Prompt string
    Message string
    TimeSent int64
    TimeRecvd int64
    Source string
    ClientID int64 
    InitiatorUID int64 `json:"-"`
}
