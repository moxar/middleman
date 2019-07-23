package contextm

import (
	"context"
	"testing"
)

func TestRange(t *testing.T) {
	ctx := context.TODO()

	type ctxKey string

	ctx = WithRange(ctx, "first Avenger", "Steve Rogers")
	ctx = WithRange(ctx, "expert tinker", "Tony Stark")
	ctx = context.WithValue(ctx, ctxKey("not an Avenger"), "Batman") // nolint: vet
	ctx = WithRange(ctx, "friendly neighbour", "Peter Parker")

	t.Run("With multiple elements to Range", func(t *testing.T) {

		var i int
		Range(ctx, func(key, val interface{}) bool {
			switch i {
			case 0:
				if key != "first Avenger" || val != "Steve Rogers" {
					t.Error("error on key", i)
				}

			case 1:
				if key != "expert tinker" || val != "Tony Stark" {
					t.Error("error on key", i)
				}

			case 2:
				if key != "friendly neighbour" || val != "Peter Parker" {
					t.Error("error on key", i)
				}
			}
			i++
			return true
		})
		if i != 3 {
			t.Fail()
			return
		}
	})

	t.Run("With multiple Keys", func(t *testing.T) {

		keys := []interface{}{
			"first Avenger",
			"expert tinker",
			"friendly neighbour",
		}

		for i, k := range Keys(ctx) {
			if keys[i] != k {
				t.Errorf("%s is different from %s", keys[i], k)
			}
		}
	})
}
