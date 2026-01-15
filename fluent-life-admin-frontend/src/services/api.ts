import axios from 'axios';

const API_BASE_URL = 'http://localhost:8082/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 请求拦截器
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('admin_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 管理员API
export const adminAPI = {
  login: async (username: string, password: string) => {
    const response = await api.post('/admin/login', { username, password });
    if (response.data.code === 0 && response.data.data?.token) {
      localStorage.setItem('admin_token', response.data.data.token);
    }
    return response.data;
  },

  // 用户管理
  getUsers: async (params: { page?: number; page_size?: number; keyword?: string }) => {
    const response = await api.get('/admin/users', { params });
    return response.data;
  },
  getUser: async (id: string) => {
    const response = await api.get(`/admin/users/${id}`);
    return response.data;
  },
  deleteUser: async (id: string) => {
    const response = await api.delete(`/admin/users/${id}`);
    return response.data;
  },

  // 帖子管理
  getPosts: async (params: { page?: number; page_size?: number; keyword?: string }) => {
    const response = await api.get('/admin/posts', { params });
    return response.data;
  },
  getPost: async (id: string) => {
    const response = await api.get(`/admin/posts/${id}`);
    return response.data;
  },
  deletePostsBatch: async (ids: string[]) => {
    const response = await api.post('/admin/posts/delete-batch', { ids });
    return response.data;
  },

  // 房间管理
  getRooms: async (params: { page?: number; page_size?: number; keyword?: string; is_active?: string }) => {
    const response = await api.get('/admin/rooms', { params });
    return response.data;
  },
  getRoom: async (id: string) => {
    const response = await api.get(`/admin/rooms/${id}`);
    return response.data;
  },
  deleteRoom: async (id: string) => {
    const response = await api.delete(`/admin/rooms/${id}`);
    return response.data;
  },
  deleteRoomsBatch: async (ids: string[]) => {
    const response = await api.post('/admin/rooms/delete-batch', { ids });
    return response.data;
  },
  toggleRoom: async (id: string) => {
    const response = await api.patch(`/admin/rooms/${id}/toggle`);
    return response.data;
  },

  // 训练统计
  getTrainingStats: async () => {
    const response = await api.get('/admin/training/stats');
    return response.data;
  },
  getTrainingRecords: async (params: { page?: number; page_size?: number; type?: string; user_id?: string }) => {
    const response = await api.get('/admin/training/records', { params });
    return response.data;
  },

  // 绕口令管理
  getTongueTwisters: async (params: { page?: number; page_size?: number; keyword?: string; level?: string; is_active?: string }) => {
    const response = await api.get('/admin/tongue-twisters', { params });
    return response.data;
  },
  getTongueTwister: async (id: string) => {
    const response = await api.get(`/admin/tongue-twisters/${id}`);
    return response.data;
  },
  createTongueTwister: async (data: any) => {
    const response = await api.post('/admin/tongue-twisters', data);
    return response.data;
  },
  updateTongueTwister: async (id: string, data: any) => {
    const response = await api.put(`/admin/tongue-twisters/${id}`, data);
    return response.data;
  },
  deleteTongueTwistersBatch: async (ids: string[]) => {
    const response = await api.post('/admin/tongue-twisters/delete-batch', { ids });
    return response.data;
  },

  // 每日朗诵文案管理
  getDailyExpressions: async (params: { page?: number; page_size?: number; keyword?: string; is_active?: string }) => {
    const response = await api.get('/admin/daily-expressions', { params });
    return response.data;
  },
  getDailyExpression: async (id: string) => {
    const response = await api.get(`/admin/daily-expressions/${id}`);
    return response.data;
  },
  createDailyExpression: async (data: any) => {
    const response = await api.post('/admin/daily-expressions', data);
    return response.data;
  },
  updateDailyExpression: async (id: string, data: any) => {
    const response = await api.put(`/admin/daily-expressions/${id}`, data);
    return response.data;
  },
  deleteDailyExpressionsBatch: async (ids: string[]) => {
    const response = await api.post('/admin/daily-expressions/delete-batch', { ids });
    return response.data;
  },
};

export default api;

