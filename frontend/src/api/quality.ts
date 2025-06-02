import apiClient from './client';
import type { ApiResponse, PageResponse } from '../types';

export interface QualityStandardRequest {
  product_id: number;
  name: string;
  type: string;
  criteria: string;
  min_value?: number;
  max_value?: number;
  unit?: string;
  is_active: boolean;
  description?: string;
}

export interface QualityStandardResponse {
  id: number;
  product_id: number;
  product?: {
    id: number;
    code: string;
    name: string;
  };
  name: string;
  type: string;
  criteria: string;
  min_value?: number;
  max_value?: number;
  unit?: string;
  is_active: boolean;
  description?: string;
  created_at: string;
  updated_at: string;
}

export interface QualityInspectionRequest {
  production_order_id: number;
  quality_standard_id: number;
  inspector_id: number;
  batch_no?: string;
  sample_size: number;
  measured_value?: number;
  result: 'pass' | 'fail';
  notes?: string;
}

export interface QualityInspectionResponse {
  id: number;
  production_order_id: number;
  production_order?: {
    id: number;
    order_no: string;
    product?: {
      id: number;
      code: string;
      name: string;
    };
  };
  quality_standard_id: number;
  quality_standard?: QualityStandardResponse;
  inspector_id: number;
  inspector?: {
    id: number;
    username: string;
    real_name: string;
  };
  batch_no?: string;
  sample_size: number;
  measured_value?: number;
  result: string;
  notes?: string;
  inspected_at: string;
  created_at: string;
  updated_at: string;
}

export interface QualityStatistics {
  total_inspections: number;
  pass_count: number;
  fail_count: number;
  pass_rate: number;
  recent_inspections: QualityInspectionResponse[];
}

/**
 * 获取质量标准列表
 */
export const getQualityStandards = (params?: any) =>
  apiClient.get<ApiResponse<PageResponse<QualityStandardResponse>>>('/api/quality/standards', { params });

/**
 * 创建质量标准
 */
export const createQualityStandard = (data: QualityStandardRequest) =>
  apiClient.post<ApiResponse<QualityStandardResponse>>('/api/quality/standards', data);

/**
 * 获取质量检测记录列表
 */
export const getQualityInspections = (params?: any) =>
  apiClient.get<ApiResponse<PageResponse<QualityInspectionResponse>>>('/api/quality/inspections', { params });

/**
 * 创建质量检测记录
 */
export const createQualityInspection = (data: QualityInspectionRequest) =>
  apiClient.post<ApiResponse<QualityInspectionResponse>>('/api/quality/inspections', data);

/**
 * 获取质量统计数据
 */
export const getQualityStatistics = (params?: any) =>
  apiClient.get<ApiResponse<QualityStatistics>>('/api/quality/statistics', { params });

/**
 * 获取质量标准类型
 */
export const getQualityStandardTypes = () =>
  apiClient.get<ApiResponse<string[]>>('/api/quality/standards/types');