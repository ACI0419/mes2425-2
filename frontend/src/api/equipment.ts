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
 * 维护记录请求接口
 */
export interface MaintenanceRecordRequest {
  equipment_id: number;
  maintainer_id: number;
  type: string; // preventive/corrective/emergency
  description: string;
  start_time: string;
  end_time?: string;
  cost?: number;
  parts_replaced?: string;
  result?: string;
  next_maintenance?: string;
  remark?: string;
}

/**
 * 维护记录响应接口
 */
export interface MaintenanceRecordResponse {
  id: number;
  equipment_id: number;
  equipment_code: string;
  equipment_name: string;
  maintainer_id: number;
  maintainer_name: string;
  type: string;
  description: string;
  start_time: string;
  end_time?: string;
  duration?: number; // 维护时长（分钟）
  cost?: number;
  parts_replaced?: string;
  result?: string;
  next_maintenance?: string;
  remark?: string;
  created_at: string;
}

/**
 * 维护记录查询参数接口
 */
export interface MaintenanceRecordQueryParams {
  page?: number;
  page_size?: number;
  equipment_id?: number;
  type?: string;
  maintainer_id?: number;
  start_date?: string;
  end_date?: string;
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
 * 获取维护记录列表
 * @param params 查询参数
 */
export const getMaintenanceRecords = (params?: MaintenanceRecordQueryParams) => 
  apiClient.get<ApiResponse<PageResponse<MaintenanceRecordResponse>>>('/api/equipments/maintenance-records', { params });

/**
 * 根据ID获取维护记录详情
 * @param id 维护记录ID
 */
export const getMaintenanceRecordById = (id: number) => 
  apiClient.get<ApiResponse<MaintenanceRecordResponse>>(`/api/equipments/maintenance-records/${id}`);

/**
 * 创建维护记录
 * @param data 维护记录数据
 */
export const createMaintenanceRecord = (data: MaintenanceRecordRequest) => 
  apiClient.post<ApiResponse<MaintenanceRecordResponse>>('/api/equipments/maintenance-records', data);

/**
 * 更新维护记录
 * @param id 维护记录ID
 * @param data 更新的维护记录数据
 */
export const updateMaintenanceRecord = (id: number, data: Partial<MaintenanceRecordRequest>) => 
  apiClient.put<ApiResponse<MaintenanceRecordResponse>>(`/api/equipments/maintenance-records/${id}`, data);

/**
 * 删除维护记录
 * @param id 维护记录ID
 */
export const deleteMaintenanceRecord = (id: number) => 
  apiClient.delete<ApiResponse<void>>(`/api/equipments/maintenance-records/${id}`);

/**
 * 获取即将到期的维护记录
 * @param days 天数
 */
export const getUpcomingMaintenance = (days: number = 7) => 
  apiClient.get<ApiResponse<MaintenanceRecordResponse[]>>(`/api/equipments/maintenance/upcoming?days=${days}`);