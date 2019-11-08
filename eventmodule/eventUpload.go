package eventmodule

import (
	_ "../../github.com/lib/pq"
	"../mqmodule"
	"../sqlmodule"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type EventResult struct {
	Result    string
	ErrorCode string
	EventID   int
}

type EventInfo struct {
	AuthUser     string     `json: "AuthUser"`
	AuthPasswd   string     `json: "AuthPasswd"`
	Timestamp    string     `json: "Timestamp"`
	EventTime    string     `json: "EventTime"`
	EventLevel   int        `json: "EventLevel"`
	EventType    int        `json: "EventType"`
	EventSubject string     `json: "EventSubject"`
	EventDesc    string     `json: "EventDesc"`
	Longitude    float64    `json: "Longitude"`
	Lantitude    float64    `json: "Lantitude"`
	AcceptUser   string     `json: "AcceptUser"`
	EventStatus  string     `json: "EventStatus, omitempty"`
	FileList     []FileInfo `json: "FileList"`
}

type MqEventUpInfo struct {
	ServiceCode  string  `json: ServiceCode`
	DataType     string  `json: DataType`
	EventID      string  `json: "EventID"`
	EventTime    string  `json: "EventTime"`
	EventLevel   int     `json: "EventLevel"`
	EventType    int     `json: "EventType"`
	EventSubject string  `json: "EventSubject"`
	Longitude    float64 `json: "Longitude"`
	Lantitude    float64 `json: "Lantitude"`
}

func failInfo() *EventResult {
	return &EventResult{"failed", "404", -1}
}

func successInfo(eid int) *EventResult {
	return &EventResult{"success", "200", eid}
}

func upEventMqInfo(eventid int, eventinfo *EventInfo) []byte {
	eventId := fmt.Sprintf("%d", eventid)
	mqInfo := &MqEventUpInfo{
		ServiceCode: "woms", DataType: "1", EventID: eventId,
		EventTime: eventinfo.EventTime, EventLevel: eventinfo.EventLevel,
		EventType: eventinfo.EventType, EventSubject: eventinfo.EventSubject,
		Longitude: eventinfo.Longitude, Lantitude: eventinfo.Lantitude}

	eventMqInfo, err := json.Marshal(mqInfo)
	if err != nil {
		fmt.Println("eventMqInfo Json Msg is error!", err.Error())
		return []byte("")
	}
	return eventMqInfo
}
func getEventID(db *sql.DB, user string, timestamp string) int {
	sqlcmd := fmt.Sprintf("select ei_id from \"T_Event_Info\" "+
		"where ei_creatuser = '%s' and ei_time = '%s';", user, timestamp)

	fmt.Println("getEventID sql:", sqlcmd)
	row, err := db.Query(sqlcmd)
	if err != nil {
		fmt.Println("db Query failed.")
		return -1
	}

	var ei_id int = -1
	for row.Next() {
		row.Scan(&ei_id)
	}

	return ei_id
}

func sqlSaveEventInfo(eventinfo *EventInfo) int {
	db := sqlmodule.ConnectDB()
	defer db.Close()

	sqlcmd := fmt.Sprintf(
		"insert into \"T_Event_Info\" " +
			"(ei_creatuser, ei_time, ei_level, ei_type, \"ei_Subject\", ei_desc, " +
			"ei_longitude, ei_latitude, ei_user, ei_status) " +
			"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)")

	fmt.Println("Insert sqlcmd: ", sqlcmd)
	smtp, err := db.Prepare(sqlcmd)
	if err != nil {
		fmt.Println("Sqlcmd Insert Prepare failed!", err.Error())
		return -1
	}

	_, err = smtp.Exec(eventinfo.AuthUser, eventinfo.EventTime,
		eventinfo.EventLevel, eventinfo.EventType, eventinfo.EventSubject, eventinfo.EventDesc,
		eventinfo.Longitude, eventinfo.Lantitude, eventinfo.AuthUser, 1)
	if err != nil {
		fmt.Println("Sqlcmd Prepare failed!", err.Error())
		return -1
	}

	return getEventID(db, eventinfo.AuthUser, eventinfo.EventTime)
}

func sqlSaveEventFile(ei_id int, ef_dir string, fileinfo []FileInfo) bool {
	db := sqlmodule.ConnectDB()
	defer db.Close()

	sqlcmd := fmt.Sprintf("insert into \"T_Event_File\" (ef_eid, ef_name, ef_size, ef_dir)" +
		"values ($1, $2, $3, $4)")

	fmt.Println("Insert File sqlcmd: ", sqlcmd)
	smtp, err := db.Prepare(sqlcmd)
	if err != nil {
		fmt.Println("Sqlcmd Insert Prepare failed!", err.Error())
		return false
	}

	for i := 0; i < len(fileinfo); i++ {
		_, err = smtp.Exec(ei_id, fileinfo[i].FileName, fileinfo[i].FileSize, ef_dir)
		if err != nil {
			fmt.Println("Sqlcmd Prepare failed!", err.Error())
			return false
		}
	}

	return true
}

func dirIsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, nil
}

func EventUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Event up interface...")

	/*	fmt.Println("Method: ", r.Method)
		fmt.Println("Url: ", r.RequestURI)
		fmt.Println("URL host: ", r.URL.Host)
		fmt.Println("URL Opaque: ", r.URL.Opaque)
		fmt.Println("URL Path: ", r.URL.Path)
		fmt.Println("URL Fragment: ", r.URL.Fragment)
		fmt.Println("URL User String: ", r.URL.User.String())

		fmt.Println("HTTP Proto: ", r.Proto)
		fmt.Println("HTTP Proto Major: ", r.ProtoMajor)
		fmt.Println("HTTP Proto Minor: ", r.ProtoMinor)

		for headKey, headVal := range r.Header {
			for _, headContent := range headVal {
				fmt.Println("HTTP Head:", headKey, ":", headContent)
			}
		}*/

	isMultipart := false
	for _, contentType := range r.Header["Content-Type"] {
		if strings.Index(contentType, "multipart/form-data") != -1 {
			isMultipart = true
		}
	}

	if isMultipart == true {
		r.ParseMultipartForm(128)
		fmt.Println("Parse Multipart Form")
	} else {
		r.ParseForm()
		fmt.Println("Parse Form")
	}

	/*	fmt.Println("ContentLen: ", r.ContentLength)
		fmt.Println("Close: ", r.Close)
		fmt.Println("Host: ", r.Host)*/

	if isMultipart == true {
		//fmt.Println("MultiPartForm: ", r.MultipartForm)
		fmt.Println("============================")

		text := r.MultipartForm.Value
		fmt.Println("INFO: ", text["info"][0])

		eventInfo := EventInfo{}
		json.Unmarshal([]byte(text["info"][0]), &eventInfo)

		fmt.Println("eventInfo: ", eventInfo)
		fmt.Println("AuthUser:", eventInfo.AuthUser)
		fmt.Println("fileInfo:", eventInfo.FileList)

		ei_id := sqlSaveEventInfo(&eventInfo)
		fmt.Println("Save event Info Success. ei_id:", ei_id)

		intEventTime, err := strconv.Atoi(eventInfo.EventTime)
		if err != nil {
			fmt.Println("String event time atoi error:", err.Error())
		}
		fmt.Println("int Event time :", intEventTime)
		timedate := time.Unix(int64(intEventTime), 0)

		ef_dir := "/var/woms/" + timedate.Format("2006-01-02")
		fmt.Println("ef_dir: ", ef_dir)
		sqlSaveEventFile(ei_id, ef_dir, eventInfo.FileList)

		dirStat, _ := dirIsExist(ef_dir)
		if !dirStat {
			os.Mkdir(ef_dir, 0777)
		}

		files := r.MultipartForm.File
		for fileKey, fileVal := range files {
			fmt.Println("fileKey: ", fileKey)
			for _, fileName := range fileVal {
				fmt.Println("filename: ", fileName.Filename, "fileSize:", fileName.Size)
				srcfile, _, err := r.FormFile(fileKey)
				if err != nil {
					fmt.Printf("FormFile failed. err:%s", err.Error())
					continue
				}
				defer srcfile.Close()
				filePath := ef_dir + "/" + fileName.Filename
				picFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0600)
				if err != nil {
					fmt.Printf("open file failed. err:%s-file:%s", err.Error(), filePath)
					continue
				}

				defer picFile.Close()
				io.Copy(picFile, srcfile)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		retInfo, err := json.Marshal(successInfo(ei_id))
		if err != nil {
			fmt.Println("success json build failed.")
		}
		w.Write(retInfo)

		mqInfo := upEventMqInfo(ei_id, &eventInfo)
		mqmodule.MqPublish(eventInfo.AcceptUser, mqInfo)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	retInfo, err := json.Marshal(failInfo())
	if err != nil {
		fmt.Println("fail json build failed.")
	}
	w.Write(retInfo)
	return
}
