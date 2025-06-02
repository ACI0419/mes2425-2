import apiClient from './client';
import type { ApiResponse } from '../types';

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user_info: {
    id: number;
    username: string;
    real_name: string;
    email: string;
    phone?: string;
    role: string;
    status: number;
    created_at: string;
    updated_at: string;
  };
}

export interface RegisterRequest {
  username: string;
  password: string;
  real_name: string;
  email: string;
  phone?: string;
  role?: string;
}

/**
 * 用户登录
 * @param data 登录请求数据
 */
export const login = (data: LoginRequest) => 
  apiClient.post<ApiResponse<LoginResponse>>('/users/login', data);

/**
 * 用户注册
 * @param data 注册请求数据
 */
export const register = (data: RegisterRequest) => 
  apiClient.post<ApiResponse<any>>('/users/register', data);

/**
 * 获取用户信息
 */
export const getUserProfile = () => 
  apiClient.get<ApiResponse<any>>('/users/profile');

/**
 * 更新用户信息
 * @param data 更新请求数据
 */
export const updateProfile = (data: any) => 
  apiClient.put<ApiResponse<any>>('/users/profile', data);

/**
 * 修改密码
 * @param data 修改密码请求数据
 */
export const changePassword = (data: any) => 
  apiClient.put<ApiResponse<any>>('/users/password', data);

/**
 * 刷新token
 * @param data 刷新token请求数据
 */
export const refreshToken = (data: any) => 
  apiClient.post<ApiResponse<any>>('/users/refresh', data);