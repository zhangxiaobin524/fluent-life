import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import {
  LayoutDashboard,
  Users,
  FileText,
  Home,
  BarChart3,
  Shield,
  Settings,
} from 'lucide-react';
import clsx from 'clsx';

interface MenuItem {
  key: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  path: string;
  children?: MenuItem[];
}

const menuItems: MenuItem[] = [
  {
    key: 'dashboard',
    label: '数据概览',
    icon: LayoutDashboard,
    path: '/',
  },
  {
    key: 'users',
    label: '用户管理',
    icon: Users,
    path: '/users',
  },
  {
    key: 'posts',
    label: '帖子管理',
    icon: FileText,
    path: '/posts',
  },
  {
    key: 'rooms',
    label: '房间管理',
    icon: Home,
    path: '/rooms',
  },
  {
    key: 'training',
    label: '训练统计',
    icon: BarChart3,
    path: '/training',
  },
  {
    key: 'permission',
    label: '权限管理',
    icon: Shield,
    path: '/permission',
  },
  {
    key: 'settings',
    label: '系统设置',
    icon: Settings,
    path: '/settings',
  },
];

interface SidebarProps {
  collapsed?: boolean;
}

const Sidebar: React.FC<SidebarProps> = ({ collapsed = false }) => {
  const location = useLocation();

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/';
    }
    return location.pathname.startsWith(path);
  };

  return (
    <aside
      className={clsx(
        'bg-white border-r border-gray-200 h-screen fixed left-0 top-0 z-40 transition-all',
        collapsed ? 'w-64' : 'w-64'
      )}
    >
      <div className="h-16 flex items-center px-6 border-b border-gray-200">
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 bg-blue-600 rounded flex items-center justify-center">
            <span className="text-white font-bold text-sm">F</span>
          </div>
          {!collapsed && (
            <div>
              <div className="text-sm font-semibold text-gray-900">流畅人生</div>
              <div className="text-xs text-gray-500">管理后台</div>
            </div>
          )}
        </div>
      </div>
      <nav className="p-4">
        {menuItems.map((item) => {
          const Icon = item.icon;
          const active = isActive(item.path);
          return (
            <Link
              key={item.key}
              to={item.path}
              className={clsx(
                'flex items-center gap-3 px-3 py-2.5 rounded text-sm font-medium mb-1 transition-colors',
                {
                  'bg-blue-50 text-blue-600': active,
                  'text-gray-700 hover:bg-gray-50': !active,
                }
              )}
            >
              <Icon className="w-5 h-5 flex-shrink-0" />
              {!collapsed && <span>{item.label}</span>}
            </Link>
          );
        })}
      </nav>
    </aside>
  );
};

export default Sidebar;

