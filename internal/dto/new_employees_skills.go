package dto

type NewEmployeesSkills struct {
	FullName       string `sheet:"Full Name"`
	Email          string `sheet:"Email"`
	DateOfBirth    string `sheet:"Date of Birth"`
	ContractType   string `sheet:"Contract Type"`
	WorkingForm    string `sheet:"Working Form"`
	OfficeTower    string `sheet:"Office Tower"`
	Position       string `sheet:"Position"`
	Level          string `sheet:"Level"`
	Project        string `sheet:"Project"`
	PM             string `sheet:"PM"`
	OnBoardingDate string `sheet:"On-Boarding Date"`
	DivisionLead   string `sheet:"Division Lead"`
	Division       string `sheet:"Division"`
	Skill          string `sheet:"Skill"`
	LinkCV         string `sheet:"Link CV"`
}
