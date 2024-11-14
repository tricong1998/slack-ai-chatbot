package dto

type UIPathGreetingNewEmployee struct {
	SkillFile     string `json:"SkillFile"`
	PersonalEmail string `json:"PersonalEmail"`
}

type UIPathTriggerResponse struct {
	Key          string `json:"key"`
	State        string `json:"state"`
	CreationTime string `json:"creationTime"`
	ID           int    `json:"id"`
}

type UIPathJobDetails struct {
	ODataContext                       string `json:"@odata.context"`
	Key                                string `json:"Key"`
	StartTime                          string `json:"StartTime"`
	EndTime                            string `json:"EndTime"`
	State                              string `json:"State"`
	JobPriority                        string `json:"JobPriority"`
	SpecificPriorityValue              int    `json:"SpecificPriorityValue"`
	ResourceOverwrites                 string `json:"ResourceOverwrites"`
	Source                             string `json:"Source"`
	SourceType                         string `json:"SourceType"`
	BatchExecutionKey                  string `json:"BatchExecutionKey"`
	Info                               string `json:"Info"`
	CreationTime                       string `json:"CreationTime"`
	StartingScheduleId                 string `json:"StartingScheduleId"`
	ReleaseName                        string `json:"ReleaseName"`
	Type                               string `json:"Type"`
	InputArguments                     string `json:"InputArguments"`
	OutputArguments                    string `json:"OutputArguments"`
	HostMachineName                    string `json:"HostMachineName"`
	HasMediaRecorded                   bool   `json:"HasMediaRecorded"`
	HasVideoRecorded                   bool   `json:"HasVideoRecorded"`
	PersistenceId                      string `json:"PersistenceId"`
	ResumeVersion                      string `json:"ResumeVersion"`
	StopStrategy                       string `json:"StopStrategy"`
	RuntimeType                        string `json:"RuntimeType"`
	RequiresUserInteraction            bool   `json:"RequiresUserInteraction"`
	ReleaseVersionId                   int    `json:"ReleaseVersionId"`
	EntryPointPath                     string `json:"EntryPointPath"`
	OrganizationUnitId                 int    `json:"OrganizationUnitId"`
	OrganizationUnitFullyQualifiedName string `json:"OrganizationUnitFullyQualifiedName"`
	Reference                          string `json:"Reference"`
	ProcessType                        string `json:"ProcessType"`
	ProfilingOptions                   string `json:"ProfilingOptions"`
	ResumeOnSameContext                bool   `json:"ResumeOnSameContext"`
	LocalSystemAccount                 string `json:"LocalSystemAccount"`
	OrchestratorUserIdentity           string `json:"OrchestratorUserIdentity"`
	RemoteControlAccess                string `json:"RemoteControlAccess"`
	StartingTriggerId                  string `json:"StartingTriggerId"`
	MaxExpectedRunningTimeSeconds      string `json:"MaxExpectedRunningTimeSeconds"`
	ServerlessJobType                  string `json:"ServerlessJobType"`
	ResumeTime                         string `json:"ResumeTime"`
	LastModificationTime               string `json:"LastModificationTime"`
	ProjectKey                         string `json:"ProjectKey"`
	ParentOperationId                  string `json:"ParentOperationId"`
	EnableAutopilotHealing             bool   `json:"EnableAutopilotHealing"`
	Id                                 int    `json:"Id"`
	AutopilotForRobots                 string `json:"AutopilotForRobots"`
}

type UIPathGreetingOutput struct {
	Position string `json:"Position"`
	Skill    string `json:"Skill"`
	Division string `json:"Division"`
	FullName string `json:"FullName"`
	Greeting string `json:"Greeting"`
}

type UIPathCheckingJobInput struct {
	JobID int `json:"jobId"`
}

type UIPathFillBuddyInput struct {
	InputSheet  string `json:"inputSheet"`
	OutputSheet string `json:"outputSheet"`
}

type UIPathFillBuddyOutput struct {
	BuddyFormName string `json:"buddyFormName"`
}

type UIPathCreateLeaveRequestInput struct {
	RequestDateFrom string `json:"request_date_from"`
	RequestDateTo   string `json:"request_date_to"`
	Description     string `json:"description"`
	CalendarId      int    `json:"calendar_id"`
	HolidayStatusId int    `json:"holiday_status_id"`
	HourFrom        int    `json:"hour_from"`
	HourTo          int    `json:"hour_to"`
	WorkEmail       string `json:"work_email"`
}

type UIPathCreateIntegrateTrainingInput struct {
	SheetURL  string `json:"sheetURL"`
	SheetName string `json:"sheetName"`
}

type UIPathCreateIntegrateTrainingOutput struct {
	ErrMessage string `json:"errMessage"`
	CalendarId string `json:"calendarId"`
}

type UIPathErrorTriggerJob struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"errorCode"`
	TraceId   string `json:"traceId"`
}

type UIPathLeaveCode struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

var AppMappingCodeLeave = []UIPathLeaveCode{
	{Code: 35, Name: "Unpaid Leave"},
	{Code: 39, Name: "Remote work"},
	{Code: 44, Name: "Maternity leave"},
	{Code: 45, Name: "Family business (Applicable to self/child/father/mother)"},
	{Code: 56, Name: "Onsite"},
	{Code: 57, Name: "Absence to deal with company affairs"},
	{Code: 59, Name: "Business travel"},
	{Code: 60, Name: "Forgot timekeeping"},
	{Code: 65, Name: "Onboard (0 remaining out of 0 days)"},
}

type UIPathWorkingTime struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

var AppMappingCodeWorkingTime = []UIPathWorkingTime{
	{Code: 35, Name: "8:00 - 17:30"},
	{Code: 36, Name: "8:30 - 18:00"},
	{Code: 37, Name: "9:00 - 18:30"},
}

type UIPathLeaveOutput struct {
	Response string `json:"response"`
}

type UIPathLeaveOutputResponse struct {
	Jsonrpc string                   `json:"jsonrpc"`
	ID      string                   `json:"id"`
	Result  *UIPathLeaveOutputResult `json:"result"`
	Error   *UIPathLeaveOutputError  `json:"error"`
}

type UIPathLeaveOutputResult struct {
	Code              int    `json:"code"`
	ResID             int    `json:"res_id"`
	CalendarID        int    `json:"calendar_id"`
	EmployeeName      string `json:"employee_name"`
	HolidayStatusID   int    `json:"holiday_status_id"`
	HolidayStatusName string `json:"holiday_status_name"`
	Period            string `json:"period"`
	RequestDateFrom   string `json:"request_date_from"`
	RequestDateTo     string `json:"request_date_to"`
	DateFrom          string `json:"date_from"`
	DateTo            string `json:"date_to"`
	HourFrom          int    `json:"hour_from"`
	HourTo            int    `json:"hour_to"`
	Status            string `json:"status"`
	Duration          string `json:"duration"`
	Approver          string `json:"approver"`
	Description       string `json:"description"`
	RefuseReason      string `json:"refuse_reason"`
	Attachment        bool   `json:"attachment"`
}

type UIPathLeaveOutputError struct {
	Message string                     `json:"message"`
	Code    int                        `json:"code"`
	Data    UIPathLeaveOutputErrorData `json:"data"`
}

type UIPathLeaveOutputErrorData struct {
	Name          string   `json:"name"`
	Debug         string   `json:"debug"`
	Message       string   `json:"message"`
	Arguments     []string `json:"arguments"`
	ExceptionType string   `json:"exception_type"`
}
