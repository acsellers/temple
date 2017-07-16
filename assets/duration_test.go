package assets

import "testing"

func TestBigDuration(t *testing.T) {
	bd, err := NewBigDuration("12y")
	if err != nil {
		t.Fatal("Can't parse years")
	}
	if bd.Years != 12 {
		t.Fatal("Can't parse years")
	}
	bd, err = NewBigDuration("12m")
	if err != nil {
		t.Fatal("Can't parse months")
	}
	if bd.Months != 12 {
		t.Fatal("Can't parse months")
	}
	bd, err = NewBigDuration("12d")
	if err != nil {
		t.Fatal("Can't parse days")
	}
	if bd.Days != 12 {
		t.Fatal("Can't parse days")
	}
	bd, err = NewBigDuration("12h")
	if err != nil {
		t.Fatal("Can't parse hours")
	}
	if bd.Hours != 12 {
		t.Fatal("Can't parse hours")
	}
	bd, err = NewBigDuration("1y5h")
	if err != nil {
		t.Fatal("Can't parse years & hours")
	}
	if bd.Years != 1 {
		t.Fatal("Can't parse years")
	}
	if bd.Hours != 5 {
		t.Fatal("Can't parse hours")
	}
	bd, err = NewBigDuration("1h5y")
	if err != nil {
		t.Fatal("Can't parse years & hours")
	}
	if bd.Hours != 1 {
		t.Fatal("Can't parse hours")
	}
	if bd.Years != 5 {
		t.Fatal("Can't parse years")
	}
}
