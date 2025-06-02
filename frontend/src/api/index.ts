// 如果这个文件存在，更新API基础URL
const API_BASE_URL = '/api/v1'; // 使用相对路径

// 或者在axios配置中
import axios from 'axios';

const api = axios.create({
  baseURL: '/api/v1', // 使用相对路径，通过Nginx代理
  timeout: 10000,
});

export default api;