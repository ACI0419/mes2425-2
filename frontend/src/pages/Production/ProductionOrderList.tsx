import React, { useState, useEffect } from 'react';
import {
  Table,
  Button,
  Space,
  Modal,
  Form,
  Input,
  Select,
  DatePicker,
  InputNumber,
  message,
  Popconfirm,
  Tag,
  Card,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getProductionOrders,
  createProductionOrder,
  updateProductionOrder,
  deleteProductionOrder,
  getProductionOrderStatuses,
  type ProductionOrderResponse,
  type ProductionOrderRequest,
  type ProductionOrderQueryParams,
} from '../../api/production';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;
const { Option } = Select;

/**
 * 生产工单列表页面组件
 */
const ProductionOrderList: React.FC = () => {
  const [orders, setOrders] = useState<ProductionOrderResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingOrder, setEditingOrder] = useState<ProductionOrderResponse | null>(null);
  const [form] = Form.useForm();
  const [statuses, setStatuses] = useState<string[]>([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [filters, setFilters] = useState<ProductionOrderQueryParams>({});

  /**
   * 获取生产工单列表
   */
  const fetchOrders = async (params?: ProductionOrderQueryParams) => {
    setLoading(true);
    try {
      const response = await getProductionOrders({
        page: pagination.current,
        page_size: pagination.pageSize,
        ...filters,
        ...params,
      });
      setOrders(response.data.data.list);
      setPagination(prev => ({
        ...prev,
        total: response.data.data.total,
      }));
    } catch (error) {
      message.error('获取生产工单列表失败');
    } finally {
      setLoading(false);
    }
  };

  /**
   * 获取工单状态列表
   */
  const fetchStatuses = async () => {
    try {
      const response = await getProductionOrderStatuses();
      setStatuses(response.data.data);
    } catch (error) {
      console.error('获取状态列表失败:', error);
    }
  };

  /**
   * 处理新增/编辑工单
   */
  const handleSubmit = async (values: any) => {
    try {
      // 先提取 dateRange，然后从 values 中移除它
      const { dateRange, ...restValues } = values;
      
      const orderData: ProductionOrderRequest = {
        ...restValues,
        start_date: dateRange[0].format('YYYY-MM-DD'),
        end_date: dateRange[1].format('YYYY-MM-DD'),
      };

      if (editingOrder) {
        await updateProductionOrder(editingOrder.id, orderData);
        message.success('更新生产工单成功');
      } else {
        await createProductionOrder(orderData);
        message.success('创建生产工单成功');
      }

      setModalVisible(false);
      form.resetFields();
      setEditingOrder(null);
      fetchOrders();
    } catch (error) {
      message.error(editingOrder ? '更新生产工单失败' : '创建生产工单失败');
    }
  };

  /**
   * 处理删除工单
   */
  const handleDelete = async (id: number) => {
    try {
      await deleteProductionOrder(id);
      message.success('删除生产工单成功');
      fetchOrders();
    } catch (error) {
      message.error('删除生产工单失败');
    }
  };

  /**
   * 打开编辑模态框
   */
  const handleEdit = (order: ProductionOrderResponse) => {
    setEditingOrder(order);
    form.setFieldsValue({
      ...order,
      dateRange: [dayjs(order.start_date), dayjs(order.end_date)],
    });
    setModalVisible(true);
  };

  /**
   * 获取状态标签颜色
   */
  const getStatusColor = (status: string) => {
    const colorMap: Record<string, string> = {
      pending: 'orange',
      in_progress: 'blue',
      completed: 'green',
      cancelled: 'red',
    };
    return colorMap[status] || 'default';
  };

  /**
   * 获取状态显示文本
   */
  const getStatusText = (status: string) => {
    const textMap: Record<string, string> = {
      pending: '待开始',
      in_progress: '进行中',
      completed: '已完成',
      cancelled: '已取消',
    };
    return textMap[status] || status;
  };

  /**
   * 表格列定义
   */
  const columns: ColumnsType<ProductionOrderResponse> = [
    {
      title: '工单号',
      dataIndex: 'order_no',
      key: 'order_no',
      width: 150,
    },
    {
      title: '产品信息',
      key: 'product',
      width: 200,
      render: (_, record) => (
        <div>
          <div>{record.product?.name}</div>
          <div style={{ fontSize: '12px', color: '#666' }}>
            {record.product?.code}
          </div>
        </div>
      ),
    },
    {
      title: '计划数量',
      dataIndex: 'quantity',
      key: 'quantity',
      width: 100,
      render: (quantity, record) => `${quantity} ${record.product?.unit || ''}`,
    },
    {
      title: '已生产',
      dataIndex: 'produced',
      key: 'produced',
      width: 100,
      render: (produced, record) => `${produced} ${record.product?.unit || ''}`,
    },
    {
      title: '进度',
      key: 'progress',
      width: 100,
      render: (_, record) => {
        const progress = record.quantity > 0 ? (record.produced / record.quantity * 100).toFixed(1) : '0';
        return `${progress}%`;
      },
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      key: 'priority',
      width: 80,
      render: (priority) => {
        const colorMap: Record<number, string> = {
          1: 'red',
          2: 'orange',
          3: 'blue',
          4: 'green',
          5: 'gray',
        };
        return <Tag color={colorMap[priority] || 'default'}>P{priority}</Tag>;
      },
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      width: 100,
      render: (status) => (
        <Tag color={getStatusColor(status)}>
          {getStatusText(status)}
        </Tag>
      ),
    },
    {
      title: '开始日期',
      dataIndex: 'start_date',
      key: 'start_date',
      width: 120,
      render: (date) => dayjs(date).format('YYYY-MM-DD'),
    },
    {
      title: '结束日期',
      dataIndex: 'end_date',
      key: 'end_date',
      width: 120,
      render: (date) => dayjs(date).format('YYYY-MM-DD'),
    },
    {
      title: '创建人',
      key: 'creator',
      width: 100,
      render: (_, record) => record.creator?.real_name || record.creator?.username,
    },
    {
      title: '操作',
      key: 'action',
      width: 150,
      fixed: 'right',
      render: (_, record) => (
        <Space size="small">
          <Button
            type="link"
            icon={<EyeOutlined />}
            size="small"
            onClick={() => {/* TODO: 查看详情 */ }}
          >
            查看
          </Button>
          <Button
            type="link"
            icon={<EditOutlined />}
            size="small"
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个生产工单吗？"
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
    fetchOrders();
    fetchStatuses();
  }, []);

  useEffect(() => {
    fetchOrders();
  }, [pagination.current, pagination.pageSize]);

  return (
    <div>
      <Card>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
          <Space>
            <Select
              placeholder="选择状态"
              allowClear
              style={{ width: 120 }}
              onChange={(value) => setFilters(prev => ({ ...prev, status: value }))}
            >
              {statuses.map(status => (
                <Option key={status} value={status}>
                  {getStatusText(status)}
                </Option>
              ))}
            </Select>
            <Button onClick={() => fetchOrders(filters)}>搜索</Button>
          </Space>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              setEditingOrder(null);
              form.resetFields();
              setModalVisible(true);
            }}
          >
            新增工单
          </Button>
        </div>

        <Table
          columns={columns}
          dataSource={orders}
          rowKey="id"
          loading={loading}
          scroll={{ x: 1200 }}
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
        title={editingOrder ? '编辑生产工单' : '新增生产工单'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
          setEditingOrder(null);
        }}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="product_id"
            label="产品"
            rules={[{ required: true, message: '请选择产品' }]}
          >
            <Select placeholder="选择产品">
              {/* TODO: 从产品API获取产品列表 */}
            </Select>
          </Form.Item>

          <Form.Item
            name="quantity"
            label="计划数量"
            rules={[{ required: true, message: '请输入计划数量' }]}
          >
            <InputNumber min={1} style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            name="priority"
            label="优先级"
            rules={[{ required: true, message: '请选择优先级' }]}
          >
            <Select placeholder="选择优先级">
              <Option value={1}>P1 - 紧急</Option>
              <Option value={2}>P2 - 高</Option>
              <Option value={3}>P3 - 中</Option>
              <Option value={4}>P4 - 低</Option>
              <Option value={5}>P5 - 最低</Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="dateRange"
            label="生产周期"
            rules={[{ required: true, message: '请选择生产周期' }]}
          >
            <RangePicker style={{ width: '100%' }} />
          </Form.Item>

          <Form.Item
            name="description"
            label="备注"
          >
            <Input.TextArea rows={3} placeholder="请输入备注信息" />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingOrder ? '更新' : '创建'}
              </Button>
              <Button onClick={() => {
                setModalVisible(false);
                form.resetFields();
                setEditingOrder(null);
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

export default ProductionOrderList;