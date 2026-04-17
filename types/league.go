package types

import "fmt"

type League struct {
	Fragment int ``
}

type LeagueParts struct {
	Name    int
	Rank    int
	Tooltip string
}

func (l *League) Split() *LeagueParts {

	name := l.Fragment / 10
	rank := l.Fragment % 10
	var text string

	switch name {
	case 1:
		text = "Bronze"
	case 2:
		text = "Silver"
	case 3:
		text = "Gold"
	case 4:
		text = "Crystal"
	case 5:
		text = "Master"
	case 6:
		text = "Champion"
	}

	return &LeagueParts{
		Name:    name,
		Rank:    rank,
		Tooltip: fmt.Sprintf("%s %d", text, rank),
	}
}