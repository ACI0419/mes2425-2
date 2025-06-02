import apiClient from './client';
import type { ApiResponse, PageResponse, Equipment } from '../types';

/**
 * 设备查询参数接口
 */
export interface EquipmentQueryParams {
  page?: number;
  page_size?: number;
  keyword?: string;
  type?: string;
  status?: string;
}

/**
 * 设备创建/更新请求接口
 */
export interface EquipmentRequest {
  code: string;
  name: string;
  type: string;
  status: string;
  model?: string;
  manufacturer?: string;
  location?: string;
  description?: string;
}

/**
 * 获取设备列表
 * @param params 查询参数
 */
export const getEquipments = (params?: EquipmentQueryParams) => 
  apiClient.get<ApiResponse<PageResponse<Equipment>>>('/equipments', { params });

/**
 * 根据ID获取设备详情
 * @param id 设备ID
 */
export const getEquipmentById = (id: number) => 
  apiClient.get<ApiResponse<Equipment>>(`/equipments/${id}`);

/**
 * 创建新设备
 * @param data 设备数据
 */
export const createEquipment = (data: EquipmentRequest) => 
  apiClient.post<ApiResponse<Equipment>>('/equipments', data);

/**
 * 更新设备信息
 * @param id 设备ID
 * @param data 更新的设备数据
 */
export const updateEquipment = (id: number, data: Partial<EquipmentRequest>) => 
  apiClient.put<ApiResponse<Equipment>>(`/equipments/${id}`, data);

/**
 * 删除设备
 * @param id 设备ID
 */
export const deleteEquipment = (id: number) => 
  apiClient.delete<ApiResponse<void>>(`/equipments/${id}`);

/**
 * 获取设备统计信息
 */
export const getEquipmentStatistics = () => 
  apiClient.get<ApiResponse<any>>('/equipments/statistics');

/**
 * 获取设备维护记录
 * @param equipmentId 设备ID
 * @param params 查询参数
 */
export const getMaintenanceRecords = (equipmentId: number, params?: any) => 
  apiClient.get<ApiResponse<PageResponse<any>>>(`/equipments/${equipmentId}/maintenance`, { params });

/**
 * 创建设备维护记录
 * @param equipmentId 设备ID
 * @param data 维护记录数据
 */
export const createMaintenanceRecord = (equipmentId: number, data: any) => 
  apiClient.post<ApiResponse<any>>(`/equipments/${equipmentId}/maintenance`, data);

/**
 * 更新设备维护记录
 * @param equipmentId 设备ID
 * @param recordId 维护记录ID
 * @param data 更新的维护记录数据
 */
export const updateMaintenanceRecord = (equipmentId: number, recordId: number, data: any) => 
  apiClient.put<ApiResponse<any>>(`/equipments/${equipmentId}/maintenance/${recordId}`, data);

/**
 * 删除设备维护记录
 * @param equipmentId 设备ID
 * @param recordId 维护记录ID
 */
export const deleteMaintenanceRecord = (equipmentId: number, recordId: number) => 
  apiClient.delete<ApiResponse<void>>(`/equipments/${equipmentId}/maintenance/${recordId}`);