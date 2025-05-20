import axios from "axios";

const api = axios.create({
    baseURL: 'http://localhost:8080/api',
    withCredentials: true
})

export const register = (userData) => api.post('/register', userData);
export const login = (credentials) => 
  api.post('/login', credentials, { 
    withCredentials: true 
  });

export const getProfile = () => api.get('/user/me');
//export const updateUserProfile = () => api.put('/user/me')

export const getWallet = () => api.get('/wallets');
export const deposit = (amount) => api.post('/wallets/deposit', { amount });

export const getMatches = () => api.get('/bet/matches');
export const placeBet = (betData) => api.post('/bet/placeBet', betData);

export const adminLogin = (credentials) => api.post('/admin/login', credentials);