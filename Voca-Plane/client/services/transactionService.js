import { fetchApi } from '../api/api.js';

export const transactionService = {
    async createTransaction(data) {
        return fetchApi('/transactions', {
            method: 'POST',
            body: data
        });
    },

    async getUserTransactions() {
        return fetchApi('/transactions');
    },

    async getTransactionByCode(code) {
        return fetchApi(`/transactions/${code}`);
    },

    async cancelTransaction(code) {
        return fetchApi(`/transactions/${code}`, {
            method: 'DELETE'
        });
    }
};
