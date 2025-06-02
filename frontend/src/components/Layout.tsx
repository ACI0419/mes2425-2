import React from 'react';
import { Layout as AntLayout, Menu, Avatar, Dropdown } from 'antd';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import {
  DashboardOutlined,
  ToolOutlined,
  ShoppingCartOutlined,
  InboxOutlined,
  SafetyCertificateOutlined,
  UserOutlined,
  LogoutOutlined,
} from '@ant-design/icons';
import type { RootState } from '../store';
import { logout } from '../store/slices/authSlice';

const { Header, Sider, Content } = AntLayout;

/**
 * 主布局组件
 */
const Layout: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useDispatch();
  const { user } = useSelector((state: RootState) => state.auth);

  /**
   * 处理用户登出
   */
  const handleLogout = () => {
    dispatch(logout());
    navigate('/login');
  };

  /**
   * 菜单项配置
   */
  const menuItems = [
    {
      key: '/dashboard',
      icon: <DashboardOutlined />,
      label: '仪表板',
      onClick: () => navigate('/dashboard'),
    },
    {
      key: '/production',
      icon: <ShoppingCartOutlined />,
      label: '生产管理',
      onClick: () => navigate('/production'),
    },
    {
      key: '/material',
      icon: <InboxOutlined />,
      label: '物料管理',
      onClick: () => navigate('/material'),
    },
    {
      key: '/quality',
      icon: <SafetyCertificateOutlined />,
      label: '质量管理',
      onClick: () => navigate('/quality'),
    },
    {
      key: 'equipment',
      icon: <ToolOutlined />,
      label: '设备管理',
      children: [
        {
          key: '/equipment',
          label: '设备列表',
          onClick: () => navigate('/equipment'),
        },
        {
          key: '/equipment/maintenance',
          label: '维护记录',
          onClick: () => navigate('/equipment/maintenance'),
        },
      ],
    },
  ];

  /**
   * 用户下拉菜单
   */
  const userMenuItems = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人信息',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ];

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider width={200} theme="dark">
        <div style={{ padding: '16px', color: 'white', textAlign: 'center' }}>
          <h3>MES系统</h3>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
        />
      </Sider>
      <AntLayout>
        <Header style={{ background: '#fff', padding: '0 16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2 style={{ margin: 0 }}>制造执行系统</h2>
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <div style={{ cursor: 'pointer', display: 'flex', alignItems: 'center', gap: 8 }}>
              <Avatar icon={<UserOutlined />} />
              <span>{user?.real_name || user?.username}</span>
            </div>
          </Dropdown>
        </Header>
        <Content style={{ margin: '16px', padding: '16px', background: '#fff' }}>
          <Outlet />
        </Content>
      </AntLayout>
    </AntLayout>
  );
};

export default Layout;