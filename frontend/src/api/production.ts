import apiClient from './client';
import type { ApiResponse, PageResponse } from '../types';

export interface ProductionOrderRequest {
  product_id: number;
  quantity: number;
  priority: number;
  start_date: string;
  end_date: string;
  description?: string;
}

export interface ProductionOrderResponse {
  id: number;
  order_no: string;
  product_id: number;
  product?: {
    id: number;
    code: string;
    name: string;
    unit: string;
  };
  quantity: number;
  produced: number;
  priority: number;
  status: string;
  start_date: string;
  end_date: string;
  created_by: number;
  creator?: {
    id: number;
    username: string;
    real_name: string;
  };
  created_at: string;
  updated_at: string;
}

export interface ProductionOrderQueryParams {
  page?: number;
  page_size?: number;
  status?: string;
  product_id?: number;
}

/**
 * 获取生产工单列表
 */
export const getProductionOrders = (params?: ProductionOrderQueryParams) =>
  apiClient.get<ApiResponse<PageResponse<ProductionOrderResponse>>>('/production/orders', { params });

/**
 * 创建生产工单
 */
export const createProductionOrder = (data: ProductionOrderRequest) =>
  apiClient.post<ApiResponse<ProductionOrderResponse>>('/production/orders', data);

/**
 * 获取生产工单详情
 */
export const getProductionOrderById = (id: number) =>
  apiClient.get<ApiResponse<ProductionOrderResponse>>(`/production/orders/${id}`);

/**
 * 更新生产工单
 */
export const updateProductionOrder = (id: number, data: ProductionOrderRequest) =>
  apiClient.put<ApiResponse<ProductionOrderResponse>>(`/production/orders/${id}`, data);

/**
 * 删除生产工单
 */
export const deleteProductionOrder = (id: number) =>
  apiClient.delete<ApiResponse<any>>(`/production/orders/${id}`);

/**
 * 获取生产工单状态列表
 */
export const getProductionOrderStatuses = () =>
  apiClient.get<ApiResponse<string[]>>('/production/orders/statuses');