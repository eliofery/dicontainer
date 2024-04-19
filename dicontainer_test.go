package dicontainer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Represents a test structures.
type testStruct struct {
	Key   string
	Value int
}

type testStructPointer struct {
	Key   string
	Value int
}

type testOtherStruct struct {
	Key   string
	Value int
}

// TestSetGet_success tests the successful setting and getting of dependencies.
func TestSetGet_success(t *testing.T) {
	cases := []struct {
		name     string
		creators []any
	}{
		{
			name: "without arguments",
			creators: []any{
				func() testStruct {
					return testStruct{
						Key:   "some-key",
						Value: 123,
					}
				},
			},
		},
		{
			name: "with arguments",
			creators: []any{
				func(ts testStruct) *testStructPointer {
					return &testStructPointer{
						Key:   ts.Key,
						Value: ts.Value,
					}
				},
				func() testStruct {
					return testStruct{
						Key:   "some-key",
						Value: 123,
					}
				},
				func(tsp *testStructPointer) *testOtherStruct {
					return &testOtherStruct{
						Key:   tsp.Key,
						Value: tsp.Value,
					}
				},
			},
		},
		{
			name: "with arguments 2",
			creators: []any{
				func() *testStructPointer {
					return &testStructPointer{
						Key:   "some-key",
						Value: 123,
					}
				},
				func() testStruct {
					return testStruct{
						Key:   "some-key",
						Value: 123,
					}
				},
				func(ts testStruct, tsp *testStructPointer) *testOtherStruct {
					return &testOtherStruct{
						Key:   ts.Key,
						Value: tsp.Value,
					}
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			di := New()

			err := di.Set(test.creators...)
			assert.NoError(t, err)

			ts := di.Get("testStruct").(testStruct)
			assert.Equal(t, ts.Key, "some-key")
			assert.Equal(t, ts.Value, 123)

			tos, ok := di.Get("testOtherStruct").(*testOtherStruct)
			if ok {
				assert.Equal(t, tos.Key, "some-key")
				assert.Equal(t, tos.Value, 123)
			}
		})
	}
}

// TestSetGet_error tests error handling during setting and getting of dependencies.
func TestSetGet_error(t *testing.T) {
	cases := []struct {
		name     string
		creators []any
	}{
		{
			name: "without arguments",
			creators: []any{
				func() (testStruct, error) {
					return testStruct{
						Key:   "some-key",
						Value: 123,
					}, errors.New("some error")
				},
			},
		},
		{
			name: "with arguments",
			creators: []any{
				func(ts testStruct) *testStructPointer {
					return &testStructPointer{
						Key:   ts.Key,
						Value: ts.Value,
					}
				},
				func() testStruct {
					return testStruct{
						Key:   "some-key",
						Value: 123,
					}
				},
				func(tsp *testStructPointer) (*testOtherStruct, error) {
					return &testOtherStruct{
						Key:   tsp.Key,
						Value: tsp.Value,
					}, errors.New("some error")
				},
			},
		},
		{
			name: "with arguments 2",
			creators: []any{
				func() *testStructPointer {
					return &testStructPointer{
						Key:   "some-key",
						Value: 123,
					}
				},
				func() testStruct {
					return testStruct{
						Key:   "some-key",
						Value: 123,
					}
				},
				func(ts testStruct, tsp *testStructPointer) (*testOtherStruct, error) {
					return &testOtherStruct{
						Key:   ts.Key,
						Value: tsp.Value,
					}, errors.New("some error")
				},
			},
		},
		{
			name: "incorrect arguments",
			creators: []any{
				func() (int, error) {
					return 123, nil
				},
				func() any {
					return testStruct{
						Key:   "some-key",
						Value: 123,
					}
				},
				func() (any, int, error) {
					return testStruct{
						Key:   "some-key",
						Value: 123,
					}, 123, errors.New("some error")
				},
				func() (any, int, error) {
					return testStruct{
						Key:   "some-key",
						Value: 123,
					}, 123, nil
				},
			},
		},
		{
			name: "missing arguments",
			creators: []any{
				func(ts testStruct) (testStruct, error) {
					ts.Key = "other-key"
					return ts, nil
				},
			},
		},
		{
			name: "invalid creator type",
			creators: []any{
				123,
				"bugagaga",
				[]int{1, 2, 3},
				struct{ foo string }{"bar"},
			},
		},
		{
			name: "invalid creator, no return value",
			creators: []any{
				func() {},
			},
		},
		{
			name: "invalid dependency type",
			creators: []any{
				func(ts testStruct) (int, error) {
					return 123, nil
				},
				func() testStruct {
					return testStruct{}
				},
			},
		},
		{
			name: "invalid result, expected 1 or 2 arguments",
			creators: []any{
				func(ts testStruct, tos testOtherStruct) (testStruct, testOtherStruct, error) {
					return ts, tos, nil
				},
				func() testStruct {
					return testStruct{}
				},
				func() testOtherStruct {
					return testOtherStruct{}
				},
			},
		},
		{
			name: "invalid result, expected 2nd argument to be an error (pointer)",
			creators: []any{
				func(ts testStruct, tos testOtherStruct) (testStruct, *testOtherStruct) {
					return ts, &tos
				},
				func() testStruct {
					return testStruct{}
				},
				func() testOtherStruct {
					return testOtherStruct{}
				},
			},
		},
		{
			name: "invalid result, expected 2nd argument to be an error",
			creators: []any{
				func(ts testStruct, tos testOtherStruct) (testStruct, testOtherStruct) {
					return ts, tos
				},
				func() testStruct {
					return testStruct{}
				},
				func() testOtherStruct {
					return testOtherStruct{}
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			di := New()

			err := di.Set(test.creators...)
			assert.Error(t, err)
		})
	}
}
