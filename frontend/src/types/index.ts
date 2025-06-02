// 通用响应类型
export interface ApiResponse<T = any> {
  token: string | null;
  user_info: any;
  code: number;
  message: string;
  data: T;
}

// 分页响应类型
export interface PageResponse<T = any> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

// 用户类型
export interface User {
  id: number;
  username: string;
  real_name: string;
  email: string;
  phone?: string;
  role: string;
  status: number;
  created_at: string;
  updated_at: string;
}

// 产品类型
export interface Product {
  id: number;
  code: string;
  name: string;
  description?: string;
  unit: string;
  price: number;
  status: number;
  created_at: string;
  updated_at: string;
}

// 生产订单类型
export interface ProductionOrder {
  id: number;
  order_no: string;
  product_id: number;
  product?: Product;
  quantity: number;
  produced: number;
  priority: number;
  status: number;
  start_date: string;
  end_date: string;
  created_by: number;
  creator?: User;
  created_at: string;
  updated_at: string;
}

// 设备类型
export interface Equipment {
  id: number;
  code: string;
  name: string;
  type: string;
  model?: string;
  manufacturer?: string;
  location?: string;
  status: string;
  purchase_date?: string;
  warranty_date?: string;
  description?: string;
  created_at: string;
  updated_at: string;
}