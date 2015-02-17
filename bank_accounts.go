package lob

import "strconv"

// BankAccount represents a bank account in lob's system.
type BankAccount struct {
	AccountAddress *Address `json:"account_address"`
	AccountNumber  string   `json:"account_number"`
	BankAddress    *Address `json:"bank_address"`
	BankCode       string   `json:"bank_code"`
	DateCreated    string   `json:"date_created"`
	DateModified   string   `json:"date_modified"`
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Object         string   `json:"object"`
	RoutingNumber  string   `json:"routing_number"`
	Signatory      string   `json:"signatory"`
}

// CreateBankAccountRequest request has the parameters needed to submit a bank account creation
// request to Lob.
type CreateBankAccountRequest struct {
	RoutingNumber    string `json:"routing_number"`
	AccountNumber    string `json:"account_number"`
	BankAddressID    string `json:"bank_address"`
	AccountAddressID string `json:"account_address"`
	Signatory        string `json:"signatory"`
}

// CreateBankAccount creates a new bank account in Lob's system.
func (lob *Lob) CreateBankAccount(account *CreateBankAccountRequest) (*BankAccount, error) {
	resp := new(BankAccount)
	return resp, Metrics.CreateBankAccount.Call(func() error {
		return lob.Post("bank_accounts/", json2form(*account), resp)
	})
}

// GetBankAccount gets information on a bank account.
func (lob *Lob) GetBankAccount(id string) (*BankAccount, error) {
	resp := new(BankAccount)
	return resp, Metrics.GetBankAccount.Call(func() error {
		return lob.Get("bank_accounts/"+id, nil, resp)
	})
}

// ListBankAccountsResponse gives the results for listing all addresses for our account.
type ListBankAccountsResponse struct {
	Data        []BankAccount `json:"data"`
	Object      string        `json:"object"`
	NextURL     string        `json:"next_url"`
	PreviousURL string        `json:"next_url"`
	Count       int           `json:"count"`
}

// ListBankAccounts lists all addresses on this account, paginated.
func (lob *Lob) ListBankAccounts(count int, offset int) (*ListBankAccountsResponse, error) {
	if count <= 0 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}

	resp := new(ListBankAccountsResponse)
	return resp, Metrics.ListBankAccounts.Call(func() error {
		return lob.Get("bank_accounts", map[string]string{
			"count":  strconv.Itoa(count),
			"offset": strconv.Itoa(offset),
		}, resp)
	})
}
