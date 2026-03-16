import { fetchApi } from '../api/api.js';

export const adminService = {
    // Dashboard Stats
    async getDashboard() {
        return fetchApi('/admin/dashboard');
    },

    // Users
    async getUsers() {
        return fetchApi('/admin/users');
    },
    async updateUserRole(id, roleData) {
        return fetchApi(`/admin/users/${id}/role`, { method: 'PATCH', body: roleData });
    },
    async deleteUser(id) {
        return fetchApi(`/admin/users/${id}`, { method: 'DELETE' });
    },
    async restoreUser(id) {
        return fetchApi(`/admin/users/${id}/restore`, { method: 'PATCH' });
    },
    async banUser(id, reasonData) {
        return fetchApi(`/admin/users/${id}/ban`, { method: 'PATCH', body: reasonData });
    },
    async unbanUser(id) {
        return fetchApi(`/admin/users/${id}/unban`, { method: 'PATCH' });
    },

    // Transactions
    async getTransactions() {
        return fetchApi('/admin/transactions');
    },

    // Flights
    async getFlights() {
        return fetchApi('/admin/flights');
    },
    async createFlight(data) {
        return fetchApi('/admin/flights', { method: 'POST', body: data });
    },
    async updateFlight(id, data) {
        return fetchApi(`/admin/flights/${id}`, { method: 'PUT', body: data });
    },
    async deleteFlight(id) {
        return fetchApi(`/admin/flights/${id}`, { method: 'DELETE' });
    },

    // Airlines
    async getAirlinesPublic() {
        return fetchApi('/airlines');
    },
    async getAirlines() {
        return fetchApi('/admin/airlines');
    },
    async createAirline(data) {
        return fetchApi('/admin/airlines', { method: 'POST', body: data });
    },
    async updateAirline(id, data) {
        return fetchApi(`/admin/airlines/${id}`, { method: 'PUT', body: data });
    },
    async deleteAirline(id) {
        return fetchApi(`/admin/airlines/${id}`, { method: 'DELETE' });
    },

    // Airports
    async getAirportsPublic() {
        return fetchApi('/airports');
    },
    async getAirports() {
        return fetchApi('/admin/airports');
    },
    async createAirport(data) {
        return fetchApi('/admin/airports', { method: 'POST', body: data });
    },
    async updateAirport(id, data) {
        return fetchApi(`/admin/airports/${id}`, { method: 'PUT', body: data });
    },
    async deleteAirport(id) {
        return fetchApi(`/admin/airports/${id}`, { method: 'DELETE' });
    },

    // Promos
    async getPromos() {
        return fetchApi('/admin/promos');
    },
    async createPromo(data) {
        return fetchApi('/admin/promos', { method: 'POST', body: data });
    },
    async updatePromo(id, data) {
        return fetchApi(`/admin/promos/${id}`, { method: 'PUT', body: data });
    },
    async deletePromo(id) {
        return fetchApi(`/admin/promos/${id}`, { method: 'DELETE' });
    }
};

export const systemService = {
    async seedData(appPassword) {
        return fetchApi('/system/seed', {
            method: 'POST',
            headers: {
                'X-App-Password': appPassword
            }
        });
    }
}
