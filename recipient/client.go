// Package recipient provides the /recipients APIs
package recipient

import (
	"net/url"
	"strconv"

	stripe "github.com/stripe/stripe-go"
)

const (
	Individual stripe.RecipientType = "individual"
	Corp       stripe.RecipientType = "corporation"

	NewAccount       stripe.BankAccountStatus = "new"
	VerifiedAccount  stripe.BankAccountStatus = "verified"
	ValidatedAccount stripe.BankAccountStatus = "validated"
	ErroredAccount   stripe.BankAccountStatus = "errored"
)

// Client is used to invoke /recipients APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New POSTs a new recipient.
// For more details see https://stripe.com/docs/api#create_recipient.
func New(params *stripe.RecipientParams) (*stripe.Recipient, error) {
	return getC().New(params)
}

func (c Client) New(params *stripe.RecipientParams) (*stripe.Recipient, error) {
	body := &url.Values{
		"name": {params.Name},
		"type": {string(params.Type)},
	}

	if params.Bank != nil {
		params.Bank.AppendDetails(body)
	}

	if len(params.Token) > 0 {
		body.Add("card", params.Token)
	} else if params.Card != nil {
		params.Card.AppendDetails(body, true)
	}

	if len(params.TaxID) > 0 {
		body.Add("tax_id", params.TaxID)
	}

	if len(params.Email) > 0 {
		body.Add("email", params.Email)
	}

	if len(params.Desc) > 0 {
		body.Add("description", params.Desc)
	}

	params.AppendTo(body)

	recipient := &stripe.Recipient{}
	err := c.B.Call("POST", "/recipients", c.Key, body, recipient)

	return recipient, err
}

// Get returns the details of a recipient.
// For more details see https://stripe.com/docs/api#retrieve_recipient.
func Get(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	return getC().Get(id, params)
}

func (c Client) Get(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}
		params.AppendTo(body)
	}

	recipient := &stripe.Recipient{}
	err := c.B.Call("GET", "/recipients/"+id, c.Key, body, recipient)

	return recipient, err
}

// Update updates a recipient's properties.
// For more details see https://stripe.com/docs/api#update_recipient.
func Update(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	return getC().Update(id, params)
}

func (c Client) Update(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	var body *url.Values

	if params != nil {
		body = &url.Values{}

		if len(params.Name) > 0 {
			body.Add("name", params.Name)
		}

		if params.Bank != nil {
			params.Bank.AppendDetails(body)
		}

		if len(params.Token) > 0 {
			body.Add("card", params.Token)
		} else if params.Card != nil {
			params.Card.AppendDetails(body, true)
		}

		if len(params.TaxID) > 0 {
			body.Add("tax_id", params.TaxID)
		}

		if len(params.DefaultCard) > 0 {
			body.Add("default_card", params.DefaultCard)
		}

		if len(params.Email) > 0 {
			body.Add("email", params.Email)
		}

		if len(params.Desc) > 0 {
			body.Add("description", params.Desc)
		}

		params.AppendTo(body)
	}

	recipient := &stripe.Recipient{}
	err := c.B.Call("POST", "/recipients/"+id, c.Key, body, recipient)

	return recipient, err
}

// Del removes a recipient.
// For more details see https://stripe.com/docs/api#delete_recipient.
func Del(id string) error {
	return getC().Del(id)
}

func (c Client) Del(id string) error {
	return c.B.Call("DELETE", "/recipients/"+id, c.Key, nil, nil)
}

// List returns a list of recipients.
// For more details see https://stripe.com/docs/api#list_recipients.
func List(params *stripe.RecipientListParams) *Iter {
	return getC().List(params)
}

func (c Client) List(params *stripe.RecipientListParams) *Iter {
	type recipientList struct {
		stripe.ListMeta
		Values []*stripe.Recipient `json:"data"`
	}

	var body *url.Values
	var lp *stripe.ListParams

	if params != nil {
		body = &url.Values{}

		if params.Verified {
			body.Add("verified", strconv.FormatBool(true))
		}

		params.AppendTo(body)
		lp = &params.ListParams
	}

	return &Iter{stripe.GetIter(lp, body, func(b url.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &recipientList{}
		err := c.B.Call("GET", "/recipients", c.Key, &b, list)

		ret := make([]interface{}, len(list.Values))
		for i, v := range list.Values {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is a iterator for list responses.
type Iter struct {
	Iter *stripe.Iter
}

// Next returns the next value in the list.
func (i *Iter) Next() (*stripe.Recipient, error) {
	r, err := i.Iter.Next()
	if err != nil {
		return nil, err
	}

	return r.(*stripe.Recipient), err
}

// Stop returns true if there are no more iterations to be performed.
func (i *Iter) Stop() bool {
	return i.Iter.Stop()
}

// Meta returns the list metadata.
func (i *Iter) Meta() *stripe.ListMeta {
	return i.Iter.Meta()
}

func getC() Client {
	return Client{stripe.GetBackend(), stripe.Key}
}
