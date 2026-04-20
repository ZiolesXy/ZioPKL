import api from "../axios";
import { clientApi } from "../client-api"

export async function CreateChatRequest(message: string): Promise<any> {
    try {
        const response = await clientApi.post(
            "/api/chat/requests", 
            { message },
        );
        
        return response;
    } catch (error: unknown) {
        throw error;
    }
}

export async function GetMySession(): Promise<any> {
    try {
        const response = await clientApi.get("/api/chat/status");
        return response;
    } catch (error: unknown) {
        throw error;
    }
}

export async function GetChatSessions(): Promise<any> {
    try {
        const response = await clientApi.get("/api/admin/chat/sessions");
        return response
    } catch (error: unknown) {
        throw error;
    }
}

export async function GetChatHistory(session_uid: string): Promise<any> {
    try {
        const response = await clientApi.get(`/api/chat/sessions/${session_uid}/messages`);
        return response;
    } catch (error: unknown) {
        throw error;
    }
}

export async function CloseChatSession(session_uid: string): Promise<any> {
    try {
        const response = await clientApi.delete(`/api/chat/sessions/${session_uid}/close`);
        return response;
    } catch (error: unknown) {
        throw error;
    }
}

export async function SearchChat(keyword: string){
    try {
        const response = await api.get(`/api/chat/search?keyword=${keyword}`);
        return response.data;
    } catch (error: unknown) {
        throw error;
    }
}