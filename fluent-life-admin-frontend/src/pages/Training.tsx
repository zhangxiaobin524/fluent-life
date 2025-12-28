import React, { useEffect, useState } from 'react';
import { adminAPI } from '../services/api';
import { TrainingRecord } from '../types/index';
import Card from '../components/common/Card';
import Table from '../components/common/Table';
import { Activity, TrendingUp, Clock } from 'lucide-react';
import { format } from 'date-fns';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';

const Training: React.FC = () => {
  const [stats, setStats] = useState<any>(null);
  const [records, setRecords] = useState<TrainingRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [filterType, setFilterType] = useState('');

  useEffect(() => {
    loadStats();
    loadRecords();
  }, [page, filterType]);

  const loadStats = async () => {
    try {
      const response = await adminAPI.getTrainingStats();
      if (response.code === 0) {
        setStats(response.data);
      }
    } catch (error) {
      console.error('加载统计失败:', error);
    }
  };

  const loadRecords = async () => {
    setLoading(true);
    try {
      const response = await adminAPI.getTrainingRecords({
        page,
        page_size: 20,
        type: filterType || undefined,
      });
      if (response.code === 0 && response.data) {
        setRecords(response.data.records || []);
        setTotal(response.data.total || 0);
      }
    } catch (error) {
      console.error('加载记录失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const getTypeName = (type: string) => {
    const types: Record<string, string> = {
      meditation: '正念冥想',
      airflow: '气流练习',
      exposure: '社会脱敏',
      practice: 'AI实战',
    };
    return types[type] || type;
  };

  const statCards = [
    {
      title: '总记录数',
      value: stats?.total_records || 0,
      icon: Activity,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
    },
    {
      title: '正念冥想',
      value: stats?.meditation_count || 0,
      icon: Activity,
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
    },
    {
      title: '气流练习',
      value: stats?.airflow_count || 0,
      icon: TrendingUp,
      color: 'text-cyan-600',
      bgColor: 'bg-cyan-50',
    },
    {
      title: '社会脱敏',
      value: stats?.exposure_count || 0,
      icon: Clock,
      color: 'text-orange-600',
      bgColor: 'bg-orange-50',
    },
  ];

  const trainingData = stats
    ? [
        { name: '正念冥想', value: stats.meditation_count || 0, color: '#8b5cf6' },
        { name: '气流练习', value: stats.airflow_count || 0, color: '#06b6d4' },
        { name: '社会脱敏', value: stats.exposure_count || 0, color: '#f59e0b' },
        { name: 'AI实战', value: stats.practice_count || 0, color: '#ec4899' },
      ]
    : [];

  const weeklyData = [
    { name: '周一', 训练: 120, 用户: 45 },
    { name: '周二', 训练: 132, 用户: 52 },
    { name: '周三', 训练: 101, 用户: 38 },
    { name: '周四', 训练: 134, 用户: 48 },
    { name: '周五', 训练: 90, 用户: 35 },
    { name: '周六', 训练: 230, 用户: 78 },
    { name: '周日', 训练: 210, 用户: 65 },
  ];

  const columns = [
    {
      key: 'type',
      title: '训练类型',
      render: (_: any, record: TrainingRecord) => (
        <span className="text-sm text-gray-900">{getTypeName(record.type)}</span>
      ),
    },
    {
      key: 'user',
      title: '用户',
      render: (_: any, record: TrainingRecord) => (
        <span className="text-sm text-gray-900">
          {record.user?.username || '未知用户'}
        </span>
      ),
    },
    {
      key: 'duration',
      title: '训练时长',
      dataIndex: 'duration' as keyof TrainingRecord,
      render: (value: number) => (
        <span className="text-sm text-gray-900">{Math.floor(value / 60)} 分钟</span>
      ),
    },
    {
      key: 'timestamp',
      title: '训练时间',
      dataIndex: 'timestamp' as keyof TrainingRecord,
      render: (value: string) => format(new Date(value), 'yyyy-MM-dd HH:mm'),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold text-gray-900">训练统计</h1>
        <p className="mt-1 text-sm text-gray-500">训练模块数据统计和管理</p>
      </div>

      {/* 统计卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {statCards.map((card, index) => {
          const Icon = card.icon;
          return (
            <Card key={index} shadow>
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-gray-600 mb-1">{card.title}</p>
                  <p className="text-2xl font-semibold text-gray-900">{card.value}</p>
                </div>
                <div className={`w-12 h-12 ${card.bgColor} rounded-lg flex items-center justify-center`}>
                  <Icon className={`w-6 h-6 ${card.color}`} />
                </div>
              </div>
            </Card>
          );
        })}
      </div>

      {/* 图表区域 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card title="周数据趋势" shadow>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={weeklyData}>
              <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
              <XAxis dataKey="name" stroke="#6b7280" />
              <YAxis stroke="#6b7280" />
              <Tooltip
                contentStyle={{
                  backgroundColor: 'white',
                  border: '1px solid #e5e7eb',
                  borderRadius: '6px',
                }}
              />
              <Legend />
              <Bar dataKey="训练" fill="#3b82f6" radius={[4, 4, 0, 0]} />
              <Bar dataKey="用户" fill="#8b5cf6" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </Card>

        <Card title="训练类型分布" shadow>
          <ResponsiveContainer width="100%" height={300}>
            <PieChart>
              <Pie
                data={trainingData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name} ${percent ? (percent * 100).toFixed(0) : 0}%`}
                outerRadius={100}
                fill="#8884d8"
                dataKey="value"
              >
                {trainingData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.color} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </Card>
      </div>

      {/* 数据表格 */}
      <Card shadow>
        <div className="mb-4">
          <select
            value={filterType}
            onChange={(e) => {
              setFilterType(e.target.value);
              setPage(1);
            }}
            className="px-3 py-2 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">全部类型</option>
            <option value="meditation">正念冥想</option>
            <option value="airflow">气流练习</option>
            <option value="exposure">社会脱敏</option>
            <option value="practice">AI实战</option>
          </select>
        </div>
        <Table
          columns={columns}
          dataSource={records}
          loading={loading}
          striped
          pagination={{
            current: page,
            pageSize: 20,
            total,
            onChange: (newPage) => setPage(newPage),
          }}
        />
      </Card>
    </div>
  );
};

export default Training;
