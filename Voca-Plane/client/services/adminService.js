import { fetchApi } from '../api/api.js';

export const adminService = {
    // Dashboard Stats
    async getDashboard() {
        return fetchApi('/admin/dashboard');
    },

    // Users
    async getUsers(page = 1, limit = 10) {
        return fetchApi(`/admin/users?page=${page}&limit=${limit}`);
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
    async getTransactions(page = 1, limit = 10) {
        return fetchApi(`/admin/transactions?page=${page}&limit=${limit}`);
    },

    // Flights
    async getFlights(page = 1, limit = 10) {
        return fetchApi(`/admin/flights?page=${page}&limit=${limit}`);
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
    async getAirlines(page = 1, limit = 10) {
        return fetchApi(`/admin/airlines?page=${page}&limit=${limit}`);
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
    async getAirports(page = 1, limit = 10) {
        return fetchApi(`/admin/airports?page=${page}&limit=${limit}`);
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
    async getPromos(page = 1, limit = 10) {
        return fetchApi(`/admin/promos?page=${page}&limit=${limit}`);
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
