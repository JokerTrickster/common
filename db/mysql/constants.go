package mysql

const (
	// 시나리오 상수
	ScenarioAll    = 0
	ScenarioCouple = iota
	ScenarioSolo
	ScenarioFamily
	ScenarioCompany
	ScenarioFriend
)

var ScenarioMap = map[string]int{
	"연인": ScenarioCouple,
	"혼밥": ScenarioSolo,
	"가족": ScenarioFamily,
	"회식": ScenarioCompany,
	"친구": ScenarioFriend,
}
