'use client';

import { ReactNode } from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { cn } from '@/lib/utils';
import { BarChart3, TrendingUp, ShoppingCart, Activity, Wallet, Search, Settings, Bell } from 'lucide-react';

interface LayoutProps {
  children: ReactNode;
}

const getNavIcon = (iconName: string) => {
  switch (iconName) {
    case 'dashboard':
      return <BarChart3 className="h-5 w-5" />;
    case 'trending_up':
      return <TrendingUp className="h-5 w-5" />;
    case 'shopping_cart':
      return <ShoppingCart className="h-5 w-5" />;
    case 'sensors':
      return <Activity className="h-5 w-5" />;
    case 'account_balance_wallet':
      return <Wallet className="h-5 w-5" />;
    default:
      return <BarChart3 className="h-5 w-5" />;
  }
};

export default function Layout({ children }: LayoutProps) {
  const pathname = usePathname();

  const navItems = [
    { href: '/', label: 'Dashboard', icon: 'dashboard' },
    { href: '/stocks', label: 'Market Intelligence', icon: 'trending_up' },
    { href: '/sales', label: 'Sales Radar', icon: 'shopping_cart' },
    { href: '/events', label: 'User Behavior', icon: 'sensors' },
    { href: '/financial', label: 'Financials', icon: 'account_balance_wallet' },
  ];

  return (
    <div className="flex h-screen overflow-hidden bg-midnight-slate text-node-white">
      <aside className="w-64 flex flex-col border-r border-border-slate bg-midnight-slate p-4 shrink-0">
        <div className="flex items-center gap-3 mb-10">
          <div className="bg-cyber-orange rounded-lg p-2 text-midnight-slate flex items-center justify-center">
            <BarChart3 className="h-6 w-6" />
          </div>
          <div>
            <h1 className="text-base font-black tracking-tighter uppercase italic text-white">
              Nexus<span className="text-cyber-orange">Analytics</span>
            </h1>
            <p className="text-slate-500 text-[10px] font-bold uppercase tracking-widest">Real-Time Dashboard</p>
          </div>
        </div>

        <nav className="flex-1 flex flex-col gap-1">
          {navItems.map((item) => {
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all",
                  isActive
                    ? 'bg-cyber-orange text-midnight-slate shadow-lg shadow-cyber-orange/20'
                    : 'text-slate-400 hover:bg-surface-slate hover:text-cyber-orange'
                )}
              >
                <span className="flex-shrink-0">{getNavIcon(item.icon)}</span>
                <span className="text-xs font-bold uppercase tracking-wider">{item.label}</span>
              </Link>
            );
          })}
        </nav>

        
      </aside>

      <div className="flex-1 flex flex-col overflow-hidden">
        <header className="flex items-center justify-end border-b border-border-slate bg-midnight-slate px-6 py-3 shrink-0">
          <div className="flex items-center gap-6">
            <div className="flex items-center bg-surface-slate rounded-full px-4 py-2 gap-3 border border-border-slate focus-within:border-cyber-orange focus-within:ring-1 focus-within:ring-cyber-orange transition-all group">
              <Search className="h-5 w-5 text-slate-500 group-focus-within:text-cyber-orange" />
              <input
                className="bg-transparent border-none focus:ring-0 text-sm w-56 text-white placeholder:text-slate-600 p-0"
                placeholder="SEARCH..."
                type="text"
              />
            </div>
            <div className="flex gap-1">
              <button className="p-2 rounded-lg text-slate-500 hover:text-cyber-orange transition-all" aria-label="Settings">
                <Settings className="h-5 w-5" />
              </button>
              <button className="p-2 rounded-lg text-slate-500 hover:text-cyber-orange transition-all relative" aria-label="Notifications">
                <Bell className="h-5 w-5" />
                <span className="absolute top-2 right-2 w-2 h-2 bg-cyber-orange rounded-full"></span>
              </button>
            </div>
          </div>
        </header>

        <main className="flex-1 overflow-y-auto bg-midnight-slate scrollbar-thin">
          {children}
        </main>
      </div>
    </div>
  );
}
