package collection

import (
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

var collectionCacheMutex sync.Mutex

// Collection static structure containing all aggregations.
type Collection struct {
	metadata  *activity.Metadata
	colmap    map[string][]interface{}
	generator *util.Generator
}

var col *Collection

func init() {
	col = new(Collection)
	col.colmap = make(map[string][]interface{})
	gen, err := util.NewGenerator()
	if err == nil {
		col.generator = gen
	}
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &Collection{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *Collection) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *Collection) Eval(context activity.Context) (done bool, err error) {

	// do eval

	return true, nil
}
