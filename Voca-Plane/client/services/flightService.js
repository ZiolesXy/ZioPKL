import { fetchApi } from '../api/api.js';

export const flightService = {
    async getAllFlights() {
        return fetchApi('/flights');
    },

    async getAllFlightsFull() {
        return fetchApi('/flights/all');
    },

    async searchFlights(params) {
        const query = new URLSearchParams(params).toString();
        return fetchApi(`/flights/search?${query}`);
    },

    async getFlightById(id) {
        return fetchApi(`/flights/${id}`);
    },

    async getFlightSeats(id) {
        return fetchApi(`/flights/${id}/seats`);
    }
};
