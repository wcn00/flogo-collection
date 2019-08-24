package collection

import (
	"encoding/json"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/mapper"
	logger "github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

type initContext struct {
	settings map[string]interface{}
}

func newInitContext(values map[string]interface{}) *initContext {
	if values == nil {
		values = make(map[string]interface{})
	}
	return &initContext{
		settings: values,
	}
}

func (i *initContext) Settings() map[string]interface{} {
	return i.settings
}

func (i *initContext) MapperFactory() mapper.Factory {
	return nil
}

func (i *initContext) Logger() logger.Logger {
	return logger.RootLogger()
}

func TestRegister(t *testing.T) {
	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)
	assert.NotNil(t, act)
}

func getObj() interface{} {
	objectJSON := []byte(`{
		"obj":{
		"name":"walter",
		"age" : 45 }
		}`)
	obj := make(map[string]interface{})
	json.Unmarshal(objectJSON, &obj)
	return obj

}
func TestEvalAppendNokey(t *testing.T) {
	act, err := New(newInitContext(map[string]interface{}{
		"operation": "append",
	}))
	if err != nil {
		t.Errorf("Could not create initial context:  %s", err)
		t.Fail()
	}

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("object", getObj())
	_, err = act.Eval(tc)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("Could not execute activty:  %s", err)
		t.Fail()
	}
	key := tc.GetOutput("key").(string)
	size := tc.GetOutput("size").(int)
	if !assert.Equal(t, 1, size) {
		t.Errorf("Activity should have returned size 1")
		t.Fail()
	}
	if !assert.NotNil(t, key) {
		t.Errorf("Activity should have returned a key")
		t.Fail()
	}
	//fmt.Printf("Append with object got key: %s and size: %d\n", key, size)
}

func TestEvalAppendNoKeyOrObj(t *testing.T) {
	act, err := New(newInitContext(map[string]interface{}{
		"operation": "append",
	}))
	if err != nil {
		t.Errorf("Could not create initial context:  %s", err)
		t.Fail()
	}
	tc := test.NewActivityContext(act.Metadata())
	ok, err := act.Eval(tc)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("Could not execute activty:  %s", err)
		t.Fail()
	}
	if !ok {
		t.Errorf("Activity returned false")
	}
	key := tc.GetOutput("key").(string)
	size := tc.GetOutput("size").(int)
	if !assert.Equal(t, 0, size) {
		t.Errorf("Activity should have returned size 0")
		t.Fail()
	}
	if !assert.NotNil(t, key) {
		t.Errorf("Activity should have returned a key")
		t.Fail()
	}
	//fmt.Printf("Append with object got key: %s and size: %d\n", key, size)
}

func TestEvalEndToEnd(t *testing.T) {
	act, err := New(newInitContext(map[string]interface{}{
		"operation": "append",
	}))
	if err != nil {
		t.Errorf("Could not create initial context:  %s", err)
		t.Fail()
	}
	tc := test.NewActivityContext(act.Metadata())
	ok, err := act.Eval(tc)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("Could not execute activty:  %s", err)
		t.Fail()
	}
	if !ok {
		t.Errorf("Activity returned false")
	}
	key := tc.GetOutput("key").(string)
	size := tc.GetOutput("size").(int)
	if !assert.Equal(t, 0, size) {
		t.Errorf("Activity should have returned size 0")
		t.Fail()
	}
	if !assert.NotNil(t, key) {
		t.Errorf("Activity should have returned a key")
		t.Fail()
	}

	//Append an obj
	tc.SetInput("object", getObj)
	tc.SetInput("key", key)
	ok, err = act.Eval(tc)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("Could not execute activty:  %s", err)
		t.Fail()
	}
	if !ok {
		t.Errorf("Activity returned false")
	}
	size = tc.GetOutput("size").(int)
	if !assert.Equal(t, 1, size) {
		t.Errorf("Activity should have returned size 1")
		t.Fail()
	}
	if !assert.Equal(t, key, tc.GetInput("key").(string)) {
		t.Errorf("Activity should have returned a key")
		t.Fail()
	}

	//Append second obj
	tc.SetInput("object", getObj)
	tc.SetInput("key", key)
	ok, err = act.Eval(tc)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("Could not execute activty:  %s", err)
		t.Fail()
	}
	if !ok {
		t.Errorf("Activity returned false")
	}
	size = tc.GetOutput("size").(int)
	if !assert.Equal(t, 2, size) {
		t.Errorf("Activity should have returned size 2")
		t.Fail()
	}
	if !assert.Equal(t, key, tc.GetInput("key").(string)) {
		t.Errorf("Activity should have returned a key")
		t.Fail()
	}

	//Get the collection
	act, err = New(newInitContext(map[string]interface{}{
		"operation": "get",
	}))
	tc.SetInput("key", key)
	ok, err = act.Eval(tc)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("Could not execute activty:  %s", err)
		t.Fail()
	}
	if !ok {
		t.Errorf("Activity returned false")
	}
	collection := tc.GetOutput("collection")
	if !assert.Equal(t, 2, len(collection.([]interface{}))) {
		t.Errorf("Returned collection length should be 2 not %d", len(collection.([]interface{})))
	}
	size = tc.GetOutput("size").(int)
	if !assert.Equal(t, 2, size) {
		t.Errorf("Activity should have returned size 2")
		t.Fail()
	}

	//Delete the collection
	act, err = New(newInitContext(map[string]interface{}{
		"operation": "delete",
	}))
	tc.SetInput("key", key)
	ok, err = act.Eval(tc)
	assert.Nil(t, err)
	if err != nil {
		t.Errorf("Could not execute activty:  %s", err)
		t.Fail()
	}
	if !ok {
		t.Errorf("Activity returned false")
	}
	size = tc.GetOutput("size").(int)
	if !assert.Equal(t, -1, size) {
		t.Errorf("Activity should have returned size -1")
		t.Fail()
	}
}
