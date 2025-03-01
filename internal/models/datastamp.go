package models // в этом файле у меня обьявление структуры которую я буду сохранять в дб,
// а так же методы для работы с ней типа вывода

import (
	// "fmt"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

var ( // я не помню где я их использую, возможно удалю позже. а может и нет.
	HULLS   = []string{"Wasp", "Hornet", "Hopper", "Viking", "Hunter", "Crusader", "Paladin", "Dictator", "Ares", "Titan", "Mammoth"}
	TURRETS = []string{"Firebird", "Freeze", "Isida", "Tesla", "Hammer", "Twins", "Ricochet", "Smoky", "Striker", "Vulcan", "Thunder", "Scorpion", "Railgun", "Magnum", "Gauss", "Shaft"}
	//DRONES  = []string{"Crisis", "Brutus", "Saboteur", "Trickster", "Mechanic", "Booster", "Defender", "Hyperion"}
)

// TODO - обьединить весь gear в одну мапу или слайс или массив структур или пошел я нахуй
type Datastamp struct {
	Timestamp      time.Time
	Name           string
	Rank           int
	Kills          int
	Deaths         int
	EarnedCrystals int
	GearScore      int
	Hulls          map[string]GearData
	Turrets        map[string]GearData
	Drones         map[string]GearData
	SuppliesUsed   map[string]int
}

type GearData struct {
	ScoreEarned   int
	SecondsPlayed int
}

type nameAndThing struct {
	Key   string
	Value GearData
}

func MapToSortedSlice(m map[string]GearData) []nameAndThing { // мапа с корпусами/пушками в список корпусов/пушек, сортированый
	// var ss []NameAndThing
	ss := make([]nameAndThing, 0, len(m)) // оптимизация хелл йеах
	for k, v := range m {
		ss = append(ss, nameAndThing{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value.SecondsPlayed > ss[j].Value.SecondsPlayed
	})
	return ss
}

func msToHours(microseconds int) int {
	return microseconds / (1000 * 60 * 60)
}

func msToSeconds(microseconds int) int {
	return microseconds / 1000
}

// дропается если ентитис больше чем 11, оно и не мудрено, да и мне пахую ес честно
func (d *Datastamp) NewPrint(howManyEntitiesPrint int) {

	//основное инфо

	briefTable := tablewriter.NewWriter(os.Stdout)
	briefTable.Append([]string{"Name", d.Name})
	briefTable.Append([]string{"K/D:", fmt.Sprintf("%.2f", (float32(d.Kills) / float32(d.Deaths)))})
	briefTable.Append([]string{"Score", fmt.Sprint(d.Rank)})
	briefTable.Render()

	// вывод таблицы пушек

	TurretTable := tablewriter.NewWriter(os.Stdout)
	TurretTable.SetHeader([]string{"Turret", "Score", "Hours played"})
	TurretSlice := MapToSortedSlice(d.Turrets)

	for _, v := range TurretSlice[:howManyEntitiesPrint] {
		TurretTable.Append([]string{v.Key, strconv.Itoa(v.Value.ScoreEarned), strconv.Itoa(msToHours(v.Value.SecondsPlayed))})
	}
	TurretTable.Render()

	// вывод таблицы корпусов

	HullTable := tablewriter.NewWriter(os.Stdout)
	HullTable.SetHeader([]string{"Hull", "Score", "Hours played"})
	HullSlice := MapToSortedSlice(d.Hulls)

	for _, v := range HullSlice[:howManyEntitiesPrint] {
		HullTable.Append([]string{v.Key, strconv.Itoa(v.Value.ScoreEarned), strconv.Itoa(msToHours(v.Value.SecondsPlayed))})
	}
	HullTable.Render()
}

func ConvertResponseToDatastamp(data *ResponseWrapper) *Datastamp {
	var d Datastamp
	r := data.Response
	d.Name, d.Rank, d.Kills, d.Deaths, d.EarnedCrystals, d.GearScore =
		r.Name, r.Rank, r.Kills, r.Deaths, r.EarnedCrystals, r.GearScore

	d.Hulls = make(map[string]GearData)
	d.Turrets = make(map[string]GearData)
	d.Drones = make(map[string]GearData)
	d.SuppliesUsed = make(map[string]int)

	for _, hull := range HULLS {
		d.Hulls[hull] = GearData{0, 0}
	}

	for _, a := range r.HullsPlayed {
		hull := d.Hulls[a.Name]
		hull.SecondsPlayed += msToSeconds(a.TimePlayed)
		hull.ScoreEarned += a.ScoreEarned
		d.Hulls[a.Name] = hull
	}

	for _, turret := range TURRETS {
		d.Turrets[turret] = GearData{0, 0}
	}
	for _, a := range r.TurretsPlayed {
		turret := d.Turrets[a.Name]
		turret.SecondsPlayed += msToSeconds(a.TimePlayed)
		turret.ScoreEarned += a.ScoreEarned
		d.Turrets[a.Name] = turret
	}

	for _, a := range r.DronesPlayed {
		d.Drones[a.Name] = GearData{a.ScoreEarned, msToSeconds(a.TimePlayed)}
	}

	for _, a := range r.SuppliesUsage {
		d.SuppliesUsed[a.Name] = a.Usages
	}

	d.Timestamp = time.Now()
	return &d

}
