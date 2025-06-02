import React, { useState, useEffect } from 'react';
import { Table, Button, Space, Input, Select, Modal, Form, message, DatePicker, InputNumber, Card, Statistic, Row, Col, Popconfirm } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, ToolOutlined, CalendarOutlined } from '@ant-design/icons';
import { useSearchParams } from 'react-router-dom';
import dayjs from 'dayjs';
import type { Equipment } from '../../types';
import {
  getMaintenanceRecords,
  createMaintenanceRecord,
  updateMaintenanceRecord,
  deleteMaintenanceRecord,
  getEquipments,
  getUpcomingMaintenance,
  type MaintenanceRecordResponse,
  type MaintenanceRecordRequest,
  type MaintenanceRecordQueryParams
} from '../../api/equipment';

const { Search } = Input;
const { Option } = Select;
const { RangePicker } = DatePicker;

/**
 * 维护记录列表组件
 * 提供维护记录的查看、新增、编辑、删除和统计功能
 */
const MaintenanceRecordList: React.FC = () => {
  const [searchParams] = useSearchParams();
  const [records, setRecords] = useState<MaintenanceRecordResponse[]>([]);
  const [equipments, setEquipments] = useState<Equipment[]>([]);
  const [upcomingRecords, setUpcomingRecords] = useState<MaintenanceRecordResponse[]>([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRecord, setEditingRecord] = useState<MaintenanceRecordResponse | null>(null);
  const [form] = Form.useForm();
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [filters, setFilters] = useState<MaintenanceRecordQueryParams>({
    equipment_id: searchParams.get('equipment_id') ? Number(searchParams.get('equipment_id')) : undefined,
    type: '',
    maintainer_id: undefined,
    start_date: '',
    end_date: '',
  });

  /**
   * 获取维护记录列表
   */
  const fetchRecords = async () => {
    setLoading(true);
    try {
      const response = await getMaintenanceRecords({
        page: pagination.current,
        page_size: pagination.pageSize,
        ...filters,
      });
      setRecords(response.data.data.list);
      setPagination(prev => ({
        ...prev,
        total: response.data.data.total,
      }));
    } catch (error) {
      message.error('获取维护记录列表失败');
    } finally {
      setLoading(false);
    }
  };

  /**
   * 获取设备列表（用于下拉选择）
   */
  const fetchEquipments = async () => {
    try {
      const response = await getEquipments({ page_size: 1000 });
      setEquipments(response.data.data.list);
    } catch (error) {
      console.error('获取设备列表失败:', error);
    }
  };

  /**
   * 获取即将到期的维护记录
   */
  const fetchUpcomingMaintenance = async () => {
    try {
      const response = await getUpcomingMaintenance(7);
      setUpcomingRecords(response.data.data);
    } catch (error) {
      console.error('获取即将到期维护记录失败:', error);
    }
  };

  useEffect(() => {
    fetchRecords();
  }, [pagination.current, pagination.pageSize, filters]);

  useEffect(() => {
    fetchEquipments();
    fetchUpcomingMaintenance();
  }, []);

  // 表格列定义
  const columns = [
    {
      title: '设备编码',
      dataIndex: 'equipment_code',
      key: 'equipment_code',
    },
    {
      title: '设备名称',
      dataIndex: 'equipment_name',
      key: 'equipment_name',
    },
    {
      title: '维护类型',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => {
        const typeMap = {
          preventive: { text: '预防性维护', color: 'blue' },
          corrective: { text: '纠正性维护', color: 'orange' },
          emergency: { text: '紧急维护', color: 'red' },
        };
        const config = typeMap[type as keyof typeof typeMap];
        return <span style={{ color: config?.color }}>{config?.text || type}</span>;
      },
    },
    {
      title: '维护人员',
      dataIndex: 'maintainer_name',
      key: 'maintainer_name',
    },
    {
      title: '开始时间',
      dataIndex: 'start_time',
      key: 'start_time',
      render: (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm'),
    },
    {
      title: '结束时间',
      dataIndex: 'end_time',
      key: 'end_time',
      render: (time: string) => time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-',
    },
    {
      title: '维护费用',
      dataIndex: 'cost',
      key: 'cost',
      render: (cost: number) => cost ? `¥${cost.toFixed(2)}` : '-',
    },
    {
      title: '下次维护',
      dataIndex: 'next_maintenance',
      key: 'next_maintenance',
      render: (time: string) => time ? dayjs(time).format('YYYY-MM-DD') : '-',
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: MaintenanceRecordResponse) => (
        <Space size="middle">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这条维护记录吗？"
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
   * 处理编辑维护记录
   * @param record 要编辑的维护记录
   */
  const handleEdit = (record: MaintenanceRecordResponse) => {
    setEditingRecord(record);
    form.setFieldsValue({
      ...record,
      start_time: dayjs(record.start_time),
      end_time: record.end_time ? dayjs(record.end_time) : undefined,
      next_maintenance: record.next_maintenance ? dayjs(record.next_maintenance) : undefined,
    });
    setModalVisible(true);
  };

  /**
   * 处理删除维护记录
   * @param id 维护记录ID
   */
  const handleDelete = async (id: number) => {
    try {
      await deleteMaintenanceRecord(id);
      message.success('删除成功');
      fetchRecords();
      fetchUpcomingMaintenance();
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
      const submitData: MaintenanceRecordRequest = {
        ...values,
        start_time: values.start_time.format('YYYY-MM-DD HH:mm:ss'),
        end_time: values.end_time ? values.end_time.format('YYYY-MM-DD HH:mm:ss') : undefined,
        next_maintenance: values.next_maintenance ? values.next_maintenance.format('YYYY-MM-DD') : undefined,
      };

      if (editingRecord) {
        await updateMaintenanceRecord(editingRecord.id, submitData);
        message.success('更新成功');
      } else {
        await createMaintenanceRecord(submitData);
        message.success('创建成功');
      }
      setModalVisible(false);
      form.resetFields();
      setEditingRecord(null);
      fetchRecords();
      fetchUpcomingMaintenance();
    } catch (error) {
      message.error('操作失败');
    }
  };

  /**
   * 处理日期范围筛选
   * @param dates 日期范围
   */
  const handleDateRangeChange = (dates: any) => {
    if (dates && dates.length === 2) {
      setFilters(prev => ({
        ...prev,
        start_date: dates[0].format('YYYY-MM-DD'),
        end_date: dates[1].format('YYYY-MM-DD'),
      }));
    } else {
      setFilters(prev => ({
        ...prev,
        start_date: '',
        end_date: '',
      }));
    }
  };

  return (
    <div>
      {/* 统计卡片 */}
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Card>
            <Statistic
              title="总维护记录"
              value={pagination.total}
              prefix={<ToolOutlined />}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="即将到期维护"
              value={upcomingRecords.length}
              prefix={<CalendarOutlined />}
              valueStyle={{ color: upcomingRecords.length > 0 ? '#cf1322' : undefined }}
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="本月维护费用"
              value={records.reduce((sum, record) => sum + (record.cost || 0), 0)}
              precision={2}
              prefix="¥"
            />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic
              title="平均维护时长"
              value={records.length > 0 ? records.reduce((sum, record) => sum + (record.duration || 0), 0) / records.length : 0}
              suffix="分钟"
              precision={0}
            />
          </Card>
        </Col>
      </Row>

      {/* 筛选和操作区域 */}
      <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
        <Space>
          <Select
            placeholder="选择设备"
            allowClear
            style={{ width: 200 }}
            value={filters.equipment_id}
            onChange={(value) => setFilters(prev => ({ ...prev, equipment_id: value }))}
          >
            {equipments.map(equipment => (
              <Option key={equipment.id} value={equipment.id}>
                {equipment.code} - {equipment.name}
              </Option>
            ))}
          </Select>
          <Select
            placeholder="维护类型"
            allowClear
            style={{ width: 120 }}
            onChange={(value) => setFilters(prev => ({ ...prev, type: value || '' }))}
          >
            <Option value="preventive">预防性维护</Option>
            <Option value="corrective">纠正性维护</Option>
            <Option value="emergency">紧急维护</Option>
          </Select>
          <RangePicker
            placeholder={['开始日期', '结束日期']}
            onChange={handleDateRangeChange}
          />
        </Space>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => {
            setEditingRecord(null);
            form.resetFields();
            setModalVisible(true);
          }}
        >
          新增维护记录
        </Button>
      </div>

      {/* 维护记录表格 */}
      <Table
        columns={columns}
        dataSource={records}
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

      {/* 新增/编辑维护记录模态框 */}
      <Modal
        title={editingRecord ? '编辑维护记录' : '新增维护记录'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
          setEditingRecord(null);
        }}
        footer={null}
        width={800}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="equipment_id"
                label="设备"
                rules={[{ required: true, message: '请选择设备' }]}
              >
                <Select placeholder="选择设备">
                  {equipments.map(equipment => (
                    <Option key={equipment.id} value={equipment.id}>
                      {equipment.code} - {equipment.name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="type"
                label="维护类型"
                rules={[{ required: true, message: '请选择维护类型' }]}
              >
                <Select placeholder="选择维护类型">
                  <Option value="preventive">预防性维护</Option>
                  <Option value="corrective">纠正性维护</Option>
                  <Option value="emergency">紧急维护</Option>
                </Select>
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="maintainer_id"
                label="维护人员ID"
                rules={[{ required: true, message: '请输入维护人员ID' }]}
              >
                <InputNumber style={{ width: '100%' }} placeholder="维护人员ID" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="cost"
                label="维护费用"
              >
                <InputNumber
                  style={{ width: '100%' }}
                  placeholder="维护费用"
                  precision={2}
                  min={0}
                  addonBefore="¥"
                />
              </Form.Item>
            </Col>
          </Row>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="start_time"
                label="开始时间"
                rules={[{ required: true, message: '请选择开始时间' }]}
              >
                <DatePicker
                  showTime
                  style={{ width: '100%' }}
                  placeholder="选择开始时间"
                  format="YYYY-MM-DD HH:mm"
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="end_time"
                label="结束时间"
              >
                <DatePicker
                  showTime
                  style={{ width: '100%' }}
                  placeholder="选择结束时间"
                  format="YYYY-MM-DD HH:mm"
                />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item
            name="description"
            label="维护描述"
            rules={[{ required: true, message: '请输入维护描述' }]}
          >
            <Input.TextArea rows={3} placeholder="详细描述维护内容" />
          </Form.Item>
          <Form.Item name="parts_replaced" label="更换部件">
            <Input.TextArea rows={2} placeholder="列出更换的部件" />
          </Form.Item>
          <Form.Item name="result" label="维护结果">
            <Input.TextArea rows={2} placeholder="维护结果和效果" />
          </Form.Item>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="next_maintenance" label="下次维护时间">
                <DatePicker
                  style={{ width: '100%' }}
                  placeholder="选择下次维护时间"
                  format="YYYY-MM-DD"
                />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item name="remark" label="备注">
            <Input.TextArea rows={2} placeholder="其他备注信息" />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingRecord ? '更新' : '创建'}
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

export default MaintenanceRecordList;