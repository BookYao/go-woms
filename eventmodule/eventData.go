package eventmodule

const (
	Event_Upload_Status = 1
	Event_Misreport_Status = 2
	Event_Finish_Status = 3
)

type EventDesc struct {
	EventTime    string  `json: "EventTime"`
	EventLevel   int     `json: "EventLevel"`
	EventType    int     `json: "EventType"`
	EventSubject string  `json: "EventSubject"`
	EventDesc    string  `json: "EventDesc"`
	Longitude    float64 `json: "Longitude"`
	Lantitude    float64 `json: "Lantitude"`
	EventStatus  string  `json: "EventStatus, omitempty"`
	EventId      string  `json: "EventID, omitempty"`
}

type FileInfo struct {
	FileName string `json: "FileName"`
	FileSize int64  `json: "FileSize"`
	FileDir  string `json: "FileDir, omitempty"`
}

type EventFile struct {
	FileList []FileInfo `json: "FileList"`
}

type EventReporter struct {
	AuthUser   string `json: "AuthUser"`
	AuthPasswd string `json: "AuthPasswd"`
	Timestamp  string `json: "Timestamp"`
}

type ResultInfo struct {
	Result  string `json: "Result"`
	ErrCode string `json: "ErrCode"`
}

type EventAccept struct {
	AcceptUser string `json: "AcceptUser"`
}

type EventQueryResponseInfo struct {
	Result       string     `json: "Result"`
	ErrCode      string     `json: "ErrCode"`
	EventUpUser  string     `json: "EventUpUser"`
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
