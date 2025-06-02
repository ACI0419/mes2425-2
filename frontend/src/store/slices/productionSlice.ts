import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { ProductionOrder } from '../../types';

interface ProductionState {
  orders: ProductionOrder[];
  loading: boolean;
  error: string | null;
}

const initialState: ProductionState = {
  orders: [],
  loading: false,
  error: null,
};

/**
 * 获取生产订单列表
 */
export const getProductionOrdersAsync = createAsyncThunk(
  'production/getOrders',
  async () => {
    // TODO: 实现API调用
    return [];
  }
);

const productionSlice = createSlice({
  name: 'production',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(getProductionOrdersAsync.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(getProductionOrdersAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.orders = action.payload;
      })
      .addCase(getProductionOrdersAsync.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || '获取生产订单失败';
      });
  },
});

export const { clearError } = productionSlice.actions;
export default productionSlice.reducer;