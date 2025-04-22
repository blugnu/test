package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

func (l Location) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`"lat:%f,lon:%f"`, l.Latitude, l.Longitude)
	return []byte(s), nil
}

type student struct {
	name     string
	age      uint8
	height   float32
	joinedAt *time.Time
}

type teacher struct {
	name     string
	age      uint8
	height   float32
	joinedAt *time.Time
}

type classRoom struct {
	name            string
	capacity        uint
	homeRoomTeacher teacher
	students        []student
}

type school struct {
	name                  string
	establishedDate       *time.Time
	classRooms            []classRoom
	descriptionUnexported []byte

	Location            Location
	DescriptionExported []byte
}

func TestMarshal(t *testing.T) {
	t.Run("should generate null when data is nil", func(t *testing.T) {
		expected := []byte(`null`)
		actual, err := Marshal(nil)

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})

	t.Run("should successfully marshal on data type", func(t *testing.T) {
		expected, _ := json.Marshal([]byte(`1234`))
		actual, err := Marshal([]byte(`1234`))

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})

	t.Run("should successfully marshal on primitive slice", func(t *testing.T) {
		expected := []byte(`[1,2,3,4]`)
		actual, err := Marshal([]int{1, 2, 3, 4})

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})

	t.Run("should successfully marshal on primitive array", func(t *testing.T) {
		expected := []byte(`[1,2,3,4]`)
		actual, err := Marshal([4]int{1, 2, 3, 4})

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})

	t.Run("should successfully marshal on an empty slice", func(t *testing.T) {
		expected, _ := json.Marshal([]school{})
		actual, err := Marshal([]school{})

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})

	t.Run("should successfully marshal given data", func(t *testing.T) {
		joinedAt := toPtr(time.Unix(1704537759, 0).UTC())
		establishedAt := toPtr(time.Unix(1504537759, 0).UTC())

		schools := []school{
			{
				name:            "School A",
				establishedDate: establishedAt,
				classRooms: []classRoom{
					{
						name:     "Class A",
						capacity: 100,
						homeRoomTeacher: teacher{
							name:     "Mrs. Diana",
							age:      30,
							height:   180.5,
							joinedAt: joinedAt,
						},
						students: []student{
							{
								name:     "Brian",
								age:      14,
								height:   170.5,
								joinedAt: joinedAt,
							},
						},
					},
				},
				descriptionUnexported: []byte(`unexported description`),
				Location: Location{
					Latitude:  1.1,
					Longitude: 2.2,
				},
				DescriptionExported: []byte(`exported description`),
			},
		}

		expected := []byte(`[{"name":"School A","establishedDate":"2017-09-04T15:09:19Z","classRooms":[{"name":"Class A","capacity":100,"homeRoomTeacher":{"name":"Mrs. Diana","age":30,"height":180.5,"joinedAt":"2024-01-06T10:42:39Z"},"students":[{"name":"Brian","age":14,"height":170.5,"joinedAt":"2024-01-06T10:42:39Z"}]}],"descriptionUnexported":"dW5leHBvcnRlZCBkZXNjcmlwdGlvbg==","Location":"lat:1.100000,lon:2.200000","DescriptionExported":"ZXhwb3J0ZWQgZGVzY3JpcHRpb24="}]`)
		actual, err := Marshal(schools)

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("got: %s, expected: %s", string(actual), string(expected))
		}

		if err != nil {
			t.Fatal("should not raise an error")
		}
	})
}

func toPtr[T any](v T) *T {
	return &v
}
