package collargo

import (
	"encoding/json"

	"github.com/satori/go.uuid"
	"reflect"
)

// AnonPayload the payload name for anonymous payload
const AnonPayload string = "__anon__"

// Signal The signal structure, represents a Signal.
// Signal is an envelope to deliver data through collar graphs
type Signal struct {
	ID      string                 // the signal id
	Seq     string                 // alias of the signal id
	Error   error                  // represent an error signal
	End     bool                   // represent an end signal
	Payload map[string]interface{} // signal payload, used for communicating from upstream to downstream
	Tags    map[string]string      // the signal tags
}

// SignalPayload Interface represents signal payload
type SignalPayload interface {
	// convert interface{} to current struct
	Convert(v interface{}) error
}

// CreateSignal Signal creator, create a signal
//
// id: the id of the signal, if "", a uuid will be assigned
func createSignal(
	id string,
	payload map[string]interface{},
	tags map[string]string,
	err error,
	end bool) Signal {

	// assign a new id if id is empty (zero for string type
	if id == "" {
		id = uuid.NewV1().String()
	}

	return Signal{
		ID:      id,
		Seq:     id,
		Error:   err,
		End:     end,
		Payload: payload,
		Tags:    tags,
	}
}

// CreateSignal create a signal with a payload (map[string]interface{})
//
// if data is of type map[string]interface{}, the signal's payload is the data
//
// if data is of type error, the signal's error is the data
//
// if data is of other type, the data is the anonymous payload
func CreateSignal(data interface{}) Signal {
	var payload map[string]interface{}
	var err error

	if reflect.TypeOf(data).Kind() == reflect.Map {
		payload = data.(map[string]interface{})
	} else if signal, ok := data.(Signal); ok {
		// do nothing, data already be signal
		return signal
	} else if e, ok := data.(error); ok {
		err = e
		payload = map[string]interface{}{}
	} else {
		payload = make(map[string]interface{})
		payload[AnonPayload] = data
	}

	id := uuid.NewV1().String()
	s := Signal{
		ID:      id,
		Seq:     id,
		Error:   err,
		End:     false,
		Payload: payload,
		Tags:    map[string]string{},
	}
	return s
}

// CreateSignalFromJSON Create a signal from json string, the payload is decoded as a map
/*
func CreateSignalFromJSON(jsonStr string) (Signal, error) {
  newSignal := Signal{}
  // general signal payload type (map and array)
  err := json.Unmarshal([]byte(jsonStr), &newSignal)
  return newSignal, err
}
*/

// Clone the current signal
func (s Signal) Clone() Signal {
	return s.New(nil)
}

// New Create a new signal from current signal, keeping the id, seq, and tags, but with different payload
func (s Signal) New(data interface{}) Signal {
	var payload map[string]interface{}
	var err error

	if data == nil {
		payload = nil
	} else if reflect.TypeOf(data).Kind() == reflect.Map {
		payload = data.(map[string]interface{})
	} else if signal, ok := data.(Signal); ok {
		// do nothing, data already be signal
		return signal
	} else if e, ok := data.(error); ok {
		err = e
		payload = map[string]interface{}{}
	} else {
		payload = make(map[string]interface{})
		payload[AnonPayload] = data
	}

	// copy the tags
	copiedTag := map[string]string{}
	for k, v := range s.Tags {
		copiedTag[k] = v
	}

	// copy payload if necessary, or get new payload from argument
	var newPayload map[string]interface{}
	if payload == nil {
		newPayload := make(map[string]interface{})
		// copy payload
		for k, v := range s.Payload {
			newPayload[k] = v
		}
	} else {
		newPayload = payload
	}

	if err == nil {
		err = s.Error
	}

	newSignal := Signal{
		ID:      s.ID,
		Seq:     s.Seq,
		Error:   err,
		End:     s.End,
		Payload: newPayload,
		Tags:    copiedTag,
	}

	return newSignal
}

// Get get the  value in the payload with a key
// return the corresponding payload and true for status, otherwise nil, false
func (s Signal) Get(name string) (v interface{}, existed bool) {
	v, existed = s.Payload[name]
	return v, existed
}

// GetValue Get a payload with name and convert it to a struct
func (s *Signal) GetValue(name string, value SignalPayload) (existed bool, err error) {
	v, existed := s.Payload[name]

	if !existed {
		return existed, nil
	}

	err = value.Convert(v)

	return existed, err
}

// Set a payload with given key, returns a new signal containing the new key - value pair as payload
func (s Signal) Set(key string, value interface{}) Signal {
	newPayload := make(map[string]interface{})

	for k, v := range s.Payload {
		newPayload[k] = v
	}

	newPayload[key] = value

	return s.New(newPayload)
}

// SetResult set the special payload result
func (s Signal) SetResult(result interface{}) Signal {
	return s.Set("__result__", result)
}

// GetResult get the result payload
func (s Signal) GetResult() (v interface{}, existed bool) {
	v, existed = s.Get("__result__")
	return v, existed
}

// SetError set error and returns a new signal represents it. The new signal keeps the original payload
func (s Signal) SetError(err error) Signal {
	newSignal := Signal{
		ID:      s.ID,
		Seq:     s.Seq,
		Error:   err,
		End:     s.End,
		Payload: s.Payload,
		Tags:    s.Tags,
	}
	return newSignal
}

// GetTag  get a tag with tag name
// return the tag value, and ok status true, otherwise "" and false
func (s Signal) GetTag(name string) (tag string, ok bool) {
	tag, ok = s.Tags[name]
	return tag, ok
}

// SetTag Set a new tag in the signal
func (s Signal) SetTag(name string, value string) Signal {
	newSignal := s.Clone()
	newSignal.Tags[name] = value
	return newSignal
}

// DelTag delete a tag with name
func (s Signal) DelTag(name string) Signal {
	newSignal := s.Clone()
	delete(newSignal.Tags, name)
	return newSignal
}

// ToJSON serialize the signal as a json string
func (s Signal) ToJSON() (string, error) {
	jsonByte, err := json.Marshal(s)
	return string(jsonByte), err
}
