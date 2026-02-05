'use client';

import { useEffect } from 'react';
import Layout from '../components/Layout';
import { PageHeader } from '../components/PageHeader';
import { useWebSocket } from '@/hooks/useWebSocket';
import { useUserEventsStore } from '@/lib/store';
import { fetchUserEvents } from '@/lib/api';
import { formatEventType } from '@/lib/utils';
import { format } from 'date-fns';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Activity, MousePointerClick, ShoppingCart, LogIn } from 'lucide-react';

const getEventIcon = (eventType: string) => {
  switch (eventType) {
    case 'page_view':
      return <MousePointerClick className="h-4 w-4" />;
    case 'purchase':
      return <ShoppingCart className="h-4 w-4" />;
    case 'login':
      return <LogIn className="h-4 w-4" />;
    default:
      return <Activity className="h-4 w-4" />;
  }
};

export default function EventsPage() {
  const { isConnected, subscribe } = useWebSocket();
  const { events, setEvents } = useUserEventsStore();

  useEffect(() => {
    subscribe('user_events');
  }, [subscribe]);

  useEffect(() => {
    const loadData = async () => {
      try {
        const end = new Date().toISOString();
        const start = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString();
        
        const eventsData = await fetchUserEvents(start, end);
        setEvents(Array.isArray(eventsData) ? eventsData : []);
      } catch (error) {
        console.error('Error loading events:', error);
      }
    };

    loadData();
  }, [setEvents]);

  const eventCounts = events.reduce((acc, event) => {
    const formattedType = formatEventType(event.event_type);
    acc[formattedType] = (acc[formattedType] || 0) + 1;
    return acc;
  }, {} as Record<string, number>);

  const recentEvents = events.slice(-20);

  return (
    <Layout>
      <div className="p-8 flex flex-col gap-8">
        <PageHeader
          title="User Behavior"
          subtitle="Real-Time User Activity Monitoring"
          isLive={isConnected}
        />

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {Object.entries(eventCounts).map(([type, count]) => (
            <Card key={type} className="bg-surface-slate border-border-slate">
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-[10px] font-bold text-slate-500 uppercase tracking-widest">{type}</CardTitle>
                <Activity className="h-5 w-5 text-cyber-orange" />
              </CardHeader>
              <CardContent>
                <div className="text-3xl font-black font-mono text-white">{count}</div>
              </CardContent>
            </Card>
          ))}
        </div>

        <Card className="bg-surface-slate border-border-slate overflow-hidden">
          <CardHeader className="px-6 py-4 border-b border-border-slate flex justify-between items-center bg-black/40">
            <div className="flex items-center gap-2">
              <span className="w-2 h-2 bg-cyber-orange rounded-full animate-pulse"></span>
              <CardTitle className="text-white text-xs font-black uppercase tracking-widest">Live Terminal Logs</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="p-0">
            <div className="overflow-x-auto font-mono">
              <Table>
                <TableHeader>
                  <TableRow className="text-slate-500 font-black border-b border-border-slate bg-black/20 hover:bg-transparent">
                    <TableHead className="px-6 py-3 uppercase tracking-[0.2em] text-[9px]">Timestamp</TableHead>
                    <TableHead className="px-6 py-3 uppercase tracking-[0.2em] text-[9px]">Event Type</TableHead>
                    <TableHead className="px-6 py-3 uppercase tracking-[0.2em] text-[9px]">User ID</TableHead>
                    <TableHead className="px-6 py-3 uppercase tracking-[0.2em] text-[9px]">Page</TableHead>
                    <TableHead className="px-6 py-3 uppercase tracking-[0.2em] text-[9px]">Country</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {recentEvents.length > 0 ? (
                    recentEvents.map((event) => (
                      <TableRow key={event.id} className="hover:bg-white/[0.03] transition-colors border-b border-border-slate/30">
                        <TableCell className="px-6 py-4 text-slate-400 font-mono text-xs">
                          {format(new Date(event.timestamp), 'HH:mm:ss.SSS')}
                        </TableCell>
                        <TableCell className="px-6 py-4">
                          <span className="flex items-center gap-2 font-black text-cyber-orange">
                            {getEventIcon(event.event_type)}
                            {formatEventType(event.event_type)}
                          </span>
                        </TableCell>
                        <TableCell className="px-6 py-4 text-slate-500 italic truncate max-w-[150px]">
                          {event.user_id || 'Anonymous'}
                        </TableCell>
                        <TableCell className="px-6 py-4 text-slate-400 truncate max-w-[200px]">
                          {event.page || 'N/A'}
                        </TableCell>
                        <TableCell className="px-6 py-4 text-cyan-accent font-bold">
                          {event.country || 'N/A'}
                        </TableCell>
                      </TableRow>
                    ))
                  ) : (
                    <TableRow>
                      <TableCell colSpan={5} className="px-6 py-8 text-center text-slate-500">
                        No events found
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      </div>
    </Layout>
  );
}
