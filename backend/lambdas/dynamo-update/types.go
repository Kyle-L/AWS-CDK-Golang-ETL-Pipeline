package main

type Item struct {
	Id                      float64 `json:"id"`
	AccountNumber           string  `json:"accountNumber"`
	CustomerId              string  `json:"customerId"`
	CreditLimit             float64 `json:"creditLimit"`
	Availablemoney          float64 `json:"availableMoney"`
	TransactionDateTime     string  `json:"transactionDateTime"`
	TransactionAmount       float64 `json:"transactionAmount"`
	MerchantName            string  `json:"merchantName"`
	AcqCountry              string  `json:"acqCountry"`
	MerchantCountryCode     string  `json:"merchantCountryCode"`
	PosEntryMode            string  `json:"posEntryMode"`
	PosConditionCode        float64 `json:"posConditionCode"`
	MerchantCategoryCode    string  `json:"merchantCategoryCode"`
	CurrentExpDate          string  `json:"currentExpDate"`
	AccountOpenDate         string  `json:"accountOpenDate"`
	DateOfLastAddressChange string  `json:"dateOfLastAddressChange"`
	CardCVV                 float64 `json:"cardCVV"`
	CardLast4Digits         float64 `json:"cardLast4Digits"`
	TransactionType         string  `json:"transactionType"`
	CurrentBalance          float64 `json:"currentBalance"`
	CardPresent             string  `json:"cardPresent"`
	IsFraud                 string  `json:"isFraud"`
	CountryCode             string  `json:"countryCode"`
}

type Body struct {
	Message string `json:"message"`
}
