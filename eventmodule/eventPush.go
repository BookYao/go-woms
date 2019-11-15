package eventmodule

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type EventUpPushInfo struct {
	EventID      string     `json: "EventID"`
	EventTime    string     `json: "EventTime"`
	EventLevel   int        `json: "EventLevel"`
	EventType    int        `json: "EventType"`
	EventSubject string     `json: "EventSubject"`
	EventDesc    string     `json: "EventDesc"`
	Longitude    float64    `json: "Longitude"`
	Lantitude    float64    `json: "Lantitude"`
	EventUpUser  string     `json: "EventUpUser"`
	AcceptUser   string     `json: "AcceptUser"`
	EventStatus  string     `json: "EventStatus, omitempty"`
	FileList     []FileInfo `json: "FileList"`
}

const (
	EVENT_PUSH_FILE = "/usr/local/woms/configs/platform.ini"
)
func getPushUrl() string {
	fp, err := os.Open(EVENT_PUSH_FILE)
	if err != nil {
		fmt.Println("Open Event Push File failed!")
		return ""
	}
	defer fp.Close()

	var url string
	content := bufio.NewReader(fp)
	for {
		line, err := content.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Read Event Push File Failed.", err.Error())
				break
			}
		}

		status := strings.Contains(line, "url=")
		if status == false {
			continue
		} else {
			url = line[len("url="):]
			strings.TrimSpace(url)
			break
		}
	}
	return url
}

func eventUpPushOtherPlat(eventUpPushInfo *EventUpPushInfo) bool {
	jsonString, err := json.Marshal(eventUpPushInfo)
	if err != nil {
		fmt.Println("Event Up Push Info Marsha Failed.")
		return false
	}
	fmt.Println("Event Push Info:", string(jsonString))

	url := getPushUrl()
	if url == "" {
		fmt.Println("Not configure Other Platform Push Url...")
		return false
	}
	fmt.Printf("Event Push Url:%s\n", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonString))
	if err != nil {
		fmt.Println("Http New Request Failed")
		return false
	}
	req.Header.Set("Content-type", "Application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Event UP Push Client DO Failed.")
		return false
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Event Up Push resp:", string(respBody))
	return true
}

func eventUpPush(eventInfo *EventInfo, filedir string, eid int, eventStatus int) bool {
	eventUpPushInfo := &EventUpPushInfo{}

	eventUpPushInfo.EventID = fmt.Sprintf("%d", eid)
	eventUpPushInfo.EventStatus = fmt.Sprintf("%d", eventStatus)
	eventUpPushInfo.AcceptUser = eventInfo.AcceptUser
	eventUpPushInfo.EventType = eventInfo.EventType
	eventUpPushInfo.EventLevel = eventInfo.EventLevel
	eventUpPushInfo.EventTime = eventInfo.EventTime
	eventUpPushInfo.EventSubject = eventInfo.EventSubject
	eventUpPushInfo.EventDesc = eventInfo.EventDesc
	eventUpPushInfo.Longitude = eventInfo.Longitude
	eventUpPushInfo.Lantitude = eventInfo.Lantitude
	eventUpPushInfo.EventUpUser = eventInfo.AuthUser
	eventUpPushInfo.FileList = eventInfo.FileList
	for i := 0; i < len(eventInfo.FileList); i++ {
		eventUpPushInfo.FileList[i].FileDir = filedir
	}

	status := eventUpPushOtherPlat(eventUpPushInfo)
	if status != true {
		fmt.Println("eventUpPushOtherPlat Failed")
		return false
	}
	return true
}
