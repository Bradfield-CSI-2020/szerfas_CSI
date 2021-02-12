package skiplist

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

const (
	Get = iota
	Put
	Delete

	NumOps
)

type action struct {
	op         int
	key, value string
}

func TestSkipList(t *testing.T) {
	params := gopter.DefaultTestParameters()
	params.MinSuccessfulTests = 10000

	properties := gopter.NewProperties(params)

	properties.Property("get / put results in same values as built-in map", prop.ForAll(
		func(actions []action) string {
			m := make(map[string]string)
			s := newSkipListOC()

			for _, action := range actions {
				if action.op == Get {
					mValue, mOk := m[action.key]
					sValue, sOk := s.Get(action.key)
					if mOk != sOk || (mOk && sOk && mValue != sValue) {
						return fmt.Sprintf("mismatch between map and skiplist for key %v", action.key)
					}
				} else if action.op == Put {
					m[action.key] = action.value
					s.Put(action.key, action.value)
				} else if action.op == Delete {
					delete(m, action.key)
					s.Delete(action.key)
				} else {
					panic(fmt.Sprintf("Unknown op %v", action.op))
				}
			}

			return ""
		},
		gen.SliceOf(gen.Struct(reflect.TypeOf(action{}), map[string]gopter.Gen{
			"op":    gen.IntRange(0, NumOps-1),
			"key":   gen.NumString(),
			"value": gen.NumString(),
		})),
	))

	properties.TestingRun(t)
}
