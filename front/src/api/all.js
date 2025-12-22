import request from '@/utils/request'

// Auth
export const login = (data) => request.post('/auth/login', data)
export const register = (data) => request.post('/users', data)

// User
export const getUser = (id) => request.get(`/users/${id}`)
export const addEmail = (userId, data) => request.post(`/users/${userId}/emails`, data)
export const deleteEmail = (userId, emailId) => request.delete(`/users/${userId}/emails/${emailId}`)
export const addPhone = (userId, data) => request.post(`/users/${userId}/phones`, data)
export const deletePhone = (userId, phoneId) => request.delete(`/users/${userId}/phones/${phoneId}`)

// Providers
export const getProviders = (params) => request.get('/providers', { params })
export const linkProvider = (userId, data) => request.post(`/users/${userId}/providers`, data)

// Appointments
export const getAppointments = (userId) => request.get(`/users/${userId}/appointments`)
export const createAppointment = (data) => request.post('/appointments', data)
export const cancelAppointment = (id, data) => request.put(`/appointments/${id}/cancel`, data)

// Challenges
export const getPopularChallenges = (limit) => request.get('/stats/popular_challenges', { params: { limit } })
export const joinChallenge = (id, data) => request.post(`/challenges/${id}/join`, data)
export const checkinChallenge = (id, data) => request.post(`/challenges/${id}/checkin`, data)
export const getDailyProgress = (id, userId) => request.get(`/challenges/${id}/progress`, { params: { user_id: userId } })
export const createChallenge = (data) => request.post('/challenges', data)
export const inviteToChallenge = (id, data) => request.post(`/challenges/${id}/invite`, data)

// Reports
export const generateReport = (data) => request.post('/reports/generate', data)

