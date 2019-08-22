package collection

import (
	"fmt"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/util"
)

var collectionCacheMutex sync.Mutex

type CollectionActivity struct {
	metadata  *activity.Metadata
	generator *util.Generator
}

// Collection static structure containing all aggregations.
type Collection struct {
	colmap map[string][]interface{}
}

var col *Collection

// newKey create a new collectin key
func (collection *CollectionActivity) newKey() (res string, err error) {
	if collection.generator == nil {
		collection.generator, err = util.NewGenerator()
		if err != nil {
			return "", err
		}
	}
	return collection.generator.NextAsString(), nil
}

func (collection *CollectionActivity) initialize() {
	if col == nil {
		col = new(Collection)
		col.colmap = make(map[string][]interface{})
	}

}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &CollectionActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (collection *CollectionActivity) Metadata() *activity.Metadata {
	return collection.metadata
}

// Eval implements activity.Activity.Eval
func (collection *CollectionActivity) Eval(context activity.Context) (done bool, err error) {

	// do eval
	key := context.GetInput("key")
	object := context.GetInput("object")
	operation := context.GetInput("operation").(string)

	switch operation {
	case "append":
		if key == nil {
			key, err = collection.newKey()
			if err != nil {
				return false, fmt.Errorf("Append with no key failed to create dynamic key for reason [%s]", err)
			}
		}
		if object == nil {
			if err != nil {
				return false, fmt.Errorf("Append called with a nil object")
			}
		}
		col.colmap[key.(string)] = append(col.colmap[key.(string)], object)

	case "get":
		if key == nil {
			return false, fmt.Errorf("Get called with no key")
		}
		array, ok := col.colmap[key.(string)]
		if !ok {
			return false, fmt.Errorf("Get called for invalid key: %s", key.(string))
		}
		context.SetOutput("collection", array)
		context.SetOutput("size", len(col.colmap[key.(string)]))
		return true, nil

	case "delete":
		if key == nil {
			return false, fmt.Errorf("Get called with no key")
		}
		delete(col.colmap, key.(string))

	default:

	}
	return true, nil
}
