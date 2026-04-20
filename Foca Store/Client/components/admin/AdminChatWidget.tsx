"use client";

import React, { useState, useEffect, useRef } from "react";
import { Search, Send, Trash2, MessageSquare, X, ChevronLeft, ChevronRight } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Card } from "@/components/ui/card";
import { useQuery } from "@tanstack/react-query";
import { GetChatSessions, GetChatHistory, CloseChatSession } from "@/lib/api/chat";
import { getWsConfig } from "@/lib/ws-config";
import { cn } from "@/lib/utils";

interface ChatMessage {
  role: string;
  message: string;
  created_at: string;
}

interface SessionEntry {
  uid: string;
  user_id: number;
  user_name: string;
  status: string;
  last_message: string;
  updated_at: string;
}

export function AdminChatWidget() {
  const [isOpen, setIsOpen] = useState(false);
  const [isListExpanded, setIsListExpanded] = useState(true); // Toggle sidebar dalam widget
  const [selectedSessionUid, setSelectedSessionUid] = useState<string | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [inputValue, setInputValue] = useState("");
  const ws = useRef<WebSocket | null>(null);
  const scrollRef = useRef<HTMLDivElement>(null);

  // ─── QUERY SESSIONS ───
  const { data: sessionData, refetch: refetchSessions } = useQuery({
    queryKey: ["adminChatSessions"],
    queryFn: GetChatSessions,
    refetchInterval: 5000,
    enabled: isOpen,
  });

  const SESSIONS: SessionEntry[] = sessionData?.data?.entries || [];
  const selectedUser = SESSIONS.find(s => s.uid === selectedSessionUid);

  // ─── CONNECTION LOGIC ───
  useEffect(() => {
    if (selectedSessionUid && isOpen) {
      loadHistory(selectedSessionUid);
      connectWs(selectedSessionUid);
    }
    return () => cleanupWs();
  }, [selectedSessionUid, isOpen]);

  const cleanupWs = () => {
    if (ws.current) {
      ws.current.close();
      ws.current = null;
    }
  };

  const loadHistory = async (uid: string) => {
    try {
      const hist = (await GetChatHistory(uid)) as any;
      if (hist?.status === "success") {
        setMessages(hist.data.entries.map((e: any) => ({
          role: e.sender_role.toLowerCase(),
          message: e.content,
          created_at: e.created_at,
        })));
      }
    } catch (e) { console.error(e); }
  };

  const connectWs = async (uid: string) => {
    cleanupWs();
    try {
      const { token, wsUrl } = await getWsConfig();
      const socket = new WebSocket(`${wsUrl}/chat/${uid}?token=${token}`);
      ws.current = socket;
      socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        const senderRole = (data.sender_role || data.payload?.sender_role || "").toLowerCase();
        const content = data.payload?.content || data.content || "";
        if (content && senderRole !== "admin") {
          setMessages(prev => [...prev, { role: "user", message: content, created_at: data.created_at || new Date().toISOString() }]);
        }
      };
    } catch (e) { console.error(e); }
  };

  const handleSendMessage = () => {
    if (!inputValue.trim() || !selectedSessionUid || !ws.current) return;
    ws.current.send(JSON.stringify({ type: "message", payload: { content: inputValue, message_type: "text" } }));
    setMessages(prev => [...prev, { role: "admin", message: inputValue, created_at: new Date().toISOString() }]);
    setInputValue("");
  };

  useEffect(() => { scrollRef.current?.scrollIntoView({ behavior: "smooth" }); }, [messages]);

  return (
    <div className="fixed bottom-6 right-6 z-[9999] flex flex-col items-end gap-3">
      {isOpen && (
        <Card className="w-[85vw] md:w-[700px] h-[550px] shadow-2xl flex overflow-hidden border-teal-100 animate-in slide-in-from-bottom-4">
          
          {/* SIDEBAR MINI (Daftar Chat) */}
          <div className={cn(
            "border-r bg-slate-50 transition-all duration-300 flex flex-col",
            isListExpanded ? "w-[250px]" : "w-[70px]"
          )}>
            <div className="p-3 border-b flex items-center justify-between bg-white">
              {isListExpanded && <span className="font-bold text-sm">Chats</span>}
              <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => setIsListExpanded(!isListExpanded)}>
                {isListExpanded ? <ChevronLeft size={16}/> : <ChevronRight size={16}/>}
              </Button>
            </div>
            
            <ScrollArea className="flex-1">
              <div className="p-2 space-y-2">
                {SESSIONS.map((session) => (
                  <div
                    key={session.uid}
                    onClick={() => setSelectedSessionUid(session.uid)}
                    className={cn(
                      "flex items-center gap-3 p-2 rounded-lg cursor-pointer transition-colors",
                      selectedSessionUid === session.uid ? "bg-teal-600 text-white" : "hover:bg-slate-200"
                    )}
                  >
                    <Avatar className="h-10 w-10 shrink-0 border">
                      <AvatarFallback className={selectedSessionUid === session.uid ? "text-teal-600" : ""}>
                        {session.user_name?.[0] || "U"}
                      </AvatarFallback>
                    </Avatar>
                    {isListExpanded && (
                      <div className="flex-1 min-w-0">
                        <p className="text-xs font-bold truncate">{session.user_name}</p>
                        <p className={cn("text-[10px] truncate opacity-80")}>{session.last_message}</p>
                      </div>
                    )}
                  </div>
                ))}
              </div>
            </ScrollArea>
          </div>

          {/* MAIN CHAT AREA */}
          <div className="flex-1 flex flex-col bg-white min-w-0">
            {selectedUser ? (
              <>
                <div className="p-3 border-b flex justify-between items-center bg-white shrink-0">
                  <div className="flex items-center gap-2 truncate">
                    <p className="font-bold text-sm truncate">{selectedUser.user_name}</p>
                    <div className="h-2 w-2 rounded-full bg-emerald-500" />
                  </div>
                  <Button variant="ghost" size="icon" className="h-8 w-8 text-red-400" onClick={() => {/* logic close session */}}>
                    <Trash2 size={14}/>
                  </Button>
                </div>

                <ScrollArea className="flex-1 bg-slate-50/50">
                  <div className="p-4 space-y-4">
                    {messages.map((msg, i) => (
                      <div key={i} className={cn("flex", msg.role === "admin" ? "justify-end" : "justify-start")}>
                        <div className={cn(
                          "max-w-[85%] p-2.5 rounded-2xl text-sm shadow-sm",
                          msg.role === "admin" ? "bg-teal-600 text-white rounded-tr-none" : "bg-white border text-slate-700 rounded-tl-none"
                        )}>
                          {msg.message}
                        </div>
                      </div>
                    ))}
                    <div ref={scrollRef} />
                  </div>
                </ScrollArea>

                <div className="p-3 border-t bg-white">
                  <div className="flex gap-2">
                    <Input 
                      placeholder="Balas chat..." 
                      value={inputValue}
                      onChange={(e) => setInputValue(e.target.value)}
                      onKeyDown={(e) => e.key === "Enter" && handleSendMessage()}
                      className="text-sm focus-visible:ring-teal-600"
                    />
                    <Button size="icon" onClick={handleSendMessage} className="bg-teal-600 hover:bg-teal-700 h-9 w-9 shrink-0">
                      <Send size={16} />
                    </Button>
                  </div>
                </div>
              </>
            ) : (
              <div className="flex-1 flex flex-col items-center justify-center text-slate-400 p-6 text-center">
                <MessageSquare className="h-12 w-12 mb-2 opacity-10" />
                <p className="text-sm">Pilih percakapan untuk membalas pesan pelanggan</p>
              </div>
            )}
          </div>
        </Card>
      )}

      {/* TRIGGER BUTTON */}
      <Button
        onClick={() => setIsOpen(!isOpen)}
        className={cn(
          "h-14 w-14 rounded-full shadow-2xl transition-all hover:scale-110",
          isOpen ? "bg-slate-800 text-white" : "bg-teal-600 text-white"
        )}
      >
        {isOpen ? <X size={24} /> : (
          <div className="relative">
            <MessageSquare size={24} />
            {SESSIONS.length > 0 && (
              <span className="absolute -top-2 -right-2 bg-red-500 text-[10px] w-5 h-5 rounded-full flex items-center justify-center border-2 border-white">
                {SESSIONS.length}
              </span>
            )}
          </div>
        )}
      </Button>
    </div>
  );
}