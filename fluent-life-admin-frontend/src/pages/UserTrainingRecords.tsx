import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { adminAPI } from '../services/api';
import { TrainingRecord, User } from '../types/index';
import Card from '../components/common/Card';
import Table from '../components/common/Table';
import { format, parseISO } from 'date-fns';
import { ArrowLeft } from 'lucide-react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';

const UserTrainingRecords: React.FC = () => {
  const { userId } = useParams<{ userId: string }>();
  const navigate = useNavigate();
  const [records, setRecords] = useState<TrainingRecord[]>([]);
const [loading, setLoading] = useState(false);
const [page, setPage] = useState(1);
const [total, setTotal] = useState(0);
const [user, setUser] = useState<User | null>(null);
const [dailyDurationData, setDailyDurationData] = useState<any[]>([]);
const [trainingTypeData, setTrainingTypeData] = useState<any[]>([]);
  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

  useEffect(() => {
    if (records.length === 0) return;

    // Process daily duration data for line chart
    const dailyMap = new Map<string, number>();
    records.forEach(record => {
      const date = format(parseISO(record.timestamp), 'yyyy-MM-dd');
      dailyMap.set(date, (dailyMap.get(date) || 0) + record.duration);
    });
    const sortedDailyData = Array.from(dailyMap.entries())
      .map(([date, duration]) => ({ date, duration }))
      .sort((a, b) => a.date.localeCompare(b.date));
    setDailyDurationData(sortedDailyData);

    // Process training type data for pie chart
    const typeMap = new Map<string, number>();
    const typeMapCN = { meditation: '冥想', airflow: '气流', exposure: '暴露', practice: '练习' };
    records.forEach(record => {
      typeMap.set(record.type, (typeMap.get(record.type) || 0) + 1);
    });
    const typeData = Array.from(typeMap.entries())
      .map(([type, value]) => ({ name: typeMapCN[type as keyof typeof typeMapCN] || type, value }));
    setTrainingTypeData(typeData);
  }, [records]);

  const loadUser = async () => {
    if (!userId) return;
    try {
      const response = await adminAPI.getUser(userId);
      if (response.code === 0 && response.data) {
        setUser(response.data);
      }
    } catch (error) {
      console.error('加载用户详情失败:', error);
    }
  };

  useEffect(() => {
    loadUser();
  }, [userId]);

  useEffect(() => {
    if (userId) {
      loadUserTrainingRecords();
    }
  }, [userId, page]);

  const loadUserTrainingRecords = async () => {
    setLoading(true);
    try {
      const response = await adminAPI.getTrainingRecords({
        page,
        page_size: 20,
        user_id: userId,
      });
      if (response.code === 0 && response.data) {
        setRecords(response.data.records || []);
        setTotal(response.data.total || 0);
      }
    } catch (error) {
      console.error('加载用户训练记录失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const columns = [
    {
      key: 'type',
      title: '训练类型',
      dataIndex: 'type' as keyof TrainingRecord,
      render: (value: string) => {
        const typeMap: { [key: string]: string } = {
          meditation: '冥想',
          airflow: '气流',
          exposure: '暴露',
          practice: '练习',
        };
        return typeMap[value] || value;
      },
    },
    {
      key: 'duration',
      title: '时长 (秒)',
      dataIndex: 'duration' as keyof TrainingRecord,
    },
    {
      key: 'timestamp',
      title: '训练时间',
      dataIndex: 'timestamp' as keyof TrainingRecord,
      render: (value: string) => format(new Date(value), 'yyyy-MM-dd HH:mm:ss'),
    },
    {
      key: 'data',
      title: '详细数据',
      dataIndex: 'data' as keyof TrainingRecord,
      render: (value: any) => (
        <pre className="text-xs bg-gray-100 p-2 rounded overflow-auto max-h-20">
          {JSON.stringify(value, null, 2)}
        </pre>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <div>
        <div className="flex items-center justify-between">
        <button
          onClick={() => navigate(-1)}
          className="text-blue-600 hover:text-blue-700 flex items-center"
        >
          <ArrowLeft className="w-4 h-4 mr-2" />
          返回
        </button>
        <h1 className="text-2xl font-semibold text-gray-900">{user ? `${user.username} 的训练记录` : '用户训练记录'}</h1>
      </div>
      <p className="mt-1 text-sm text-gray-500">用户ID: {userId}</p>
      </div>

      <Card shadow className="mb-6">
        <h2 className="text-xl font-semibold mb-4">每日训练时长</h2>
        <ResponsiveContainer width="100%" height={200}>
          <LineChart data={dailyDurationData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Line type="monotone" dataKey="duration" stroke="#8884d8" activeDot={{ r: 8 }} name="时长 (秒)" />
          </LineChart>
        </ResponsiveContainer>
      </Card>

      <Card shadow className="mb-6">
        <h2 className="text-xl font-semibold mb-4">训练类型分布</h2>
        <ResponsiveContainer width="100%" height={200}>
          <PieChart>
            <Pie
              data={trainingTypeData}
              cx="50%"
              cy="50%"
              labelLine={false}
              outerRadius={80}
              fill="#8884d8"
              dataKey="value"
              label={({ name, percent }) => `${name} ${((percent || 0) * 100).toFixed(0)}%`}
            >
              {trainingTypeData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip />
            <Legend />
          </PieChart>
        </ResponsiveContainer>
      </Card>

      <Card shadow>
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

export default UserTrainingRecords;
