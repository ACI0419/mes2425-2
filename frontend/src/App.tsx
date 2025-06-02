import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { Provider } from 'react-redux';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { store } from './store';
import Layout from './components/Layout';
import Login from './pages/Auth/Login';
import Dashboard from './pages/Dashboard';
import EquipmentList from './pages/Equipment/EquipmentList';
import ProductionOrderList from './pages/Production/ProductionOrderList';
import MaterialList from './pages/Material/MaterialList';
import QualityInspectionList from './pages/Quality/QualityInspectionList';
import MaintenanceRecordList from './pages/Equipment/MaintenanceRecordList';

// 在路由配置中添加
<Route path="/equipment/maintenance" element={<MaintenanceRecordList />} />

function App() {
  return (
    <Provider store={store}>
      <ConfigProvider locale={zhCN}>
        <Router>
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/" element={<Layout />}>
              <Route index element={<Navigate to="/dashboard" replace />} />
              <Route path="dashboard" element={<Dashboard />} />
              <Route path="equipment" element={<EquipmentList />} />
              {/* 其他路由... */}
              // 在路由配置中添加：
              <Route path="/production" element={<ProductionOrderList />} />
              <Route path="/material" element={<MaterialList />} />
              <Route path="/quality" element={<QualityInspectionList />} />
            </Route>
          </Routes>
        </Router>
      </ConfigProvider>
    </Provider>
  );
}

export default App;