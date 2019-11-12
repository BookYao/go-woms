package  eventmodule

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"../sqlmodule"
)

type eventQueryInfo struct {
	AuthUser string `json: "authUser"`
	AuthPasswd string  `json: "authPasswd"`
	TimeStamp string `json: "timeStamp"`
	EventID string `json: "eventID"`
}

func getEventDesc(id string, eventRes *EventQueryResponseInfo) bool {
	eid, _:= strconv.Atoi(id)
	sqlcmd := fmt.Sprintf("select ei_creatuser, ei_time, ei_level, ei_type, " +
		" \"ei_Subject\", ei_desc, ei_longitude, ei_latitude, ei_user, ei_status " +
		" from \"T_Event_Info\" where ei_id = '%d';", eid)

	fmt.Println("event query sql:", sqlcmd)
	db := sqlmodule.ConnectDB()
	defer db.Close()

	row, err := db.Query(sqlcmd)
	if err != nil {
		fmt.Println("Sqlcmd query failed!", err.Error())
		return false
	}

	for row.Next() {
		row.Scan(&eventRes.EventUpUser, &eventRes.EventTime, &eventRes.EventLevel,
			&eventRes.EventType, &eventRes.EventSubject,
			&eventRes.EventDesc, &eventRes.Longitude,
			&eventRes.Lantitude, &eventRes.AcceptUser, &eventRes.EventStatus)
	}
	return true
}

func getEventFile(id string, eventRes *EventQueryResponseInfo) bool {
	eid, _:= strconv.Atoi(id)
	sqlcmd := fmt.Sprintf("select ef_name, ef_dir, ef_size " +
		" from \"T_Event_File\" where ef_eid = '%d';", eid)

	db := sqlmodule.ConnectDB()
	defer db.Close()

	var fileinfo FileInfo
	row, err := db.Query(sqlcmd)
	if err != nil {
		fmt.Println("Sqlcmd query eventfile failed!", err.Error())
		return false
	}

	for row.Next() {
		row.Scan(&fileinfo.FileName, &fileinfo.FileDir, &fileinfo.FileSize)
		eventRes.FileList = append(eventRes.FileList, fileinfo)
	}
	return true
}

func getEventQueryRes(id string) *EventQueryResponseInfo {
	eventQueryRes := &EventQueryResponseInfo{}
	err := getEventDesc(id, eventQueryRes)
	if err != true {
		fmt.Println("getEventQueyRes faled")
		return nil
	}

	err = getEventFile(id, eventQueryRes)
	if err != true {
		fmt.Println("getEventQueyRes failed")
		return nil
	}
	return eventQueryRes
}

func getEventJsonByID(id string) []byte {
	EventQueryResponse := getEventQueryRes(id)

	EventQueryResponse.Result = "succeed"
	EventQueryResponse.ErrCode = "200"

	fmt.Println("eventQueryResponse: ", EventQueryResponse)
	jsonString, err := json.Marshal(&EventQueryResponse)
	if err != nil {
		fmt.Println("Marshal Event Query Json failed!")
		return nil
	}
	fmt.Println("jsonString: ", string(jsonString))

	return jsonString
}

func eventQueryFailedInfo() []byte {
	var EventQueryResponse EventQueryResponseInfo
	EventQueryResponse.Result = "failed"
	EventQueryResponse.ErrCode = "404"

	jsonString, err := json.Marshal(&EventQueryResponse)
	if err != nil  {
		fmt.Println("eventQueryFailedInfo build Failed, err:", err.Error())
		return nil
	}
	return jsonString
}

func QuerySingleEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}


/*	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("body: ", string(body))*/

	len := r.ContentLength
	body := make([]byte, len)

	defer r.Body.Close()
	r.Body.Read(body)

	fmt.Println("body:", string(body))
	eventQueryReqInfo := new(eventQueryInfo)
	json.Unmarshal(body, eventQueryReqInfo)

	resString := getEventJsonByID(eventQueryReqInfo.EventID)
	if resString == nil {
		fmt.Println("getEventJsonByID failed!")
		resString = eventQueryFailedInfo()
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(resString)
}
