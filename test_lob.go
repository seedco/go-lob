package lob

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/pborman/uuid"
)

var Non200Error = errors.New("Non-200 Status code returned")

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

func (t *fakeLob) ListChecks(count int) (*ListChecksResponse, error) {
	if count <= 0 {
		count = 10
	}

	resp := new(ListChecksResponse)

	data := make([]Check, len(t.checks))

	for _, check := range t.checks {
		data = append(data, *check)
	}

	resp.Data = data
	resp.Count = count
	return resp, nil
}

// Addresses

func (t *fakeLob) CreateAddress(address *Address) (*Address, error) {
	var message string
	var status int
	if address.Name != nil && len(*address.Name) > 40 {
		message = "name length must be less than or equal to 40 characters long"
		status = 422
		address.Error = &Error{
			Message:    message,
			StatusCode: status,
		}
		return address, Non200Error
	}
	if len(address.AddressLine1) > 200 {
		message = "address_line1 length must be less than or equal to 200 characters long"
		status = 422
		address.Error = &Error{
			Message:    message,
			StatusCode: status,
		}
		return address, Non200Error
	}
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

func (t *fakeLob) ListAddresses(count int) (*ListAddressesResponse, error) {
	if count <= 0 {
		count = 10
	}

	resp := new(ListAddressesResponse)

	data := make([]Address, len(t.addresses))

	for _, address := range t.addresses {
		data = append(data, *address)
	}

	resp.Data = data
	resp.Count = count
	return resp, nil
}

func (t *fakeLob) VerifyUSAddress(address *Address) (*USAddressVerificationResponse, error) {
	resp := new(USAddressVerificationResponse)

	if address != nil {
		resp.Id = fmt.Sprintf("us_ver_%v", uuid.New())
		resp.PrimaryLine = address.AddressLine1
		resp.SecondaryLine = *address.AddressLine2
		resp.LastLine = fmt.Sprintf("%v %v %v-%v", address.AddressCity, address.AddressState, address.AddressZip, "0000")
		resp.Deliverability = "no_match"
		for _, a := range t.addresses {
			if a.AddressCity == address.AddressCity &&
				a.AddressCountry == address.AddressCountry &&
				a.AddressLine1 == address.AddressLine1 &&
				a.AddressLine2 == address.AddressLine2 &&
				a.AddressState == address.AddressState &&
				a.AddressZip == address.AddressZip {
				resp.Deliverability = "deliverable"
				break
			}
		}
	}

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

func (t *fakeLob) ListBankAccounts(count int) (*ListBankAccountsResponse, error) {
	if count <= 0 {
		count = 10
	}

	resp := new(ListBankAccountsResponse)

	data := make([]BankAccount, len(t.bankAccounts))

	for _, bankAccount := range t.bankAccounts {
		data = append(data, *bankAccount)
	}

	resp.Data = data
	resp.Count = count
	return resp, nil
}
