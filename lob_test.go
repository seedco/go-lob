package lob

import (
	"os"
	"testing"
)

const testUserAgent = "Test/1.0"

var testAPIKey = os.Getenv("TEST_LOB_API_KEY")

var testAddress = &Address{
	Name:           nullString("Lobster Test"),
	Email:          nullString("lobtest@example.com"),
	Phone:          nullString("5555555555"),
	AddressLine1:   "1005 W Burnside St", // Powell's City of Books, the best book store in the world.
	AddressCity:    nullString("Portland"),
	AddressState:   nullString("OR"),
	AddressZip:     nullString("97209"),
	AddressCountry: nullString("US"),
}

func nullString(s string) *string {
	return &s
}

func TestAddresses(t *testing.T) {
	lob := NewLob(BaseAPI, testAPIKey, testUserAgent)

	verify, err := lob.VerifyUSAddress(testAddress)
	if err != nil {
		t.Errorf("Error verifying address: %s", err.Error())
	}
	t.Logf("Verification = %+v", verify)

	address, err := lob.CreateAddress(testAddress)
	if err != nil {
		t.Fatalf("Could not create address: %s", err.Error())
	}

	address, err = lob.GetAddress(address.ID)
	if err != nil {
		t.Errorf("Could not get address: %s", err.Error())
	}

	addresses, err := lob.ListAddresses(-1)
	if err != nil {
		t.Errorf("Could not list addresses: %s", err.Error())
	}
	t.Logf("Address list = %+v", addresses)

	err = lob.DeleteAddress(address.ID)
	if err != nil {
		t.Errorf("Error deleting address: %s", err.Error())
	}
}
func TestAddressError(t *testing.T) {
	lob := NewLob(BaseAPI, testAPIKey, testUserAgent)

	address, err := lob.CreateAddress(&Address{
		Name:           nullString("Name that is way too long to be printed on the check so that it will error."),
		Email:          nullString("lobtest@example.com"),
		Phone:          nullString("5555555555"),
		AddressLine1:   "1005 W Burnside St", // Powell's City of Books, the best book store in the world.
		AddressCity:    nullString("Portland"),
		AddressState:   nullString("OR"),
		AddressZip:     nullString("97209"),
		AddressCountry: nullString("US"),
	})
	if err == nil {
		t.Error("error should not have been nil")
	}
	if address.Error.Message == "" {
		t.Error("Expected human readable error message")
	}
	if address.Error.StatusCode != 422 {
		t.Error("Expected status code to be 422")
	}
}

func TestBankAccounts(t *testing.T) {
	lob := NewLob(BaseAPI, testAPIKey, testUserAgent)

	address, err := lob.CreateAddress(testAddress)
	if err != nil {
		t.Fatalf("Could not create address: %s", err.Error())
	}

	bankAccount, err := lob.CreateBankAccount(&CreateBankAccountRequest{
		RoutingNumber: "255077370",
		AccountNumber: "1234",
		Signatory:     "Lobster Test",
		AccountType:   "company",
	})

	if err != nil {
		t.Fatalf("Could not create bank account: %s", err.Error())
	}
	t.Logf("Bank account = %+v", bankAccount)

	bankAccount, err = lob.GetBankAccount(bankAccount.ID)
	if err != nil {
		t.Errorf("Error retrieving bank account")
	}
	t.Logf("Bank account = %+v", bankAccount)

	resp, err := lob.ListBankAccounts(-1)
	if err != nil {
		t.Errorf("Could not list bank accounts: %s", err.Error())
	}
	t.Logf("Bank accounts = %+v", resp)

	if err != nil {
		t.Fatalf("Could not create bank account: %s", err.Error())
	}

	err = lob.DeleteAddress(address.ID)
	if err != nil {
		t.Errorf("Error deleting address: %s", err.Error())
	}
}

func TestChecks(t *testing.T) {
	lob := NewLob(BaseAPI, testAPIKey, testUserAgent)

	address, err := lob.CreateAddress(testAddress)
	if err != nil {
		t.Fatalf("Could not create address: %s", err.Error())
	}

	bankAccount, err := lob.CreateBankAccount(&CreateBankAccountRequest{
		RoutingNumber: "255077370",
		AccountNumber: "1234",
		Signatory:     "Lobster Test",
		AccountType:   "company",
	})

	if err != nil {
		t.Fatalf("Could not create bank account: %s", err.Error())
	}

	check, err := lob.CreateCheck(&CreateCheckRequest{
		CheckNumber:   nullString("12345"),
		BankAccountID: bankAccount.ID,
		FromAddressID: address.ID,
		ToAddressID:   address.ID,
		Amount:        987.65,
		Message:       nullString("Some message"),
		Memo:          nullString("A memo"),
	})

	if err != nil {
		t.Fatalf("Could not create check: %s", err.Error())
	}
	t.Logf("Check = %+v", check)

	_, err = lob.GetCheck(check.ID)
	if err != nil {
		t.Errorf("Could not get check: %s", err.Error())
	}

	resp, err := lob.ListChecks(-1)
	if err != nil {
		t.Errorf("Could not list checks: %s", err.Error())
	}
	t.Logf("List checks = %+v", resp)

	err = lob.DeleteAddress(address.ID)
	if err != nil {
		t.Errorf("Error deleting address: %s", err.Error())
	}
}

func TestGetStates(t *testing.T) {
	lob := NewLob(BaseAPI, testAPIKey, testUserAgent)
	list, err := lob.GetStates()
	if err != nil {
		t.Fatalf("Error retrieving state list: %s", err.Error())
	}
	if len(list.Data) < 50 || len(list.Data) > 80 {
		t.Errorf("Expected at least 50 US states, got %d", len(list.Data))
	}
}

func TestGetCountries(t *testing.T) {
	lob := NewLob(BaseAPI, testAPIKey, testUserAgent)
	list, err := lob.GetCountries()
	if err != nil {
		t.Fatalf("Error retrieving countries list: %s", err.Error())
	}
	if len(list.Data) < 200 || len(list.Data) > 400 {
		t.Errorf("Expected at least 200 countries, got %d", len(list.Data))
	}
}
