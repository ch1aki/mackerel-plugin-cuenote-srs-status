package mpcuenotesrs

import (
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

func TestGraphDefinitionEnableGroup(t *testing.T) {
	cuenoteSrs := CuenoteSrsStatPlugin{EnableGroupStats: true}

	graphdef := cuenoteSrs.GraphDefinition()
	if len(graphdef) != 4 {
		t.Errorf("GetTempfilename: %d should be 16", len(graphdef))
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

func TestParseNowGroup(t *testing.T) {
	var cuenoteSrs CuenoteSrsStatPlugin
	stub := `group1	delivering	0
group1	undelivered	1111
group1	resend	2222
group1	success	111
group1	failure	222
group1	dnsfail	333
group1	exclusion	444
group1	bounce_unique	555
group1	canceled	11
group1	expired	22
group1	deferral	33
group1	dnsdeferral	44
group1	connfail	55
group1	bounce	66
group1	exception	0
group1	delivering	0
#bounce	delivering	0
#bounce	undelivered	19
#bounce	resend	1212
#bounce	success	0
#bounce	failure	0
#bounce	dnsfail	0
#bounce	exclusion	0
#bounce	bounce_unique	0
#bounce	canceled	0
#bounce	expired	0
#bounce	deferral	0
#bounce	dnsdeferral	0
#bounce	connfail	0
#bounce	bounce	0
#bounce	exception	0
#relay	delivering	0
#relay	undelivered	0
#relay	resend	0
#relay	success	0
#relay	failure	0
#relay	dnsfail	0
#relay	exclusion	0
#relay	bounce_unique	0
#relay	canceled	0
#relay	expired	0
#relay	deferral	0
#relay	dnsdeferral	0
#relay	connfail	0
#relay	bounce	0
#relay	exception	0
#default	delivering	0
#default	undelivered	0
#default	resend	0
#default	success	0
#default	failure	0
#default	dnsfail	0
#default	exclusion	0
#default	bounce_unique	0
#default	canceled	0
#default	expired	0
#default	deferral	0
#default	dnsdeferral	0
#default	connfail	0
#default	bounce	0
#default	exception	0
#all	delivering	806
#all	undelivered	0
#all	resend	0
#all	success	0
#all	failure	0
#all	dnsfail	0
#all	exclusion	0
#all	bounce_unique	0
#all	canceled	0
#all	expired	0
#all	deferral	0
#all	dnsdeferral	0
#all	connfail	0
#all	bounce	0
#all	exception	0
`

	cuenoteStat := strings.NewReader(stub)

	stat, err := cuenoteSrs.parseNowGroup(cuenoteStat)
	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["queue_group.undelivered.group1"]).String(), "float64")
	assert.EqualValues(t, stat["queue_group.undelivered.group1"], 1111)
	assert.EqualValues(t, reflect.TypeOf(stat["queue_group.resend.bounce"]).String(), "float64")
	assert.EqualValues(t, stat["queue_group.resend.bounce"], 1212)
}

func TestParseDeliveryGroup(t *testing.T) {
	var cuenoteSrs CuenoteSrsStatPlugin
	stub := `
group1	success	100
group1	failure	10
group1	dnsfail	0
group1	exclusion	0
group1	bounce_unique	0
group1	canceled	0
group1	expired	2
group1	deferral	58
group1	dnsdeferral	0
group1	connfail	0
group1	bounce	3
group1	exception	0
#all	success	100
#all	failure	0
#all	dnsfail	0
#all	exclusion	0
#all	bounce_unique	0
#all	canceled	0
#all	expired	0
#all	deferral	0
#all	dnsdeferral	950
#all	connfail	1459
#all	bounce	0
#all	exception	210
`

	cuenoteStat := strings.NewReader(stub)

	stat, err := cuenoteSrs.parseDeliveryGroup(cuenoteStat)
	assert.Nil(t, err)
	assert.EqualValues(t, reflect.TypeOf(stat["delivery_group.success.all"]).String(), "float64")
	assert.EqualValues(t, stat["delivery_group.success.all"], 100)
	assert.EqualValues(t, reflect.TypeOf(stat["delivery_group.failure.group1"]).String(), "float64")
	assert.EqualValues(t, stat["delivery_group.failure.group1"], 10)
}
