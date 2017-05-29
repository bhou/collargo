package collargo

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type userType struct {
	Name string
	Age  int
}

func (u *userType) Convert(v interface{}) (err error) {
	user, ok := v.(userType)
	if !ok {
		err = errors.New("Can't convert interface{} to userType")
	}

	u.Name = user.Name
	u.Age = user.Age

	return err
}

/* Examples */

func ExampleSignal_GetValue() {
	s := createSignal(
		"",
		make(map[string]interface{}),
		make(map[string]string),
		nil,
		false)

	s = s.Set("user", userType{"Mike", 25})

	// Method 1. get interface{} type payload, and convert it to struct later
	v, existed := s.Get("user")
	u1 := v.(userType)
	fmt.Println(u1.Name, u1.Age, existed) // Mike 25 true

	// Method 2. get struct instance directly with GetValue
	u2 := userType{}
	existed, err := s.GetValue("user", &u2)
	fmt.Println(u2.Name, u2.Age, err, existed) // Mike 25 <nil> true

	// u1 & u2 are copies, they are not the same object
	u2.Name = "Jack"
	fmt.Println(u1.Name, u2.Name) // Mike Jack

	// Output:
	// Mike 25 true
	// Mike 25 <nil> true
	// Mike Jack
}

/* Tests */

func TestCreateSignal(t *testing.T) {
	s := createSignal(
		"id",
		map[string]interface{}{"payload1": "payloadValue1"},
		map[string]string{"tag1": "value1"},
		nil,
		false,
	)

	assert.Equal(t, "id", s.ID)

	v1, ok := s.Get("payload1")
	assert.Equal(t, "payloadValue1", v1)
	assert.True(t, ok)

	t1, ok := s.GetTag("tag1")
	assert.Equal(t, "value1", t1)
	assert.True(t, ok)
}

func TestSetPayload(t *testing.T) {
	s := createSignal(
		"id",
		map[string]interface{}{"payload1": "payloadValue1"},
		map[string]string{"tag1": "value1"},
		nil,
		false,
	)

	assert.Equal(t, "id", s.ID)

	newS := s.Set("test", "test-value")

	// test original signal payload
	payload1, existed := s.Get("payload1")
	assert.True(t, existed)
	assert.Equal(t, "payloadValue1", payload1)

	// original signal should not be modified
	_, existed = s.Get("test")
	assert.False(t, existed)

	// new signal contains the new payload
	testPayload, existed := newS.Get("test")
	assert.True(t, existed)
	assert.Equal(t, "test-value", testPayload)
}

func TestCreateSignalWithPayload(t *testing.T) {
	user := userType{"John", 25}

	// test anonymous payload
	s1 := CreateSignal(user)

	u1, existed := s1.Get(AnonPayload)

	assert.Equal(t, "John", u1.(userType).Name)
	assert.Equal(t, 25, u1.(userType).Age)

	// test map payload
	m := map[string]interface{}{
		"user": user,
	}

	s2 := CreateSignal(m)
	u2, existed := s2.Get("user")

	assert.Equal(t, "John", u2.(userType).Name)
	assert.Equal(t, 25, u2.(userType).Age)

	// test signal payload
	s3 := CreateSignal(s2)
	u3, existed := s3.Get("user")

	assert.Equal(t, "John", u3.(userType).Name)
	assert.Equal(t, 25, u3.(userType).Age)

	// test error signal
	err := errors.New("Test error")
	s4 := CreateSignal(err)

	_, existed = s4.Get("user")
	assert.False(t, existed)
	assert.Equal(t, err, s4.Error)

	s5 := CreateSignal("string payload")
	strPayload, existed := s5.Get(AnonPayload)
	assert.True(t, existed)
	assert.Equal(t, "string payload", strPayload)
}

func TestSetStructPayload(t *testing.T) {
	s := createSignal(
		"id",
		map[string]interface{}{},
		map[string]string{},
		nil,
		false,
	)

	user := userType{"John", 25}

	s = s.Set("user", user)

	v, existed := s.Get("user")
	assert.True(t, existed)
	assert.NotNil(t, v)

	// convert type after calling Get
	u1, ok := v.(userType)
	assert.True(t, ok)
	assert.Equal(t, "John", u1.Name)

	// get user with GetValue
	u2 := userType{}
	ok, err := s.GetValue("user", &u2)
	assert.Nil(t, err)
	assert.Equal(t, "John", u2.Name)
	assert.EqualValues(t, 25, u2.Age)

	// the 2 user are not the same
	u1.Name = "Lock"
	u2.Age = 32
	assert.Equal(t, "Lock", u1.Name)
	assert.Equal(t, "John", u2.Name)
	assert.Equal(t, 25, u1.Age)
	assert.Equal(t, 32, u2.Age)

	// get unexited user value
	u3 := userType{}
	ok, _ = s.GetValue("not_existed", &u3)
	assert.False(t, ok)
}

func TestSetError(t *testing.T) {
	err := errors.New("test error")
	s := createSignal(
		"id",
		map[string]interface{}{"payload1": "payloadValue1"},
		map[string]string{"tag1": "value1"},
		nil,
		false,
	)

	assert.Nil(t, s.Error)

	s1 := s.SetError(err)

	assert.Nil(t, s.Error)
	assert.Equal(t, err, s1.Error)
}

func TestToJSON(t *testing.T) {
	s := createSignal(
		"id",
		map[string]interface{}{"payload1": "payloadValue1"},
		map[string]string{"tag1": "value1"},
		nil,
		false,
	)

	jsonStr, err := s.ToJSON()

	assert.Nil(t, err)
	assert.Equal(t, `{"ID":"id","Seq":"id","Error":null,"End":false,"Payload":{"payload1":"payloadValue1"},"Tags":{"tag1":"value1"}}`, jsonStr)

	// test interface{} payload
	type user struct {
		Name    string
		Age     int8
		Address string
	}
	s2 := s.Set("interfacePayload", user{"John", 27, "7 Rue John Smith, New York"})
	jsonStr2, err := s2.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, `{"ID":"id","Seq":"id","Error":null,"End":false,"Payload":{"interfacePayload":{"Name":"John","Age":27,"Address":"7 Rue John Smith, New York"},"payload1":"payloadValue1"},"Tags":{"tag1":"value1"}}`, jsonStr2)
}

func TestManipulateTag(t *testing.T) {
	s := createSignal(
		"id",
		map[string]interface{}{"payload1": "payloadValue1"},
		map[string]string{"tag1": "value1"},
		nil,
		false,
	)

	// check original signal tags
	tag1, existed := s.GetTag("tag1")
	assert.True(t, existed)
	assert.Equal(t, "value1", tag1)

	// add a new tag
	s2 := s.SetTag("tag2", "value2")
	tag1, existed = s2.GetTag("tag1")
	assert.True(t, existed)
	assert.Equal(t, "value1", tag1)
	tag2, existed := s2.GetTag("tag2")
	assert.True(t, existed)
	assert.Equal(t, "value2", tag2)

	// remove the tag
	s3 := s2.DelTag("tag1")
	tag1, existed = s3.GetTag("tag1")
	assert.False(t, existed)
	tag2, existed = s2.GetTag("tag2")
	assert.True(t, existed)
	assert.Equal(t, "value2", tag2)
}

/*
func TestCreateGenericSignalFromJSON(t *testing.T) {
  s := createSignal(
    "id",
    map[string]interface{}{"payload1": "payloadValue1"},
    map[string]string{"tag1": "value1"},
    nil,
    false,
  )

  jsonStr, _ := s.ToJSON()

  s1, _ := createSignalFromJSON(jsonStr, nil)

  assert.Equal(t, "id", s1.ID)
  payload1, _ := s1.Get("payload1")
  assert.Equal(t, "payloadValue1", payload1)
  tag1, _ := s1.GetTag("tag1")
  assert.Equal(t, "value1", tag1)
  assert.Nil(t, s1.Error)
  assert.False(t, s1.End)
}

func TestCreateTypedSignalFromJSON(t *testing.T) {
  type user struct {
    Name string
    Age  int8
  }

  s := createSignal(
    "id",
    map[string]interface{}{
      "payload1": "payloadValue1",
      "user": user{
        Name: "John",
        Age:  27,
      },
      "strField": "Hello",
      "intField": 100,
    },
    map[string]string{"tag1": "value1"},
    nil,
    false,
  )

  jsonStr, err := s.ToJSON()

  assert.Nil(t, err)

  newPayload := map[string]interface{}{
    "user": user{},
  }
  CreateSignalFromJSON(jsonStr, &newPayload)
}
*/
