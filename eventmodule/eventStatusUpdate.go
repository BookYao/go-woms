package eventmodule

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"../sqlmodule"
	"../mqmodule"
)

type EventStatusUpdateReqInfo struct  {
	AuthUser     string     `json: "AuthUser"`
	AuthPasswd   string     `json: "AuthPasswd"`
	Timestamp    string     `json: "Timestamp"`
	EventID   string `json: "EventID"`
	EventStatus int `json: "EventStatus"`
}

type EventStatusMqInfo struct {
	ServiceCode string `json: "ServiceCode"`
	DataType string `json: "DataType"`
	EventID   string `json: "EventID"`
	EventStatus int `json: "EventStatus"`
}

func eventStatusFailedInfo() []byte {
	var  resultInfo ResultInfo
	resultInfo.Result = "failed"
	resultInfo.ErrCode = "404"

	jsonString, err := json.Marshal(&resultInfo)
	if err != nil  {
		fmt.Println("eventQueryFailedInfo build Failed, err:", err.Error())
		return nil
	}
	return jsonString
}

func EventStatusSuccessInfo() []byte {
	var  resultInfo ResultInfo
	resultInfo.Result = "succeed"
	resultInfo.ErrCode = "200"

	jsonString, err := json.Marshal(&resultInfo)
	if err != nil  {
		fmt.Println("eventQueryFailedInfo build Failed, err:", err.Error())
		return nil
	}
	return jsonString
}

func updateEventStatus(eventUpdateInfo *EventStatusUpdateReqInfo) bool {
	eid, _:= strconv.Atoi(eventUpdateInfo.EventID)
	sqlcmd := fmt.Sprintf("Update \"T_Event_Info\" set ei_status = $1 where ei_id = $2;")

	conn := sqlmodule.ConnectDB()
	defer conn.Close()

	smtp, err := conn.Prepare(sqlcmd)
	if err != nil {
		fmt.Println("Update Event Status Prepare failed.", err.Error())
		return false
	}

	_, err = smtp.Exec(eventUpdateInfo.EventStatus, eid)
	if err != nil {
		fmt.Println("Update Event Status Exec failed!", err.Error())
		return false
	}
	return true
}

func getMqNotifyUser(eventUpdateInfo *EventStatusUpdateReqInfo) string {
	var user string
	eid, _:= strconv.Atoi(eventUpdateInfo.EventID)
	sqlcmd := fmt.Sprintf("select ei_creatuser from \"T_Event_Info\" where ei_id = %d;", eid)

	conn := sqlmodule.ConnectDB()
	if conn == nil {
		fmt.Println("getMqNotifyUser DB connect Failed.")
		return ""
	}

	fmt.Println("sql:", sqlcmd)
	defer conn.Close()
	row := conn.QueryRow(sqlcmd)

	row.Scan(&user)
	fmt.Println("Notify User:", user)
	return user
}

func eventStatusNotifyMq(eventUpdateInfo *EventStatusUpdateReqInfo) bool {
	var eventStatusMqInfo EventStatusMqInfo
	eventStatusMqInfo.EventID = eventUpdateInfo.EventID
	eventStatusMqInfo.EventStatus = eventUpdateInfo.EventStatus
	eventStatusMqInfo.ServiceCode = "woms"
	eventStatusMqInfo.DataType = "2"

	jsonString, err := json.Marshal(&eventStatusMqInfo)
	if err != nil {
		fmt.Println("Build Event Status MQ Info Failed.")
		return false
	}
	fmt.Println("Event Update MQ Json:", string(jsonString))

	notifyUser := getMqNotifyUser(eventUpdateInfo)
	if notifyUser == "" {
		fmt.Println("getMqNotifyUser is Nil")
		return false
	}
	mqmodule.MqPublish(notifyUser, jsonString)
	return true
}

func EventStatusUpdate(w http.ResponseWriter, r *http.Request) {
	var retString []byte
	retString = EventStatusSuccessInfo()

	w.Header().Set("Content-Type", "application/json")
	if r.Method != "POST" {
		fmt.Println("EventStatusUpdate Request Method Error!")
		retString = eventStatusFailedInfo()
		w.Write(retString)
		return
	}
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)

	eventStatusUpdateReqInfo := new(EventStatusUpdateReqInfo)
	json.Unmarshal(body, eventStatusUpdateReqInfo)

	status := updateEventStatus(eventStatusUpdateReqInfo)
	if status != true {
		fmt.Println("updateEventStatus failed.")
		retString = eventStatusFailedInfo()
		w.Write(retString)
		return
	}

	status = eventStatusNotifyMq(eventStatusUpdateReqInfo)
	if status != true {
		fmt.Println("Event Status Notify MQ failed.")
		retString = eventQueryFailedInfo()
		w.Write(retString)
		return
	}
	retString = EventStatusSuccessInfo()
	w.Write(retString)
	return
}