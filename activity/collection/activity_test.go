package collection

import (
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&CollectionActivity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

// func getActivityMetadata() *activity.Metadata {

// 	if activityMetadata == nil {
// 		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
// 		if err != nil {
// 			panic("No Json Metadata found for activity.json path")
// 		}

// 		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
// 	}

// 	return activityMetadata
// }

func TestCreate(t *testing.T) {

	settings := &Settings{operation: "append"}
	iCtx := test.NewActivityInitContext(settings, nil)

	act, err := New(iCtx)
	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())

	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs

	act.Eval(tc)

	//check result attr
}
