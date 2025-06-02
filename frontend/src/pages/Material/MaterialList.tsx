import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  InputNumber,
  message,
  Popconfirm,
  Tag,
  Card,
  Row,
  Col,
  Statistic,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, WarningOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getMaterials,
  createMaterial,
  updateMaterial,
  deleteMaterial,
  getMaterialTypes,
  getLowStockMaterials,
  type MaterialResponse,
  type MaterialRequest,
  type MaterialQueryParams,
} from '../../api/material';

const { Option } = Select;
const { Search } = Input;

/**
 * 物料管理页面组件
 */
const MaterialList: React.FC = () => {
  const [materials, setMaterials] = useState<MaterialResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingMaterial, setEditingMaterial] = useState<MaterialResponse | null>(null);
  const [form] = Form.useForm();
  const [materialTypes, setMaterialTypes] = useState<string[]>([]);
  const [lowStockCount, setLowStockCount] = useState(0);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [filters, setFilters] = useState<MaterialQueryParams>({});

  /**
   * 获取物料列表
   */
  const fetchMaterials = async (params?: MaterialQueryParams) => {
    setLoading(true);
    try {
      const response = await getMaterials({
        page: pagination.current,
        page_size: pagination.pageSize,
        ...filters,
        ...params,
      });
      setMaterials(response.data.data.list);
      setPagination(prev => ({
        ...prev,
        total: response.data.data.total,
      }));
    } catch (error) {
      message.error('获取物料列表失败');
    } finally {
      setLoading(false);
    }
  };

  /**
   * 获取物料类型列表
   */
  const fetchMaterialTypes = async () => {
    try {
      const response = await getMaterialTypes();
      setMaterialTypes(response.data.data);
    } catch (error) {
      console.error('获取物料类型失败:', error);
    }
  };

  /**
   * 获取低库存物料数量
   */
  const fetchLowStockCount = async () => {
    try {
      const response = await getLowStockMaterials();
      setLowStockCount(response.data.data.length);
    } catch (error) {
      console.error('获取低库存物料失败:', error);
    }
  };

  /**
   * 处理新增/编辑物料
   */
  const handleSubmit = async (values: MaterialRequest) => {
    try {
      if (editingMaterial) {
        await updateMaterial(editingMaterial.id, values);
        message.success('更新物料成功');
      } else {
        await createMaterial(values);
        message.success('创建物料成功');
      }

      setModalVisible(false);
      form.resetFields();
      setEditingMaterial(null);
      fetchMaterials();
      fetchLowStockCount();
    } catch (error) {
      message.error(editingMaterial ? '更新物料失败' : '创建物料失败');
    }
  };

  /**
   * 处理删除物料
   */
  const handleDelete = async (id: number) => {
    try {
      await deleteMaterial(id);
      message.success('删除物料成功');
      fetchMaterials();
      fetchLowStockCount();
    } catch (error) {
      message.error('删除物料失败');
    }
  };

  /**
   * 打开编辑模态框
   */
  const handleEdit = (material: MaterialResponse) => {
    setEditingMaterial(material);
    form.setFieldsValue(material);
    setModalVisible(true);
  };

  /**
   * 获取库存状态
   */
  const getStockStatus = (material: MaterialResponse) => {
    if (material.current_stock <= material.min_stock) {
      return { status: 'error', text: '库存不足' };
    } else if (material.current_stock <= material.min_stock * 1.2) {
      return { status: 'warning', text: '库存偏低' };
    } else if (material.current_stock >= material.max_stock) {
      return { status: 'processing', text: '库存过多' };
    }
    return { status: 'success', text: '库存正常' };
  };

  /**
   * 表格列定义
   */
  const columns: ColumnsType<MaterialResponse> = [
    {
      title: '物料编码',
      dataIndex: 'code',
      key: 'code',
      width: 120,
    },
    {
      title: '物料名称',
      dataIndex: 'name',
      key: 'name',
      width: 150,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      width: 100,
    },
    {
      title: '单位',
      dataIndex: 'unit',
      key: 'unit',
      width: 80,
    },
    {
      title: '单价',
      dataIndex: 'price',
      key: 'price',
      width: 100,
      render: (price) => `¥${price.toFixed(2)}`,
    },
    {
      title: '当前库存',
      dataIndex: 'current_stock',
      key: 'current_stock',
      width: 100,
      render: (stock, record) => `${stock} ${record.unit}`,
    },
    {
      title: '库存范围',
      key: 'stock_range',
      width: 120,
      render: (_, record) => (
        <div>
          <div style={{ fontSize: '12px' }}>最小: {record.min_stock}</div>
          <div style={{ fontSize: '12px' }}>最大: {record.max_stock}</div>
        </div>
      ),
    },
    {
      title: '库存状态',
      key: 'stock_status',
      width: 100,
      render: (_, record) => {
        const { status, text } = getStockStatus(record);
        return <Tag color={status}>{text}</Tag>;
      },
    },
    {
      title: '供应商',
      dataIndex: 'supplier',
      key: 'supplier',
      width: 120,
    },
    {
      title: '操作',
      key: 'action',
      width: 120,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EditOutlined />}
            size="small"
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个物料吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
              size="small"
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  useEffect(() => {
    fetchMaterials();
    fetchMaterialTypes();
    fetchLowStockCount();
  }, []);

  useEffect(() => {
    fetchMaterials();
  }, [pagination.current, pagination.pageSize]);

  return (
    <div>
      {/* 统计卡片 */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="物料总数"
              value={pagination.total}
              valueStyle={{ color: '#1890ff' }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="低库存物料"
              value={lowStockCount}
              valueStyle={{ color: '#cf1322' }}
              prefix={<WarningOutlined />}
            />
          </Card>
        </Col>
      </Row>

      <Card>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
          <Space>
            <Select
              placeholder="选择类型"
              allowClear
              style={{ width: 120 }}
              onChange={(value) => setFilters(prev => ({ ...prev, type: value }))}
            >
              {materialTypes.map(type => (
                <Option key={type} value={type}>
                  {type}
                </Option>
              ))}
            </Select>
            <Search
              placeholder="搜索物料名称或编码"
              style={{ width: 200 }}
              onSearch={(value) => {
                setFilters(prev => ({ ...prev, keyword: value }));
                fetchMaterials({ keyword: value });
              }}
            />
          </Space>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              setEditingMaterial(null);
              form.resetFields();
              setModalVisible(true);
            }}
          >
            新增物料
          </Button>
        </div>

        <Table
          columns={columns}
          dataSource={materials}
          rowKey="id"
          loading={loading}
          scroll={{ x: 1000 }}
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
      </Card>

      <Modal
        title={editingMaterial ? '编辑物料' : '新增物料'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
          setEditingMaterial(null);
        }}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="code"
                label="物料编码"
                rules={[{ required: true, message: '请输入物料编码' }]}
              >
                <Input placeholder="请输入物料编码" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="name"
                label="物料名称"
                rules={[{ required: true, message: '请输入物料名称' }]}
              >
                <Input placeholder="请输入物料名称" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="type"
                label="物料类型"
                rules={[{ required: true, message: '请选择物料类型' }]}
              >
                <Select placeholder="选择物料类型">
                  {materialTypes.map(type => (
                    <Option key={type} value={type}>
                      {type}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="unit"
                label="单位"
                rules={[{ required: true, message: '请输入单位' }]}
              >
                <Input placeholder="如：个、kg、米等" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="price"
                label="单价"
                rules={[{ required: true, message: '请输入单价' }]}
              >
                <InputNumber
                  min={0}
                  precision={2}
                  style={{ width: '100%' }}
                  placeholder="请输入单价"
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="current_stock"
                label="当前库存"
                rules={[{ required: true, message: '请输入当前库存' }]}
              >
                <InputNumber
                  min={0}
                  style={{ width: '100%' }}
                  placeholder="请输入当前库存"
                />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="min_stock"
                label="最小库存"
                rules={[{ required: true, message: '请输入最小库存' }]}
              >
                <InputNumber
                  min={0}
                  style={{ width: '100%' }}
                  placeholder="请输入最小库存"
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="max_stock"
                label="最大库存"
                rules={[{ required: true, message: '请输入最大库存' }]}
              >
                <InputNumber
                  min={0}
                  style={{ width: '100%' }}
                  placeholder="请输入最大库存"
                />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="supplier"
            label="供应商"
          >
            <Input placeholder="请输入供应商名称" />
          </Form.Item>

          <Form.Item
            name="description"
            label="描述"
          >
            <Input.TextArea rows={3} placeholder="请输入物料描述" />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingMaterial ? '更新' : '创建'}
              </Button>
              <Button onClick={() => {
                setModalVisible(false);
                form.resetFields();
                setEditingMaterial(null);
              }}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default MaterialList;