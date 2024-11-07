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

type UIPathGreetingJobInput struct {
	JobID int `json:"jobId"`
}
