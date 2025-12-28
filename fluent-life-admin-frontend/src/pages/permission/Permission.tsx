import React, { useState } from 'react';
import Card from '../../components/common/Card';
import Table from '../../components/common/Table';
import Button from '../../components/form/Button';
import { Role, Menu } from '../../types/index';
import { mockRoles, mockMenus } from '../../mock/data';
import { Plus, Edit, Trash2, Shield, List } from 'lucide-react';

const Permission: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'roles' | 'menus'>('roles');
  const [roles] = useState<Role[]>(mockRoles);
  const [menus] = useState<Menu[]>(mockMenus);

  const roleColumns = [
    {
      key: 'name',
      title: '角色名称',
      dataIndex: 'name' as keyof Role,
    },
    {
      key: 'code',
      title: '角色代码',
      dataIndex: 'code' as keyof Role,
    },
    {
      key: 'description',
      title: '描述',
      dataIndex: 'description' as keyof Role,
      render: (value: string) => value || '-',
    },
    {
      key: 'permissions',
      title: '权限',
      render: (_: any, record: Role) => (
        <div className="flex flex-wrap gap-1">
          {record.permissions.map((perm) => (
            <span
              key={perm}
              className="px-2 py-1 bg-blue-50 text-blue-700 rounded text-xs"
            >
              {perm === '*' ? '全部权限' : perm}
            </span>
          ))}
        </div>
      ),
    },
    {
      key: 'actions',
      title: '操作',
      render: () => (
        <div className="flex items-center gap-2">
          <button className="text-blue-600 hover:text-blue-700">
            <Edit className="w-4 h-4" />
          </button>
          <button className="text-red-600 hover:text-red-700">
            <Trash2 className="w-4 h-4" />
          </button>
        </div>
      ),
    },
  ];

  const menuColumns = [
    {
      key: 'name',
      title: '菜单名称',
      dataIndex: 'name' as keyof Menu,
    },
    {
      key: 'path',
      title: '路径',
      dataIndex: 'path' as keyof Menu,
    },
    {
      key: 'icon',
      title: '图标',
      dataIndex: 'icon' as keyof Menu,
      render: (value: string) => value || '-',
    },
    {
      key: 'sort',
      title: '排序',
      dataIndex: 'sort' as keyof Menu,
    },
    {
      key: 'actions',
      title: '操作',
      render: () => (
        <div className="flex items-center gap-2">
          <button className="text-blue-600 hover:text-blue-700">
            <Edit className="w-4 h-4" />
          </button>
          <button className="text-red-600 hover:text-red-700">
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
          <h1 className="text-2xl font-semibold text-gray-900">权限管理</h1>
          <p className="mt-1 text-sm text-gray-500">管理系统角色和菜单权限</p>
        </div>
        <Button variant="primary">
          <Plus className="w-4 h-4 mr-2" />
          新增{activeTab === 'roles' ? '角色' : '菜单'}
        </Button>
      </div>

      <div className="flex gap-4 border-b border-gray-200">
        <button
          onClick={() => setActiveTab('roles')}
          className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
            activeTab === 'roles'
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-600 hover:text-gray-900'
          }`}
        >
          <Shield className="w-4 h-4 inline mr-2" />
          角色管理
        </button>
        <button
          onClick={() => setActiveTab('menus')}
          className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
            activeTab === 'menus'
              ? 'border-blue-600 text-blue-600'
              : 'border-transparent text-gray-600 hover:text-gray-900'
          }`}
        >
          <List className="w-4 h-4 inline mr-2" />
          菜单管理
        </button>
      </div>

      <Card shadow>
        {activeTab === 'roles' ? (
          <Table columns={roleColumns} dataSource={roles} striped />
        ) : (
          <Table columns={menuColumns} dataSource={menus} striped />
        )}
      </Card>
    </div>
  );
};

export default Permission;

