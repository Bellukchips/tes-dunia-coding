import axios from "axios";

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export interface Task {
    id: number;
    title: string;
    description: string;
    due_date: string;
    status: 'pending' | 'in-progress' | 'completed';
    created_at?: string;
    updated_at?: string;
}

export interface TaskInput {
    title: string;
    description: string;
    due_date: string;
    status: 'pending' | 'in-progress' | 'completed';
}

export interface ApiResponse<T> {
    data?: T;
    success: boolean;
    error?: string;
}

export async function getAllTasks(): Promise<{ tasks: Task[]; success: boolean }> {
    try {
        const response = await axios.get<Task[]>(`${API_URL}/tasks`);
        return { tasks: response.data || [], success: true };
    } catch (error) {
        console.error("Error fetching tasks:", error);
        return { tasks: [], success: false };
    }
}

export async function getTaskById(id: number): Promise<{ task: Task | null; success: boolean }> {
    try {
        const response = await axios.get<Task>(`${API_URL}/tasks/${id}`);
        return { task: response.data, success: true };
    } catch (error) {
        console.error(`Error fetching task ${id}:`, error);
        return { task: null, success: false };
    }
}

export async function createTask(task: TaskInput): Promise<{ task: Task | null; success: boolean; error?: string }> {
    try {
        const response = await axios.post<Task>(`${API_URL}/tasks`, task);
        return { task: response.data, success: true };
    } catch (error: any) {
        console.error("Error creating task:", error);
        return { 
            task: null, 
            success: false, 
            error: error.response?.data?.error || 'Failed to create task'
        };
    }
}

export async function updateTask(id: number, task: Partial<TaskInput>): Promise<{ task: Task | null; success: boolean; error?: string }> {
    try {
        const response = await axios.put<Task>(`${API_URL}/tasks/${id}`, task);
        return { task: response.data, success: true };
    } catch (error: any) {
        console.error(`Error updating task ${id}:`, error);
        return { 
            task: null, 
            success: false, 
            error: error.response?.data?.error || 'Failed to update task'
        };
    }
}

export async function deleteTask(id: number): Promise<{ success: boolean; error?: string }> {
    try {
        await axios.delete(`${API_URL}/tasks/${id}`);
        return { success: true };
    } catch (error: any) {
        console.error(`Error deleting task ${id}:`, error);
        return { 
            success: false, 
            error: error.response?.data?.error || 'Failed to delete task'
        };
    }
}