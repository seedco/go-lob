# go-lob
[Lob.com](http://lob.com) API client in Go.

## Install

Install with

```sh
go get github.com/seedco/go-lob
```

## Use

Use by creating a `Lob` struct with `NewLob`, and calling the methods it offers.

```go
// fill in your API key and user agent here
lob := lob.NewLob(lob.BaseAPI, apiKey, userAgent)

testAddress := &Address{
  Name:           "Lobster Test",
  Email:          "lobtest@example.com",
  Phone:          "5555555555",
  AddressLine1:   "1005 W Burnside St", // Powell's City of Books, the best book store in the world.
  AddressCity:    "Portland",
  AddressState:   "OR",
  AddressZip:     "97209",
  AddressCountry: "US",
}

verify, err := lob.VerifyAddress(testAddress)
// ...
```

## Test

You can run the tests if you set the `TEST_LOB_API_KEY` environment variable, i.e.,

```sh
TEST_LOB_API_KEY=test_yourtestkeyhere go test .
```

## License

Licensed under the MIT license. See [LICENSE](LICENSE) for more details.
