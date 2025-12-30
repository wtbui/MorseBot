package lightsync

type LEffect struct {
	Id int `json:"id"`
	ParamId int `json:"paramId"`
}

var LColors = map[string]int{ 
	"red": 16711680,
	"blue": 255,
	"green": 65280,
	"purple": 16711935,
	"orange": 16753920,
	"white": 16777215,
	"yellow": 16776960,
}

var LTemps = map[string]int {
	"warm": 2000,
	"cool": 4000,
	"daylight": 9000,
}

var LEffects = map[string]LEffect{
	"valosignal": LEffect{3296, 3146}, 
}
