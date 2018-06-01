package vultr

import (
	"testing"
)

func TestSplitFirewallRule(t *testing.T) {
	cases := []struct {
		portRange string
		from      int
		to        int
		err       bool
	}{
		{
			portRange: "",
			from:      0,
			to:        0,
			err:       false,
		},
		{
			portRange: ":",
			from:      0,
			to:        0,
			err:       true,
		},
		{
			portRange: "-",
			from:      0,
			to:        0,
			err:       true,
		},
		{
			portRange: "22",
			from:      22,
			to:        22,
			err:       false,
		},
		{
			portRange: "foo",
			from:      0,
			to:        0,
			err:       true,
		},
		{
			portRange: "22:23",
			from:      0,
			to:        0,
			err:       true,
		},
		{
			portRange: "22 - 23",
			from:      22,
			to:        23,
			err:       false,
		},
		{
			portRange: "80-81",
			from:      80,
			to:        81,
			err:       false,
		},
	}

	for i, c := range cases {
		from, to, err := splitFirewallRule(c.portRange)
		if (err != nil) != c.err {
			no := "no"
			if c.err {
				no = "an"
			}
			t.Errorf("test case %d: expected %s error, got %v", i, no, err)
		}
		if from != c.from || to != c.to {
			t.Errorf("test case %d: expected range %d:%d, got %d:%d", i, c.from, c.to, from, to)
		}
	}
}
