import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { login, getUserProfile } from '../../api/auth';
import type { User } from '../../types';

interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  loading: boolean;
}

const initialState: AuthState = {
  user: null,
  token: localStorage.getItem('token'),
  isAuthenticated: false,
  loading: false,
};

/**
 * 异步登录action
 * @param credentials 登录凭据
 */
export const loginAsync = createAsyncThunk(
  'auth/login',
  async (credentials: { username: string; password: string }) => {
    const response = await login(credentials);
    // 修复：正确访问嵌套在 data.data 中的 token
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
    }
    return response.data;
  }
);

/**
 * 异步获取用户信息action
 */
export const getUserProfileAsync = createAsyncThunk(
  'auth/getUserProfile',
  async () => {
    const response = await getUserProfile();
    return response.data;
  }
);

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    /**
     * 用户登出
     */
    logout: (state) => {
      state.user = null;
      state.token = null;
      state.isAuthenticated = false;
      localStorage.removeItem('token');
    },
    /**
     * 从本地存储设置认证状态
     */
    setAuthFromToken: (state) => {
      const token = localStorage.getItem('token');
      if (token) {
        state.token = token;
        state.isAuthenticated = true;
      }
    },
  },
  extraReducers: (builder) => {
    builder
      // 登录相关
      .addCase(loginAsync.pending, (state) => {
        state.loading = true;
      })
      .addCase(loginAsync.fulfilled, (state, action) => {
        state.loading = false;
        // action.payload 是 response.data (LoginResponse)
        state.user = action.payload.user_info;
        state.token = action.payload.token;
        state.isAuthenticated = true;
      })
      .addCase(loginAsync.rejected, (state) => {
        state.loading = false;
        state.isAuthenticated = false;
      })
      // 获取用户信息相关
      .addCase(getUserProfileAsync.pending, (state) => {
        state.loading = true;
      })
      .addCase(getUserProfileAsync.fulfilled, (state, action) => {
        state.loading = false;
        state.user = action.payload.data;
        state.isAuthenticated = true;
      })
      .addCase(getUserProfileAsync.rejected, (state) => {
        state.loading = false;
        // 如果获取用户信息失败，可能token已过期
        state.user = null;
        state.token = null;
        state.isAuthenticated = false;
        localStorage.removeItem('token');
      });
  },
});

export const { logout, setAuthFromToken } = authSlice.actions;
export default authSlice.reducer;