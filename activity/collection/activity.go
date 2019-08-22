package flogo-collection

import (
	"fmt"
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

// newKey create a new collectin key
func (collection *Collection) newKey() (string, error) {
	if collection.generator == nil {
		return "", fmt.Errorf("Generator not initialized")
	}
	return collection.generator.NextAsString(), nil
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
	key := context.GetInput("key")
	object := context.GetInput("object")
	operation := context.GetInput("operation").(string)

	switch operation {
	case "append":
		if key == nil {
			key, err = a.newKey()
			if err != nil {
				return false, fmt.Errorf("Append with no key failed to create dynamic key for reason [%s]", err)
			}
		}
		if object == nil {
			if err != nil {
				return false, fmt.Errorf("Append called with a nil object")
			}
		}
		a.colmap[key.(string)] = append(a.colmap[key.(string)], object)

	case "get":
		if key == nil {
			return false, fmt.Errorf("Get called with no key")
		}
		col, ok := a.colmap[key.(string)]
		if !ok {
			return false, fmt.Errorf("Get called for invalid key: %s", key.(string))
		}
		context.SetOutput("collection", col)
		context.SetOutput("size", len(a.colmap[key.(string)]))
		return true, nil

	case "delete":
		if key == nil {
			return false, fmt.Errorf("Get called with no key")
		}
		delete(a.colmap, key.(string))

	default:

	}
	return true, nil
}
