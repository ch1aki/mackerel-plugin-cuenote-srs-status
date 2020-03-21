package mpcuenotesrs

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var cuenoteSrs CuenoteSrsStatPlugin

	graphdef := cuenoteSrs.GraphDefinition()
	if len(graphdef) != 1 {
		t.Errorf("GetTempfilename: %d should be 1", len(graphdef))
	}
}

func TestParseNowTotal(t *testing.T) {
	var cuenoteSrs CuenoteSrsStatPlugin
	stub := `delivering	803
undelivered	4680
resend	2869
success	0
failure	0
dnsfail	0
exclusion	0
bounce_unique	0
canceled	0
expired	0
deferral	0
dnsdeferral	0
connfail	0
bounce	0
exception	0
`

	cuenoteStat := strings.NewReader(stub)

	stat, err := cuenoteSrs.parseNowTotal(cuenoteStat)
	fmt.Println(stat)
	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["delivering"]).String(), "float64")
	assert.EqualValues(t, stat["delivering"], 803)
	assert.EqualValues(t, reflect.TypeOf(stat["undelivered"]).String(), "float64")
	assert.EqualValues(t, stat["undelivered"], 4680)
	assert.EqualValues(t, reflect.TypeOf(stat["undelivered"]).String(), "float64")
	assert.EqualValues(t, stat["undelivered"], 4680)
	assert.EqualValues(t, reflect.TypeOf(stat["resend"]).String(), "float64")
	assert.EqualValues(t, stat["resend"], 2869)
}
