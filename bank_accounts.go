package lob

import "strconv"

// BankAccount represents a bank account in lob's system.
type BankAccount struct {
	Error         *Error            `json:"error"`
	AccountNumber string            `json:"account_number"`
	BankName      string            `json:"bank_name"`
	DateCreated   string            `json:"date_created"`
	DateModified  string            `json:"date_modified"`
	Description   *string           `json:"description"`
	ID            string            `json:"id"`
	Metadata      map[string]string `json:"metadata"`
	Object        string            `json:"object"`
	RoutingNumber string            `json:"routing_number"`
	Signatory     string            `json:"signatory"`
	Verified      bool              `json:"verified"`
}

// CreateBankAccountRequest request has the parameters needed to submit a bank account creation
// request to Lob.
type CreateBankAccountRequest struct {
	Description   *string           `json:"description"`
	RoutingNumber string            `json:"routing_number"`
	AccountNumber string            `json:"account_number"`
	Signatory     string            `json:"signatory"`
	Metadata      map[string]string `json:"metadata"`
}

// CreateBankAccount creates a new bank account in Lob's system.
func (l *lob) CreateBankAccount(account *CreateBankAccountRequest) (*BankAccount, error) {
	resp := new(BankAccount)
	if err := l.post("bank_accounts/", json2form(*account), resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetBankAccount gets information on a bank account.
func (l *lob) GetBankAccount(id string) (*BankAccount, error) {
	resp := new(BankAccount)
	if err := l.get("bank_accounts/"+id, nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
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
func (l *lob) ListBankAccounts(count int, offset int) (*ListBankAccountsResponse, error) {
	if count <= 0 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}

	resp := new(ListBankAccountsResponse)
	if err := l.get("bank_accounts", map[string]string{
		"limit":  strconv.Itoa(count),
		"offset": strconv.Itoa(offset),
	}, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
