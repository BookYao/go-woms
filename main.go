package main

import (
	"./eventmodule"
	"fmt"
	"log"
	"net/http"
)

const (
	Event_UP_Interface  = "/woms/api/UpEvent.php"
	Query_Single_Event  = "/woms/api/QuerySingleEvent.php"
	Event_Status_UpDate = "/woms/api/UpEventStatus.php"
)

func main() {
	fmt.Println("WOMS Http Interface start...")

	http.HandleFunc(Event_UP_Interface, eventmodule.EventUpload)
	http.HandleFunc(Query_Single_Event, eventmodule.QuerySingleEvent)
	http.HandleFunc(Event_Status_UpDate, eventmodule.EventStatusUpdate)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("Http Listen Error.", err.Error())
		log.Fatal("ListenAndServe error:", err)
	}
}
