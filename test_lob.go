package lob

import (
	"errors"
	"math/rand"
	"time"

	"github.com/pborman/uuid"
)

type fakeLob struct {
	checks       map[string]*Check
	addresses    map[string]*Address
	objectLists  map[string]*NamedObjectList
	bankAccounts map[string]*BankAccount
}

func NewFakeLob() *fakeLob {
	return &fakeLob{
		checks:       make(map[string]*Check),
		addresses:    make(map[string]*Address),
		objectLists:  make(map[string]*NamedObjectList),
		bankAccounts: make(map[string]*BankAccount),
	}
}

func (t *fakeLob) CreateCheck(request *CreateCheckRequest) (*Check, error) {

	bankAccount, ok := t.bankAccounts[request.BankAccountID]
	if !ok {
		return nil, errors.New("bank account not found")
	}

	address, ok := t.addresses[request.ToAddressID]
	if !ok {
		return nil, errors.New("address not found")
	}
	check := &Check{
		ID:                   uuid.New(),
		Amount:               request.Amount,
		BankAccount:          bankAccount,
		CheckNumber:          rand.Int(),
		ExpectedDeliveryDate: time.Now().Add(3 * 24 * time.Hour).Format("1/2/2006"),
		SendDate:             time.Now().Add(1 * 24 * time.Hour),
		To:                   address,
	}
	t.checks[check.ID] = check
	return check, nil
}

func (t *fakeLob) GetCheck(id string) (*Check, error) {
	check, ok := t.checks[id]
	if !ok {
		return nil, errors.New("no check found")
	}
	return check, nil
}

func (t *fakeLob) CancelCheck(id string) (*CancelCheckResponse, error) {
	delete(t.checks, id)
	return &CancelCheckResponse{
		ID:      id,
		Deleted: true,
	}, nil
}

func (t *fakeLob) ListChecks(count, offset int) (*ListChecksResponse, error) {
	if count <= 0 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}
	resp := new(ListChecksResponse)

	data := make([]Check, len(t.checks))

	for _, check := range t.checks {
		data = append(data, *check)
	}

	resp.Data = data[offset:count]
	resp.Count = count
	return resp, nil
}

// Addresses

func (t *fakeLob) CreateAddress(address *Address) (*Address, error) {
	if address.ID == "" {
		address.ID = uuid.New()
	}
	t.addresses[address.ID] = address
	return address, nil
}

func (t *fakeLob) GetAddress(id string) (*Address, error) {
	address, ok := t.addresses[id]
	if !ok {
		return nil, errors.New("address not found")
	}
	return address, nil
}

func (t *fakeLob) DeleteAddress(id string) error {
	delete(t.addresses, id)
	return nil
}

func (t *fakeLob) ListAddresses(count, offset int) (*ListAddressesResponse, error) {
	if count <= 0 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}
	resp := new(ListAddressesResponse)

	data := make([]Address, len(t.addresses))

	for _, address := range t.addresses {
		data = append(data, *address)
	}

	resp.Data = data[offset:count]
	resp.Count = count
	return resp, nil
}

func (t *fakeLob) VerifyAddress(address *Address) (*AddressVerificationResponse, error) {
	resp := new(AddressVerificationResponse)

	for _, a := range t.addresses {
		if a.AddressCity == address.AddressCity &&
			a.AddressCountry == address.AddressCountry &&
			a.AddressLine1 == address.AddressLine1 &&
			a.AddressLine2 == address.AddressLine2 &&
			a.AddressState == address.AddressState &&
			a.AddressZip == address.AddressZip {
			resp.Address = *a
			return resp, nil
		}
	}
	resp.Errors = []ErrorMessage{ErrorMessage{Message: "could not find address", StatusCode: 400}}
	return resp, nil
}

func (t *fakeLob) GetStates() (*NamedObjectList, error) {
	return &NamedObjectList{}, nil
}

func (t *fakeLob) GetCountries() (*NamedObjectList, error) {
	return &NamedObjectList{}, nil
}

func (t *fakeLob) CreateBankAccount(request *CreateBankAccountRequest) (*BankAccount, error) {
	bankAccount := &BankAccount{
		AccountNumber: request.AccountNumber,
		BankName:      "Fake Bank",
		DateCreated:   time.Now().Format("1/2/2006"),
		DateModified:  time.Now().Format("1/2/2006"),
		ID:            uuid.New(),
		Metadata:      request.Metadata,
		Object:        "",
		RoutingNumber: request.RoutingNumber,
		Signatory:     request.Signatory,
		Verified:      true,
	}
	t.bankAccounts[bankAccount.ID] = bankAccount
	return bankAccount, nil
}

func (t *fakeLob) GetBankAccount(id string) (*BankAccount, error) {
	bankAccount, ok := t.bankAccounts[id]
	if !ok {
		return nil, errors.New("bank account not found")
	}
	return bankAccount, nil
}

func (t *fakeLob) ListBankAccounts(count, offset int) (*ListBankAccountsResponse, error) {
	if count <= 0 {
		count = 10
	}
	if offset < 0 {
		offset = 0
	}
	resp := new(ListBankAccountsResponse)

	data := make([]BankAccount, len(t.bankAccounts))

	for _, bankAccount := range t.bankAccounts {
		data = append(data, *bankAccount)
	}

	resp.Data = data[offset:count]
	resp.Count = count
	return resp, nil
}
