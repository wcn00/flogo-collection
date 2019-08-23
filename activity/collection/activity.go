package collection

import (
	"fmt"
	"sync"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/support"
)

var collectionCacheMutex sync.Mutex

//CollectionActivity is structure for collection parms
type Activity struct {
	Operation string      `md:"operation"`
	Key       string      `md:"key"`
	Object    interface{} `md:"object"`
}

type ActivityOutput struct {
	Key        string        `md:"key"`
	Collection []interface{} `md:"collection"`
	Size       int           `md:"size"`
}

// FromMap converts the values from a map into the struct Output
func (o *ActivityOutput) FromMap(values map[string]interface{}) error {
	key, err := coerce.ToString(values["key"])
	if err != nil {
		return err
	}
	o.Key = key
	collection, err := coerce.ToArray(values["collection"])
	if err != nil {
		return err
	}
	o.Collection = collection
	size, err := coerce.ToInt(values["size"])
	if err != nil {
		return err
	}
	o.Size = size
	return nil
}

// ToMap converts the struct Output into a map
func (o *ActivityOutput) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"key":    o.Key,
		"":       o.Collection,
		"result": o.Size,
	}
}

// Collection static structure containing all aggregations.
type Collection struct {
	colmap    map[string][]interface{}
	generator *support.Generator
}

var col *Collection

func init() {
	col = new(Collection)
	col.colmap = make(map[string][]interface{})
	col.generator, _ = support.NewGenerator()
	_ = activity.Register(&Activity{})
}

// newKey create a new collectin key
func (collection *Activity) newKey() (res string, err error) {
	if col.generator == nil {
		col.generator, err = support.NewGenerator()
		if err != nil {
			return "", fmt.Errorf("Failed to generate a dynamic key for collection for reason [%s]", err)
		}
	}
	return col.generator.NextAsString(), nil
}

var collectionActivityMd = activity.ToMetadata()

// Metadata implements activity.Activity.Metadata
func (collection *Activity) Metadata() *activity.Metadata {
	return collectionActivityMd
}

// Eval implements activity.Activity.Eval
func (collection *Activity) Eval(context activity.Context) (done bool, err error) {
	collectionCacheMutex.Lock()
	defer collectionCacheMutex.Unlock()

	// do eval
	key := context.GetInput("key")
	object := context.GetInput("object")
	operation := context.GetInput("operation")
	output := &ActivityOutput{}
	switch operation.(string) {
	case "append":
		if key == nil {
			key, err = collection.newKey()
			if err != nil {
				return false, fmt.Errorf("Append with no key failed to create dynamic key for reason [%s]", err)
			}
		}
		if object != nil {
			col.colmap[key.(string)] = append(col.colmap[key.(string)], object)
		}

		output.Size = len(col.colmap[key.(string)])
		output.Key = key.(string)
		context.SetOutputObject(output)
		return true, nil

	case "get":
		if key == nil {
			return false, fmt.Errorf("Get called with no key")
		}
		array, ok := col.colmap[key.(string)]
		if !ok {
			return false, fmt.Errorf("Get called for invalid key: %s", key.(string))
		}
		output.Size = len(col.colmap[key.(string)])
		output.Key = key.(string)
		output.Collection = array
		context.SetOutputObject(output)
		return true, nil

	case "delete":
		if key == nil {
			return false, fmt.Errorf("Get called with no key")
		}
		delete(col.colmap, key.(string))
		context.SetOutputObject(output)
		return true, nil

	default:

	}
	return true, nil
}
