"use client";

import React, { useState, useEffect, useRef } from "react";
import { Search, Send, Trash2 } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Card } from "@/components/ui/card";
import { useQuery } from "@tanstack/react-query";
import { GetChatSessions, GetChatHistory, CloseChatSession } from "@/lib/api/chat";

import { getWsConfig } from "@/lib/ws-config";

interface ChatMessage {
  role: string;
  message: string;
  created_at: string;
}

interface SessionEntry {
  uid: string;
  user_id: number;
  user_name: string;
  user_avatar: string;
  status: string;
  last_message: string;
  last_message_at: string;
  created_at: string;
  updated_at: string;
}

export default function AdminChatPage() {
  const [selectedSessionUid, setSelectedSessionUid] = useState<string | null>(null);
  const [selectedUser, setSelectedUser] = useState<SessionEntry | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
   const [inputValue, setInputValue] = useState("");
  const [searchTerm, setSearchTerm] = useState("");
  const ws = useRef<WebSocket | null>(null);
  const scrollRef = useRef<HTMLDivElement>(null);

  // GET ALL CHAT SESSIONS
  const { data: sessionData, refetch: refetchSessions } = useQuery({
    queryKey: ["adminChatSessions"],
    queryFn: GetChatSessions,
    refetchInterval: 5000,
  });

  const getSessions = (): SessionEntry[] => {
    const data = sessionData as any;
    if (data?.status === "success" && Array.isArray(data.data?.entries)) {
      return data.data.entries;
    }
    return [];
  };

   const SESSIONS = getSessions();
  const filteredSessions = SESSIONS.filter((session) => {
    const searchLower = searchTerm.toLowerCase();
    return (
      (session.user_name || "").toLowerCase().includes(searchLower) ||
      (session.last_message || "").toLowerCase().includes(searchLower)
    );
  });

  // HANDLE SESSION SELECTION
  useEffect(() => {
    if (selectedSessionUid) {
      const user = SESSIONS.find((s) => s.uid === selectedSessionUid) || null;
      setSelectedUser(user);
      loadHistory(selectedSessionUid);
      connectWs(selectedSessionUid);
    }
    return () => {
      if (ws.current) {
        ws.current.close();
        ws.current = null;
      }
    };
  }, [selectedSessionUid]);

  // LOAD CHAT HISTORY
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
      } else {
        setMessages([]);
      }
    } catch (e) {
      console.error("Failed to load chat history:", e);
      setMessages([]);
    }
  };

  // CONNECT WEBSOCKET
  const connectWs = async (uid: string) => {
    if (ws.current) {
      ws.current.close();
      ws.current = null;
    }
    try {
      const { token, wsUrl } = await getWsConfig();
      const socket = new WebSocket(`${wsUrl}/chat/${uid}?token=${token}`);
      ws.current = socket;
      socket.onopen = () => {
        console.log("✅ WS Admin connected");
      };
      socket.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          const senderRole = (data.sender_role || data.payload?.sender_role || "").toLowerCase();
          const content = data.payload?.content || data.content || data.message || "";
          if (!content) return;
          if (senderRole !== "admin") {
            setMessages((prev) => [
              ...prev,
              {
                role: "user",
                message: content,
                created_at: data.created_at || new Date().toISOString(),
              },
            ]);
          }
        } catch (e) {
          console.error("Failed to parse WS message:", e);
        }
      };
      socket.onclose = () => {
        console.log("WS Admin disconnected");
      };
      socket.onerror = (err) => {
        console.error("WS Admin error:", err);
      };
    } catch (e) {
      console.error("Failed to connect WS:", e);
    }
  };

  // SEND MESSAGE
  const handleSendMessage = () => {
    if (!inputValue.trim() || !selectedSessionUid) return;
    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      const payload = {
        type: "message",
        payload: {
          content: inputValue,
          message_type: "text",
        },
      };
      ws.current.send(JSON.stringify(payload));
      setMessages((prev) => [
        ...prev,
        {
          role: "admin",
          message: inputValue,
          created_at: new Date().toISOString(),
        },
      ]);
      setInputValue("");
    } else {
      console.warn("WebSocket not connected");
    }
  };

  // CLOSE SESSION
  const handleCloseSession = async () => {
    if (!selectedSessionUid) return;
    try {
      await CloseChatSession(selectedSessionUid);
      if (ws.current) {
        ws.current.close();
        ws.current = null;
      }
      refetchSessions();
      setSelectedSessionUid(null);
      setSelectedUser(null);
      setMessages([]);
    } catch (e) {
      console.error("Failed to close session:", e);
    }
  };

  // AUTO SCROLL
  useEffect(() => {
    scrollRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  return (
    <div className="p-8 w-full h-[calc(100vh-130px)] pr-5">
      <h1 className="text-2xl font-extrabold text-slate-900 tracking-tight">Customer Service</h1>
      <p className="text-slate-500 mb-5">Utility untuk menangani percakapan antara admin dan pelanggan.</p>
      <Card className="flex flex-row h-full w-full overflow-hidden border shadow-sm bg-white">
        <div className="w-[350px] lg:w-[380px] flex-shrink-0 border-r flex flex-col h-full overflow-hidden bg-white">
          <div className="p-4 border-b shrink-0">
            <h2 className="text-xl font-bold mb-4">Messages</h2>
             <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-slate-400" />
              <Input 
                placeholder="Cari pelanggan..." 
                className="pl-9" 
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
            </div>
          </div>
          <div className="flex-1 overflow-hidden">
            <ScrollArea className="h-full flex-1">
               <div className="p-2 space-y-1">
                {filteredSessions.length > 0 ? (
                  filteredSessions.map((session) => (
                    <div
                      key={session.uid}
                      onClick={() => setSelectedSessionUid(session.uid)}
                      className={`flex items-center gap-3 p-3 rounded-xl cursor-pointer transition-all ${
                        selectedSessionUid === session.uid ? "bg-slate-100" : "hover:bg-slate-50"
                      }`}
                    >
                      <Avatar className="h-12 w-12 border">
                        <AvatarImage src={session.user_avatar} alt={session.user_name} />
                        <AvatarFallback>
                          {(session.user_name || "P")[0]}
                        </AvatarFallback>
                      </Avatar>
                      <div className="flex-1 min-w-0">
                        <div className="flex justify-between items-center">
                          <span className="font-bold text-sm truncate">{session.user_name || "User"}</span>
                          <span className="text-[10px] text-slate-400">
                            {new Date(session.updated_at || session.created_at).toLocaleTimeString([], {
                              hour: "2-digit",
                              minute: "2-digit",
                            })}
                          </span>
                        </div>
                        <p className="text-xs text-slate-500 truncate">
                          {session.last_message || "Active chat..."}
                        </p>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="text-center py-10 text-slate-400">
                    <p className="text-sm">Tidak ada chat ditemukan</p>
                  </div>
                )}
              </div>
            </ScrollArea>
          </div>
        </div>
        {/* MAIN: Chat Area */}
        <div className="flex-1 min-w-0 flex flex-col h-full overflow-hidden bg-white">
          {selectedUser ? (
            <>
              {/* Header */}
              <div className="p-4 border-b flex justify-between items-center shrink-0">
                <div className="flex items-center gap-3">
                  <Avatar>
                    <AvatarImage src={selectedUser.user_avatar} alt={selectedUser.user_name} />
                    <AvatarFallback>
                      {(selectedUser.user_name || "U")[0]}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <p className="font-bold text-sm">{selectedUser.user_name || "User"}</p>
                    <p className="text-[11px] text-emerald-500 font-medium">{selectedUser.status}</p>
                  </div>
                </div>
                <div className="flex gap-1">
                  <Button variant="ghost" size="icon" onClick={handleCloseSession} title="Tutup Sesi">
                    <Trash2 className="h-4 w-4 text-red-500" />
                  </Button>
                </div>
              </div>
              {/* Messages */}
                            <div className="flex-1 overflow-hidden">
                <ScrollArea className="h-full bg-slate-50/50">
                  <div className="space-y-4 p-4">
                    {messages.map((msg, index) => {
                      const isAdmin = msg.role === "admin";
                      return (
                        <div key={index} className={`flex ${isAdmin ? "justify-end" : "justify-start"}`}>
                          <div
                            className={`max-w-[80%] p-3 rounded-2xl shadow-sm text-sm ${
                              isAdmin
                                ? "bg-teal-600 text-white rounded-tr-none"
                                : "bg-white border text-slate-700 rounded-tl-none"
                            }`}
                          >
                            {msg.message}
                          </div>
                        </div>
                      );
                    })}
                    <div ref={scrollRef} />
                  </div>
                </ScrollArea>
              </div>
              {/* Input */}
              <div className="p-4 border-t shrink-0">
                <div className="flex gap-2">
                  <Input
                    placeholder="Tulis pesan..."
                    className="flex-1"
                    value={inputValue}
                    onChange={(e) => setInputValue(e.target.value)}
                    onKeyDown={(e) => e.key === "Enter" && handleSendMessage()}
                  />
                  <Button className="bg-teal-600 hover:bg-teal-700" onClick={handleSendMessage}>
                    <Send className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            </>
          ) : (
            <div className="h-full flex items-center justify-center text-slate-500 flex-col gap-3">
              <Search className="h-10 w-10 text-slate-300" />
              <p>Pilih chat dari list untuk memulai percakapan</p>
            </div>
          )}
        </div>
      </Card>
    </div>
  );
}
