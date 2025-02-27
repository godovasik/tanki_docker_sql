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

type Datastamp struct {
	Timestamp      time.Time
	Name           string
	Rank           int
	Kills          int
	Deaths         int
	EarnedCrystals int
	GearScore      int
	Hulls          map[string]Thing
	Turrets        map[string]Thing
	Drones         map[string]Thing
	SuppliesUsed   map[string]int
}

type Thing struct {
	ScoreEarned int
	TimePlayed  int
}

type NameAndThing struct {
	Key   string
	Value Thing
}

func MapToSortedSlice(m map[string]Thing) []NameAndThing { // мапа с корпусами/пушками в список корпусов/пушек, сортированый
	// var ss []NameAndThing
	ss := make([]NameAndThing, 0, len(m)) // оптимизация хелл йеах
	for k, v := range m {
		ss = append(ss, NameAndThing{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value.TimePlayed > ss[j].Value.TimePlayed
	})
	return ss
}

func msToHours(microseconds int) int {
	return microseconds / (1000 * 60 * 60)
}

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
		TurretTable.Append([]string{v.Key, strconv.Itoa(v.Value.ScoreEarned), strconv.Itoa(msToHours(v.Value.TimePlayed))})
	}
	TurretTable.Render()

	// вывод таблицы корпусов

	HullTable := tablewriter.NewWriter(os.Stdout)
	HullTable.SetHeader([]string{"Hull", "Score", "Hours played"})
	HullSlice := MapToSortedSlice(d.Hulls)

	for _, v := range HullSlice[:howManyEntitiesPrint] {
		HullTable.Append([]string{v.Key, strconv.Itoa(v.Value.ScoreEarned), strconv.Itoa(msToHours(v.Value.TimePlayed))})
	}
	HullTable.Render()
}

func (d *Datastamp) ConvertResponseToDatastamp(data ResponseWrapper) {
	r := data.Response
	d.Name, d.Rank, d.Kills, d.Deaths, d.EarnedCrystals, d.GearScore =
		r.Name, r.Rank, r.Kills, r.Deaths, r.EarnedCrystals, r.GearScore

	d.Hulls = make(map[string]Thing)
	d.Turrets = make(map[string]Thing)
	d.Drones = make(map[string]Thing)
	d.SuppliesUsed = make(map[string]int)

	for _, hull := range HULLS {
		d.Hulls[hull] = Thing{0, 0}
	}

	for _, a := range r.HullsPlayed {
		hull := d.Hulls[a.Name]
		hull.TimePlayed += a.TimePlayed
		hull.ScoreEarned += a.ScoreEarned
		d.Hulls[a.Name] = hull
	}

	for _, turret := range TURRETS {
		d.Turrets[turret] = Thing{0, 0}
	}
	for _, a := range r.TurretsPlayed {
		turret := d.Turrets[a.Name]
		turret.TimePlayed += a.TimePlayed
		turret.ScoreEarned += a.ScoreEarned
		d.Turrets[a.Name] = turret
	}

	for _, a := range r.DronesPlayed {
		d.Drones[a.Name] = Thing{a.ScoreEarned, a.TimePlayed}
	}

	for _, a := range r.SuppliesUsage {
		d.SuppliesUsed[a.Name] = a.Usages
	}

	d.Timestamp = time.Now().Truncate(24 * time.Hour)

}
