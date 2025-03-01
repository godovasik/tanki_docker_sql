package models

// по-хорошему мне нужно переделывать датастамп, чтоб там был id gear, но мне лень. буду так делать.

var gearIDs = map[string]int{
	// Hulls (от 0)
	"Wasp":     0,
	"Hornet":   1,
	"Hopper":   2,
	"Viking":   3,
	"Hunter":   4,
	"Crusader": 5,
	"Paladin":  6,
	"Dictator": 7,
	"Ares":     8,
	"Titan":    9,
	"Mammoth":  10,

	// Turrets (от 100)
	"Firebird": 100,
	"Freeze":   101,
	"Isida":    102,
	"Tesla":    103,
	"Hammer":   104,
	"Twins":    105,
	"Ricochet": 106,
	"Smoky":    107,
	"Striker":  108,
	"Vulcan":   109,
	"Thunder":  110,
	"Scorpion": 111,
	"Railgun":  112,
	"Magnum":   113,
	"Gauss":    114,
	"Shaft":    115,
}

func GetGearId(name string) (int, bool) {
	id, exists := gearIDs[name]
	return id, exists
}
