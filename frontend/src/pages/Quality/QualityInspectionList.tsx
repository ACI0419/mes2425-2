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
  DatePicker,
} from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import type { ColumnsType } from 'antd/es/table';
import {
  getQualityInspections,
  createQualityInspection,
  getQualityStatistics,
  type QualityInspectionResponse,
  type QualityInspectionRequest,
  type QualityStatistics,
  type QualityInspectionQueryParams, // 添加这个导入
} from '../../api/quality';
import dayjs from 'dayjs';

const { Option } = Select;

/**
 * 质量检测记录页面组件
 */
const QualityInspectionList: React.FC = () => {
  const [inspections, setInspections] = useState<QualityInspectionResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [statistics, setStatistics] = useState<QualityStatistics | null>(null);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  // 将 filters 的类型从 any 改为明确的类型
  const [filters, setFilters] = useState<Partial<QualityInspectionQueryParams>>({});

  /**
   * 获取质量检测记录列表
   */
  const fetchInspections = async (params?: Partial<QualityInspectionQueryParams>) => {
    setLoading(true);
    try {
      const response = await getQualityInspections({
        page: pagination.current,
        page_size: pagination.pageSize,
        ...filters,
        ...params,
      });
      setInspections(response.data.data.list);
      setPagination(prev => ({
        ...prev,
        total: response.data.data.total,
      }));
    } catch (error) {
      message.error('获取质量检测记录失败');
    } finally {
      setLoading(false);
    }
  };

  /**
   * 获取质量统计数据
   */
  const fetchStatistics = async () => {
    try {
      const response = await getQualityStatistics({});
      setStatistics(response.data.data);
    } catch (error) {
      console.error('获取质量统计数据失败:', error);
    }
  };

  /**
   * 处理新增检测记录
   */
  const handleSubmit = async (values: any) => {
    try {
      const inspectionData: QualityInspectionRequest = {
        ...values,
        inspector_id: 1, // TODO: 从当前用户获取
      };

      await createQualityInspection(inspectionData);
      message.success('创建质量检测记录成功');

      setModalVisible(false);
      form.resetFields();
      fetchInspections();
      fetchStatistics();
    } catch (error) {
      message.error('创建质量检测记录失败');
    }
  };

  /**
   * 获取检测结果标签颜色
   */
  const getResultColor = (result: string) => {
    return result === 'pass' ? 'success' : 'error';
  };

  /**
   * 获取检测结果文本
   */
  const getResultText = (result: string) => {
    return result === 'pass' ? '合格' : '不合格';
  };

  /**
   * 表格列定义
   */
  const columns: ColumnsType<QualityInspectionResponse> = [
    {
      title: '工单号',
      key: 'order_no',
      width: 150,
      render: (_, record) => record.production_order?.order_no,
    },
    {
      title: '产品信息',
      key: 'product',
      width: 150,
      render: (_, record) => (
        <div>
          <div>{record.production_order?.product?.name}</div>
          <div style={{ fontSize: '12px', color: '#666' }}>
            {record.production_order?.product?.code}
          </div>
        </div>
      ),
    },
    {
      title: '质量标准',
      key: 'quality_standard',
      width: 150,
      render: (_, record) => record.quality_standard?.name,
    },
    {
      title: '批次号',
      dataIndex: 'batch_no',
      key: 'batch_no',
      width: 120,
    },
    {
      title: '样本数量',
      dataIndex: 'sample_size',
      key: 'sample_size',
      width: 100,
    },
    {
      title: '测量值',
      dataIndex: 'measured_value',
      key: 'measured_value',
      width: 100,
      render: (value, record) => {
        if (value !== null && value !== undefined) {
          return `${value} ${record.quality_standard?.unit || ''}`;
        }
        return '-';
      },
    },
    {
      title: '检测结果',
      dataIndex: 'result',
      key: 'result',
      width: 100,
      render: (result) => (
        <Tag
          color={getResultColor(result)}
          icon={result === 'pass' ? <CheckCircleOutlined /> : <CloseCircleOutlined />}
        >
          {getResultText(result)}
        </Tag>
      ),
    },
    {
      title: '检测员',
      key: 'inspector',
      width: 100,
      render: (_, record) => record.inspector?.real_name || record.inspector?.username,
    },
    {
      title: '检测时间',
      dataIndex: 'inspected_at',
      key: 'inspected_at',
      width: 150,
      render: (date) => dayjs(date).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '备注',
      dataIndex: 'notes',
      key: 'notes',
      width: 150,
      ellipsis: true,
    },
  ];

  useEffect(() => {
    fetchInspections();
    fetchStatistics();
  }, []);

  useEffect(() => {
    fetchInspections();
  }, [pagination.current, pagination.pageSize]);

  return (
    <div>
      {/* 统计卡片 */}
      {statistics && (
        <Row gutter={16} style={{ marginBottom: 16 }}>
          <Col span={6}>
            <Card>
              <Statistic
                title="总检测次数"
                value={statistics.total_inspections}
                valueStyle={{ color: '#1890ff' }}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title="合格次数"
                value={statistics.pass_count}
                valueStyle={{ color: '#52c41a' }}
                prefix={<CheckCircleOutlined />}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title="不合格次数"
                value={statistics.fail_count}
                valueStyle={{ color: '#f5222d' }}
                prefix={<CloseCircleOutlined />}
              />
            </Card>
          </Col>
          <Col span={6}>
            <Card>
              <Statistic
                title="合格率"
                value={statistics.pass_rate}
                precision={2}
                suffix="%"
                valueStyle={{
                  color: statistics.pass_rate >= 95 ? '#52c41a' :
                    statistics.pass_rate >= 90 ? '#faad14' : '#f5222d'
                }}
              />
            </Card>
          </Col>
        </Row>
      )}

      <Card>
        <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
          <Space>
            <Select
              placeholder="检测结果"
              allowClear
              style={{ width: 120 }}
              onChange={(value) => setFilters(prev => ({ ...prev, result: value }))}
            >
              <Option value="pass">合格</Option>
              <Option value="fail">不合格</Option>
            </Select>
            <Button onClick={() => fetchInspections(filters)}>搜索</Button>
          </Space>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => {
              form.resetFields();
              setModalVisible(true);
            }}
          >
            新增检测记录
          </Button>
        </div>

        <Table
          columns={columns}
          dataSource={inspections}
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
        title="新增质量检测记录"
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
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
            name="production_order_id"
            label="生产工单"
            rules={[{ required: true, message: '请选择生产工单' }]}
          >
            <Select placeholder="选择生产工单">
              {/* TODO: 从生产工单API获取列表 */}
            </Select>
          </Form.Item>

          <Form.Item
            name="quality_standard_id"
            label="质量标准"
            rules={[{ required: true, message: '请选择质量标准' }]}
          >
            <Select placeholder="选择质量标准">
              {/* TODO: 从质量标准API获取列表 */}
            </Select>
          </Form.Item>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="batch_no"
                label="批次号"
              >
                <Input placeholder="请输入批次号" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="sample_size"
                label="样本数量"
                rules={[{ required: true, message: '请输入样本数量' }]}
              >
                <InputNumber min={1} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="measured_value"
                label="测量值"
              >
                <InputNumber style={{ width: '100%' }} placeholder="请输入测量值" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="result"
                label="检测结果"
                rules={[{ required: true, message: '请选择检测结果' }]}
              >
                <Select placeholder="选择检测结果">
                  <Option value="pass">合格</Option>
                  <Option value="fail">不合格</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="notes"
            label="备注"
          >
            <Input.TextArea rows={3} placeholder="请输入备注信息" />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                创建
              </Button>
              <Button onClick={() => {
                setModalVisible(false);
                form.resetFields();
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

export default QualityInspectionList;