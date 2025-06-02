import React, { useState, useEffect } from 'react';
import { Table, Button, Space, Input, Select, Modal, Form, message, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, ToolOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import type { Equipment } from '../../types';
import { getEquipments, createEquipment, updateEquipment, deleteEquipment } from '../../api/equipment';

const { Search } = Input;
const { Option } = Select;

/**
 * 设备列表组件
 * 提供设备的查看、新增、编辑、删除和维护记录管理功能
 */
const EquipmentList: React.FC = () => {
  const navigate = useNavigate();
  const [equipments, setEquipments] = useState<Equipment[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingEquipment, setEditingEquipment] = useState<Equipment | null>(null);
  const [form] = Form.useForm();
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [filters, setFilters] = useState({
    keyword: '',
    type: '',
    status: '',
  });

  /**
   * 获取设备列表数据
   */
  const fetchEquipments = async () => {
    setLoading(true);
    try {
      const response = await getEquipments({
        page: pagination.current,
        page_size: pagination.pageSize,
        ...filters,
      });
      setEquipments(response.data.data.list);
      setPagination(prev => ({
        ...prev,
        total: response.data.data.total,
      }));
    } catch (error) {
      message.error('获取设备列表失败');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEquipments();
  }, [pagination.current, pagination.pageSize, filters]);

  // 表格列定义
  const columns = [
    {
      title: '设备编码',
      dataIndex: 'code',
      key: 'code',
    },
    {
      title: '设备名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '设备类型',
      dataIndex: 'type',
      key: 'type',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        const statusMap = {
          running: { text: '运行中', color: 'green' },
          stopped: { text: '已停止', color: 'orange' },
          maintenance: { text: '维护中', color: 'blue' },
          fault: { text: '故障', color: 'red' },
        };
        const config = statusMap[status as keyof typeof statusMap];
        return <span style={{ color: config?.color }}>{config?.text}</span>;
      },
    },
    {
      title: '位置',
      dataIndex: 'location',
      key: 'location',
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Equipment) => (
        <Space size="middle">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Button
            type="link"
            icon={<ToolOutlined />}
            onClick={() => {
              // 使用 React Router 进行页面跳转
              navigate(`/equipment/maintenance?equipment_id=${record.id}`);
            }}
          >
            维护记录
          </Button>
          <Popconfirm
            title="确定要删除这个设备吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  /**
   * 处理编辑设备
   * @param equipment 要编辑的设备信息
   */
  const handleEdit = (equipment: Equipment) => {
    setEditingEquipment(equipment);
    form.setFieldsValue(equipment);
    setModalVisible(true);
  };

  /**
   * 处理删除设备
   * @param id 设备ID
   */
  const handleDelete = async (id: number) => {
    try {
      await deleteEquipment(id);
      message.success('删除成功');
      fetchEquipments();
    } catch (error) {
      message.error('删除失败');
    }
  };

  /**
   * 处理表单提交
   * @param values 表单数据
   */
  const handleSubmit = async (values: any) => {
    try {
      if (editingEquipment) {
        await updateEquipment(editingEquipment.id, values);
        message.success('更新成功');
      } else {
        await createEquipment(values);
        message.success('创建成功');
      }
      setModalVisible(false);
      form.resetFields();
      setEditingEquipment(null);
      fetchEquipments();
    } catch (error) {
      message.error('操作失败');
    }
  };

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Space>
          <Search
            placeholder="搜索设备"
            allowClear
            style={{ width: 200 }}
            onSearch={(value) => setFilters(prev => ({ ...prev, keyword: value }))}
          />
          <Select
            placeholder="设备类型"
            allowClear
            style={{ width: 120 }}
            onChange={(value) => setFilters(prev => ({ ...prev, type: value || '' }))}
          >
            <Option value="production">生产设备</Option>
            <Option value="testing">检测设备</Option>
            <Option value="auxiliary">辅助设备</Option>
          </Select>
          <Select
            placeholder="设备状态"
            allowClear
            style={{ width: 120 }}
            onChange={(value) => setFilters(prev => ({ ...prev, status: value || '' }))}
          >
            <Option value="running">运行中</Option>
            <Option value="stopped">已停止</Option>
            <Option value="maintenance">维护中</Option>
            <Option value="fault">故障</Option>
          </Select>
        </Space>
        <Button 
          type="primary" 
          icon={<PlusOutlined />}
          onClick={() => {
            setEditingEquipment(null);
            form.resetFields();
            setModalVisible(true);
          }}
        >
          新增设备
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={equipments}
        rowKey="id"
        loading={loading}
        pagination={{
          current: pagination.current,
          pageSize: pagination.pageSize,
          total: pagination.total,
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条记录`,
          onChange: (page, pageSize) => {
            setPagination(prev => ({
              ...prev,
              current: page,
              pageSize: pageSize || 10,
            }));
          },
        }}
      />

      <Modal
        title={editingEquipment ? '编辑设备' : '新增设备'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
          setEditingEquipment(null);
        }}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="code"
            label="设备编码"
            rules={[{ required: true, message: '请输入设备编码' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="name"
            label="设备名称"
            rules={[{ required: true, message: '请输入设备名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="type"
            label="设备类型"
            rules={[{ required: true, message: '请选择设备类型' }]}
          >
            <Select>
              <Option value="production">生产设备</Option>
              <Option value="testing">检测设备</Option>
              <Option value="auxiliary">辅助设备</Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="status"
            label="设备状态"
            rules={[{ required: true, message: '请选择设备状态' }]}
          >
            <Select>
              <Option value="running">运行中</Option>
              <Option value="stopped">已停止</Option>
              <Option value="maintenance">维护中</Option>
              <Option value="fault">故障</Option>
            </Select>
          </Form.Item>
          <Form.Item name="model" label="设备型号">
            <Input />
          </Form.Item>
          <Form.Item name="manufacturer" label="制造商">
            <Input />
          </Form.Item>
          <Form.Item name="location" label="设备位置">
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingEquipment ? '更新' : '创建'}
              </Button>
              <Button onClick={() => setModalVisible(false)}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default EquipmentList;