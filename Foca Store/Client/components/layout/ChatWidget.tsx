"use client";

import { Send, MessageCircle, X, Loader2 } from "lucide-react"; // Tambah icon baru
import { useState, useEffect, useRef } from "react";
import {
  CreateChatRequest,
  GetMySession,
  GetChatHistory,
} from "@/lib/api/chat";
import { getWsConfig } from "@/lib/ws-config";
import { useMutation, useQuery } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils"; // Pastikan shadcn utility ini ada

interface ChatMessage {
  role: string;
  message: string;
  created_at: string;
}

export default function UserChatWidget() {
  const [isOpen, setIsOpen] = useState(false); // State untuk buka/tutup widget
  const [inputValue, setInputValue] = useState("");
  const [activeSession, setActiveSession] = useState<string | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const ws = useRef<WebSocket | null>(null);
  const scrollRef = useRef<HTMLDivElement>(null);
  const activeSessionRef = useRef<string | null>(null);

  // ─── GET SESSION (polling) ───
  const { data: sessionData } = useQuery({
    queryKey: ["myChatSession"],
    queryFn: GetMySession,
    refetchInterval: 5000,
    enabled: isOpen, // Hanya polling jika widget terbuka
  });

  // ─── HANDLE SESSION CHANGE (Logic tetap sama) ───
  useEffect(() => {
    const data = sessionData as any;
    if (data?.status === "success" && data?.data?.uid && data?.data?.status === "active") {
      const uid = data.data.uid;
      if (activeSessionRef.current !== uid) {
        activeSessionRef.current = uid;
        setActiveSession(uid);
        loadHistory(uid);
        connectWs(uid);
      }
    } else if (activeSessionRef.current) {
      cleanupWs();
      activeSessionRef.current = null;
      setActiveSession(null);
      setMessages([]);
    }
  }, [sessionData]);

  const loadHistory = async (uid: string) => {
    try {
      const hist = (await GetChatHistory(uid)) as any;
      if (hist?.status === "success" && Array.isArray(hist.data?.entries)) {
        const mapped: ChatMessage[] = hist.data.entries.map((entry: any) => ({
          role: entry.sender_role.toLowerCase(),
          message: entry.content,
          created_at: entry.created_at,
        }));
        setMessages(mapped);
      }
    } catch (e) {
      console.error("Failed to load chat history:", e);
    }
  };

  const cleanupWs = () => {
    if (ws.current) {
      ws.current.close();
      ws.current = null;
    }
  };

  const connectWs = async (uid: string) => {
    cleanupWs();
    try {
      const { wsUrl, token } = await getWsConfig();
      const socket = new WebSocket(`${wsUrl}/chat/${uid}?token=${token}`);
      ws.current = socket;

      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          const senderRole = (data.sender_role || data.payload?.sender_role || "").toLowerCase();
          const content = data.payload?.content || data.content || data.message || "";

          if (!content) return;
          if (senderRole !== "user") {
            setMessages((prev) => [
              ...prev,
              {
                role: "admin",
                message: content,
                created_at: data.created_at || new Date().toISOString(),
              },
            ]);
          }
        } catch (e) { console.error(e); }
      };
    } catch (e) { console.error(e); }
  };

  const mutation = useMutation({
    mutationFn: (msg: string) => CreateChatRequest(msg),
    onSuccess: async (response: any) => {
      if (response.status === "success") {
        const uid = response.data.uid;
        activeSessionRef.current = uid;
        setActiveSession(uid);
        await connectWs(uid);
        setMessages([{
          role: "user",
          message: inputValue,
          created_at: new Date().toISOString(),
        }]);
        setInputValue("");
      }
    },
  });

  const handleSendMessage = () => {
    if (!inputValue.trim()) return;
    if (!activeSession) {
      mutation.mutate(inputValue);
      return;
    }

    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({
        type: "message",
        payload: { content: inputValue, message_type: "text" },
      }));
      setMessages((prev) => [...prev, {
        role: "user",
        message: inputValue,
        created_at: new Date().toISOString(),
      }]);
      setInputValue("");
    }
  };

  useEffect(() => {
    scrollRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, isOpen]);

  return (
    <div className="fixed bottom-6 right-6 z-50 flex flex-col items-end gap-4">
      {isOpen && (
        <Card className="w-[350px] sm:w-[400px] h-[500px] shadow-2xl flex flex-col overflow-hidden border-teal-100 animate-in slide-in-from-bottom-5 duration-300">
          {/* Header */}
          <div className="bg-teal-600 p-4 text-white flex justify-between items-center">
            <div>
              <h3 className="font-semibold">Chat Admin</h3>
              <p className="text-xs text-teal-100">Kami siap membantu Anda</p>
            </div>
            <Button 
              variant="ghost" 
              size="icon" 
              onClick={() => setIsOpen(false)}
              className="text-white hover:bg-teal-700"
            >
              <X className="h-5 w-5" />
            </Button>
          </div>

          <div className="flex-1 overflow-hidden bg-slate-50">
            <ScrollArea className="h-full">
              <div className="p-4 space-y-4">
                {messages.length === 0 && !activeSession && (
                  <div className="flex flex-col items-center justify-center h-40 text-slate-400 text-center px-4">
                    <MessageCircle className="h-8 w-8 mb-2 opacity-20" />
                    <p className="text-sm">Halo! Ada yang bisa kami bantu? Silakan kirim pesan.</p>
                  </div>
                )}

                {messages.map((msg, i) => {
                  const isUser = msg.role === "user";
                  return (
                    <div key={i} className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
                      <div className={cn(
                        "p-3 rounded-2xl text-sm max-w-[85%] shadow-sm",
                        isUser 
                          ? "bg-teal-600 text-white rounded-tr-none" 
                          : "bg-white border text-slate-700 rounded-tl-none"
                      )}>
                        {msg.message}
                      </div>
                    </div>
                  );
                })}
                <div ref={scrollRef} />
              </div>
            </ScrollArea>
          </div>

          {/* Input Area */}
          <div className="p-4 border-t bg-white flex gap-2">
            <Input
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSendMessage()}
              placeholder="Tulis pesan..."
              disabled={mutation.isPending}
              className="focus-visible:ring-teal-600"
            />
            <Button 
              onClick={handleSendMessage} 
              disabled={mutation.isPending}
              className="bg-teal-600 hover:bg-teal-700"
            >
              {mutation.isPending ? <Loader2 className="animate-spin h-4 w-4" /> : <Send className="h-4 w-4" />}
            </Button>
          </div>
        </Card>
      )}

      {/* FLOATING TRIGGER BUTTON */}
      <Button
        onClick={() => setIsOpen(!isOpen)}
        className={cn(
          "h-14 w-14 rounded-full shadow-lg transition-all duration-300 transform hover:scale-110",
          isOpen ? "bg-slate-200 text-slate-600 hover:bg-slate-300" : "bg-teal-600 text-white hover:bg-teal-700"
        )}
      >
        {isOpen ? <X className="h-6 w-6" /> : <MessageCircle className="h-6 w-6" />}
      </Button>
    </div>
  );
}