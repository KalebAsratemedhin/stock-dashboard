'use client';

import { useEffect, useRef, useState } from 'react';
import { WebSocketMessage } from '@/types';
import { useStockStore, useSalesStore, useUserEventsStore, useFinancialStore } from '@/lib/store';
import { StockQuote, Sale, UserEvent, FinancialMetric } from '@/types';

const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080';

export function useWebSocket() {
  const [isConnected, setIsConnected] = useState(false);
  const [messages, setMessages] = useState<WebSocketMessage[]>([]);
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | undefined>(undefined);
  const subscriptionsRef = useRef<Set<string>>(new Set());

  // Store hooks
  const { addQuote } = useStockStore();
  const { addSale } = useSalesStore();
  const { addEvent } = useUserEventsStore();
  const { addMetric } = useFinancialStore();

  useEffect(() => {
    function connect() {
      try {
        const ws = new WebSocket(`${WS_URL}/ws`);
        wsRef.current = ws;

        ws.onopen = () => {
          setIsConnected(true);
          // Resubscribe to all channels
          subscriptionsRef.current.forEach((channel) => {
            ws.send(JSON.stringify({ type: 'subscribe', channel }));
          });
        };

        ws.onmessage = (event) => {
          try {
            const message: WebSocketMessage = JSON.parse(event.data);
            
            // Handle subscription confirmations
            if (message.type === 'subscribed' || message.type === 'unsubscribed') {
              return;
            }

            // Handle data messages
            if (message.type === 'data' && message.channel && message.data) {
              setMessages((prev) => [...prev.slice(-99), message]);
              
              // Update stores based on channel
              switch (message.channel) {
                case 'stock_quotes':
                  if (message.data && typeof message.data === 'object') {
                    const quote = message.data as StockQuote;
                    addQuote(quote);
                  }
                  break;
                case 'sales':
                  if (message.data && typeof message.data === 'object') {
                    const sale = message.data as Sale;
                    addSale(sale);
                  }
                  break;
                case 'user_events':
                  if (message.data && typeof message.data === 'object') {
                    const event = message.data as UserEvent;
                    addEvent(event);
                  }
                  break;
                case 'financial_metrics':
                  if (message.data && typeof message.data === 'object') {
                    const metric = message.data as FinancialMetric;
                    addMetric(metric);
                  }
                  break;
              }
            }
          } catch (error) {
            console.error('Error parsing WebSocket message:', error);
          }
        };

        ws.onerror = (error) => {
          console.error('WebSocket error:', error);
        };

        ws.onclose = () => {
          setIsConnected(false);
          reconnectTimeoutRef.current = setTimeout(() => {
            connect();
          }, 3000);
        };
      } catch (error) {
        console.error('Error connecting WebSocket:', error);
        reconnectTimeoutRef.current = setTimeout(() => {
          connect();
        }, 3000);
      }
    }

    connect();

    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [addQuote, addSale, addEvent, addMetric]);

  const subscribe = (channel: string) => {
    subscriptionsRef.current.add(channel);
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type: 'subscribe', channel }));
    }
  };

  const unsubscribe = (channel: string) => {
    subscriptionsRef.current.delete(channel);
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type: 'unsubscribe', channel }));
    }
  };

  return {
    isConnected,
    messages,
    subscribe,
    unsubscribe,
  };
}
