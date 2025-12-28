import React, { useEffect, useState } from 'react';
import { adminAPI } from '../../services/api';
import { User } from '../../types/index';
import Card from '../../components/common/Card';
import Table from '../../components/common/Table';
import Button from '../../components/form/Button';
import UserModal from './UserModal';
import { Plus, Edit, Trash2 } from 'lucide-react';
import { format } from 'date-fns';

const UserList: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  useEffect(() => {
    loadUsers();
  }, [page]);

  const loadUsers = async () => {
    setLoading(true);
    try {
      const response = await adminAPI.getUsers({
        page,
        page_size: 20,
      });
      if (response.code === 0 && response.data) {
        setUsers(response.data.users || []);
        setTotal(response.data.total || 0);
      }
    } catch (error) {
      console.error('加载用户失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setEditingUser(null);
    setModalVisible(true);
  };

  const handleEdit = (user: User) => {
    setEditingUser(user);
    setModalVisible(true);
  };

  const handleDelete = async (id: string) => {
    if (!confirm('确定要删除这个用户吗？此操作不可恢复！')) return;
    try {
      const response = await adminAPI.deleteUser(id);
      if (response.code === 0) {
        loadUsers();
      } else {
        alert(response.message || '删除失败');
      }
    } catch (error) {
      console.error('删除用户失败:', error);
      alert('删除失败，请重试');
    }
  };

  const handleModalClose = () => {
    setModalVisible(false);
    setEditingUser(null);
    loadUsers();
  };

  const columns = [
    {
      key: 'username',
      title: '用户名',
      dataIndex: 'username' as keyof User,
    },
    {
      key: 'email',
      title: '邮箱',
      dataIndex: 'email' as keyof User,
      render: (value: string) => value || '-',
    },
    {
      key: 'phone',
      title: '手机号',
      dataIndex: 'phone' as keyof User,
      render: (value: string) => value || '-',
    },
    {
      key: 'role',
      title: '角色',
      dataIndex: 'role' as keyof User,
      render: (value: string) => {
        const roleMap: Record<string, string> = {
          admin: '管理员',
          editor: '编辑',
          user: '普通用户',
        };
        return roleMap[value] || value;
      },
    },
    {
      key: 'status',
      title: '状态',
      dataIndex: 'status' as keyof User,
      render: (value: string) => (
        <span
          className={`px-2 py-1 rounded text-xs font-medium ${
            value === 'active'
              ? 'bg-green-100 text-green-800'
              : 'bg-gray-100 text-gray-800'
          }`}
        >
          {value === 'active' ? '活跃' : '非活跃'}
        </span>
      ),
    },
    {
      key: 'created_at',
      title: '注册时间',
      dataIndex: 'created_at' as keyof User,
      render: (value: string) => format(new Date(value), 'yyyy-MM-dd HH:mm'),
    },
    {
      key: 'actions',
      title: '操作',
      render: (_: any, record: User) => (
        <div className="flex items-center gap-2">
          <button
            onClick={() => handleEdit(record)}
            className="text-blue-600 hover:text-blue-700"
          >
            <Edit className="w-4 h-4" />
          </button>
          <button
            onClick={() => handleDelete(record.id)}
            className="text-red-600 hover:text-red-700"
          >
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-gray-900">用户管理</h1>
          <p className="mt-1 text-sm text-gray-500">管理系统用户信息</p>
        </div>
        <Button variant="primary" onClick={handleAdd}>
          <Plus className="w-4 h-4 mr-2" />
          新增用户
        </Button>
      </div>

      <Card shadow>
        <Table
          columns={columns}
          dataSource={users}
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

      {modalVisible && (
        <UserModal
          visible={modalVisible}
          user={editingUser}
          onClose={handleModalClose}
        />
      )}
    </div>
  );
};

export default UserList;

