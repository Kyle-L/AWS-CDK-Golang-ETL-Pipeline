export type TransactionQueryResponse = {
    items: Transaction[];
    count: number;
    paginationToken: string;
  }

export type Transaction = {
    deleted: boolean;
    id: number;
    accountNumber: string;
    customerId: string;
    creditLimit: number;
    availableMoney: number;
    transactionDateTime: string;
    transactionAmount: number;
    merchantName: string;
    acqCountry: string;
    merchantCountryCode: string;
    posEntryMode: string;
    posConditionCode: number;
    merchantCategoryCode: string;
    currentExpDate: string;
    accountOpenDate: string;
    dateOfLastAddressChange: string;
    cardCVV: number;
    cardLast4Digits: number;
    transactionType: string;
    currentBalance: number;
    cardPresent: string;
    isFraud: string;
    countryCode: string;
};