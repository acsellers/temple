package assets

import (
	"errors"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	NotFound error
	configs  map[string]*Config
	hosts    map[string]string
	master   *sync.RWMutex
)

func init() {
	NotFound = errors.New("Item missing")
	configs = make(map[string]*Config)
	hosts = make(map[string]string)
	master = &sync.RWMutex{}
}

type Config struct {
	Name         string
	Hosts        map[string]*Host
	Lock         *sync.RWMutex
	AssetFolders map[string]*AssetFolder
}

type AssetFolder struct {
	Assets  map[string][]byte
	ETags   map[string]string
	Headers map[string]string
	Expires *BigDuration
}

type BigDuration struct {
	Hours  int
	Days   int
	Months int
	Years  int
}

// marvel at the basicness of the Andrew regex, IDK lol
var bdr = regexp.MustCompile(`(\d+)([ymdh])(\d+)?([ymdh])?(\d+)?([ymdh])?(\d+)?([ymdh])?`)

func NewBigDuration(s string) (*BigDuration, error) {
	args := bdr.FindStringSubmatch(s)
	if len(args) < 2 {
		return nil, NotFound
	}

	bd := &BigDuration{}
	args = args[1:]
	for len(args) >= 2 && args[0] != "" {
		i, _ := strconv.Atoi(args[0])
		if args[1] == "y" {
			bd.Years = i
		}
		if args[1] == "m" {
			bd.Months = i
		}
		if args[1] == "d" {
			bd.Days = i
		}
		if args[1] == "h" {
			bd.Hours = i
		}
		args = args[2:]
	}

	return bd, nil
}

func (bd BigDuration) FromNow() time.Time {
	return time.Now().
		AddDate(bd.Years, bd.Months, bd.Days).
		Add(time.Duration(bd.Hours) * time.Hour)
}

type Host struct {
	Name      string
	Overrides map[string]*AssetFolder
}
