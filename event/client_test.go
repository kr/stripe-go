package event

import (
	"testing"

	stripe "github.com/stripe/stripe-go"
	. "github.com/stripe/stripe-go/utils"
)

func init() {
	stripe.Key = GetTestKey()
}

func TestEvent(t *testing.T) {
	params := &stripe.EventListParams{}
	params.Filters.AddFilter("limit", "", "5")
	params.Single = true
	params.Type = "charge.*"

	i := List(params)
	for !i.Stop() {
		e, err := i.Next()

		if err != nil {
			t.Error(err)
		}

		if e == nil {
			t.Error("No nil values expected")
		}

		if len(e.ID) == 0 {
			t.Errorf("ID is not set\n")
		}

		if e.Created == 0 {
			t.Errorf("Created date is not set\n")
		}

		if len(e.Type) == 0 {
			t.Errorf("Type is not set\n")
		}

		if len(e.Req) == 0 {
			t.Errorf("Request is not set\n")
		}

		if e.Data == nil {
			t.Errorf("Data is not set\n")
		}

		if len(e.Data.Obj) == 0 {
			t.Errorf("Object is empty\n")
		}

		target, err := Get(e.ID)

		if err != nil {
			t.Error(err)
		}

		if e.ID != target.ID {
			t.Errorf("ID %q does not match expected id %q\n", e.ID, target.ID)
		}

		targetVal := e.GetObjValue("card", "last4")
		val := target.Data.Obj["card"].(map[string]interface{})["last4"]

		if targetVal != val {
			t.Errorf("Value %q does not match expected value %q\n", targetVal, val)
		}

		// no need to actually check the value, we're just validating this doesn't bomb
		e.GetObjValue("does not exist")
	}
}
