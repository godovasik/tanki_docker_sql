package models

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"sort"
	"strconv"
	"time"
)

var (
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

func MapToSortedSlice(m map[string]Thing) []NameAndThing {
	var ss []NameAndThing
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

func (d *Datastamp) Print() {
	line := "------------------------------------------------\n"
	fmt.Print("Time:", d.Timestamp)
	fmt.Print("\nName: ", d.Name, "\nGS: ", d.GearScore)
	fmt.Print("\nKills: ", d.Kills, "\nDeaths: ", d.Deaths, "\nK/D:", float32(d.Kills)/float32(d.Deaths))

	fmt.Print("\nTurret\t\tScore\t\tTime played, h\n", line)
	turrets := MapToSortedSlice(d.Turrets)
	for _, a := range turrets[:5] {
		fmt.Print(a.Key, "\t\t", a.Value.ScoreEarned, "\t\t", msToHours(a.Value.TimePlayed), "\n")
	}

	hulls := MapToSortedSlice(d.Hulls)
	fmt.Print("\nHull\t\tScore\t\tTime played, h\n", line)
	for _, a := range hulls[:5] {
		fmt.Print(a.Key, "\t\t", a.Value.ScoreEarned, "\t\t", msToHours(a.Value.TimePlayed), "\n")
	}

}

func (d *Datastamp) NewPrint() {
	howManyEntitiesPrint := 5

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Turret", "Score", "Time played, h"})
	TurretSlice := MapToSortedSlice(d.Turrets)

	for _, v := range TurretSlice[:howManyEntitiesPrint] {
		table.Append([]string{v.Key, strconv.Itoa(v.Value.ScoreEarned), strconv.Itoa(v.Value.TimePlayed)})
	}

	table.Render()
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
