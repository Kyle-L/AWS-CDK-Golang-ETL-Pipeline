import axios from 'axios';
import { Transaction, TransactionQueryResponse } from '../types/types';

function getTransactions(filter: { [key: string]: any }, callback: (transactions: TransactionQueryResponse) => void) {
    axios.get(`${process.env.API_ENDPOINT}/transactions`, { params: filter })
        .then(response => {
            callback(response.data);
        });
}

function putTransaction(transaction: Transaction, callback: (transaction: Transaction) => void) {
    axios.put(`${process.env.API_ENDPOINT}/transactions/${transaction.id}`, transaction)
        .then(response => {
            callback(response.data);
        });
}

export { getTransactions, putTransaction };