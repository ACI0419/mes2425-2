import apiClient from './client';
import type { ApiResponse, PageResponse } from '../types';

export interface MaterialRequest {
  code: string;
  name: string;
  type: string;
  unit: string;
  price: number;
  current_stock: number;
  min_stock: number;
  max_stock: number;
  supplier?: string;
  description?: string;
}

export interface MaterialResponse {
  id: number;
  code: string;
  name: string;
  type: string;
  unit: string;
  price: number;
  current_stock: number;
  min_stock: number;
  max_stock: number;
  supplier?: string;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface MaterialTransactionRequest {
  material_id: number;
  type: 'in' | 'out';
  quantity: number;
  unit_price: number;
  reference_no?: string;
  notes?: string;
}

export interface MaterialTransactionResponse {
  id: number;
  material_id: number;
  material?: MaterialResponse;
  type: string;
  quantity: number;
  unit_price: number;
  total_amount: number;
  reference_no?: string;
  notes?: string;
  created_by: number;
  creator?: {
    id: number;
    username: string;
    real_name: string;
  };
  created_at: string;
}

export interface MaterialQueryParams {
  page?: number;
  page_size?: number;
  type?: string;
  keyword?: string;
}

export interface MaterialTransactionQueryParams {
  page?: number;
  page_size?: number;
  material_id?: number;
  type?: string;
}

/**
 * 获取物料列表
 */
export const getMaterials = (params?: MaterialQueryParams) =>
  apiClient.get<ApiResponse<PageResponse<MaterialResponse>>>('/api/materials', { params });

/**
 * 创建物料
 */
export const createMaterial = (data: MaterialRequest) =>
  apiClient.post<ApiResponse<MaterialResponse>>('/api/materials', data);

/**
 * 获取物料详情
 */
export const getMaterialById = (id: number) =>
  apiClient.get<ApiResponse<MaterialResponse>>(`/api/materials/${id}`);

/**
 * 更新物料
 */
export const updateMaterial = (id: number, data: MaterialRequest) =>
  apiClient.put<ApiResponse<MaterialResponse>>(`/api/materials/${id}`, data);

/**
 * 删除物料
 */
export const deleteMaterial = (id: number) =>
  apiClient.delete<ApiResponse<any>>(`/api/materials/${id}`);

/**
 * 获取物料类型列表
 */
export const getMaterialTypes = () =>
  apiClient.get<ApiResponse<string[]>>('/api/materials/types');

/**
 * 获取低库存物料
 */
export const getLowStockMaterials = () =>
  apiClient.get<ApiResponse<MaterialResponse[]>>('/api/materials/low-stock');

/**
 * 获取物料交易记录
 */
export const getMaterialTransactions = (params?: MaterialTransactionQueryParams) =>
  apiClient.get<ApiResponse<PageResponse<MaterialTransactionResponse>>>('/api/materials/transactions', { params });

/**
 * 创建物料交易记录
 */
export const createMaterialTransaction = (data: MaterialTransactionRequest) =>
  apiClient.post<ApiResponse<MaterialTransactionResponse>>('/api/materials/transactions', data);