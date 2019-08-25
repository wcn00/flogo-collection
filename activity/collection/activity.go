package collection

import (
	"fmt"
	"sync"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/core/support"
)

var collectionCacheMutex sync.Mutex

//ActivitySettings is structure for collection parms
type Settings struct {
	Operation string `md:"operation"`
}

// FromMap converts the values from a map into the struct Output
func (o *Settings) FromMap(values map[string]interface{}) error {
	operation, err := coerce.ToString(values["operation"])
	if err != nil {
		return err
	}
	o.Operation = operation
	return nil
}

// ToMap converts the struct Output into a map
func (o *Settings) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"operation": o.Operation,
	}
}

//ActivitySettings is structure for collection parms
type ActivityInput struct {
	Key    string      `md:"key"`
	Object interface{} `md:"object"`
}

// FromMap converts the values from a map into the struct Output
func (o *ActivityInput) FromMap(values map[string]interface{}) error {
	key, err := coerce.ToString(values["key"])
	if err != nil {
		return err
	}
	o.Key = key
	object, err := coerce.ToObject(values["object"])
	if err != nil {
		return err
	}
	o.Object = object
	return nil
}

// ToMap converts the struct Output into a map
func (o *ActivityInput) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"key":    o.Key,
		"object": o.Object,
	}
}

// ActivityOutput activity output
type Output struct {
	Key        string        `md:"key"`
	Collection []interface{} `md:"collection"`
	Size       int           `md:"size"`
}

var activityMd = activity.ToMetadata(&Settings{}, &Output{})

// FromMap converts the values from a map into the struct Output
func (o *Output) FromMap(values map[string]interface{}) error {
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
func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"key":        o.Key,
		"collection": o.Collection,
		"size":       o.Size,
	}
}

// Collection static structure containing all aggregations.
type Collection struct {
	colmap map[string][]interface{}
}

var col *Collection

func init() {
	col = new(Collection)
	col.colmap = make(map[string][]interface{})
	_ = activity.Register(&Activity{}, New)
}

// Activity base activty type
type Activity struct {
	operation string
	generator *support.Generator
}

func (collection *Activity) newKey() (res string, err error) {
	if collection.generator == nil {
		collection.generator, err = support.NewGenerator()
		if err != nil {
			return "", fmt.Errorf("Failed to generate a dynamic key for collection for reason [%s]", err)
		}
	}
	return collection.generator.NextAsString(), nil
}

// New creates a new javascript activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	settings := Settings{}
	err := metadata.MapToStruct(ctx.Settings(), &settings, true)
	if err != nil {
		return nil, err
	}
	act := Activity{
		operation: settings.Operation,
	}
	return &act, nil
}

var collectionActivityMd = activity.ToMetadata(&Settings{}, &Output{})

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
	//	output := &ActivityOutput{}
	switch collection.operation {
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

		err = context.SetOutput("size", len(col.colmap[key.(string)]))
		if err != nil {
			return false, fmt.Errorf("Append failed to set output \"size\" for reason [%s]", err)
		}
		err = context.SetOutput("key", key)
		if err != nil {
			return false, fmt.Errorf("Append failed to set output \"key\" for reason [%s]", err)
		}
		return true, nil

	case "get":
		if key == nil {
			return false, fmt.Errorf("Get called with no key")
		}
		array, ok := col.colmap[key.(string)]
		if !ok {
			return false, fmt.Errorf("Get called for invalid key: %s", key.(string))
		}
		context.SetOutput("size", len(col.colmap[key.(string)]))
		context.SetOutput("key", key)
		context.SetOutput("collection", array)
		return true, nil

	case "delete":
		if key == nil {
			return false, fmt.Errorf("Get called with no key")
		}
		delete(col.colmap, key.(string))
		context.SetOutput("size", -1)
		return true, nil

	default:

	}
	return true, nil
}
