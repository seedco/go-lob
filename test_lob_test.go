package lob

import (
	"testing"

	"github.com/pborman/uuid"
)

func TestFakeLobChecks(t *testing.T) {
	lob := NewFakeLob()

	city := "Davis"
	country := "USA"
	address1 := "1234 Seasame St."
	zip := "95616"
	address := &Address{
		ID:             uuid.New(),
		AddressCity:    &city,
		AddressCountry: &country,
		AddressLine1:   address1,
		AddressZip:     &zip,
	}
	var err error
	if _, err = lob.CreateAddress(address); err != nil {
		t.Error("create address had an error")
	}

	bankAccountRequest := &CreateBankAccountRequest{
		AccountNumber: "1132234455",
		RoutingNumber: "00000000",
	}
	var bankAccount *BankAccount
	if bankAccount, err = lob.CreateBankAccount(bankAccountRequest); err != nil {
		t.Error("create address had an error")
	}

	var check *Check
	if check, err = lob.CreateCheck(&CreateCheckRequest{
		Amount:        100,
		BankAccountID: bankAccount.ID,
		ToAddressID:   address.ID,
	}); err != nil {
		t.Error("create check had an error")
	}

	var retrievedCheck *Check
	if retrievedCheck, err = lob.GetCheck(check.ID); err != nil {
		t.Error("get check had an error")
	}

	if retrievedCheck.Amount != check.Amount {
		t.Errorf("expected check amount to be %d, got %d", check.Amount, retrievedCheck.Amount)
	}
}
