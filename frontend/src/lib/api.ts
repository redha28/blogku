import { Blog } from "@/types/blog";
import Cookies from 'js-cookie';
import { getToken } from './auth';

const API_BASE_URL = "http://localhost:8080/api/v1";

// Common fetch options to include credentials
const getCommonOptions = () => ({
  credentials: 'include' as RequestCredentials,
  headers: {
    'Content-Type': 'application/json'
  }
});

// Helper to get auth headers
const getAuthHeaders = (token?: string) => {
  const authToken = token || getToken();
  return authToken ? {
    'Authorization': `Bearer ${authToken}`
  } : {};
};

// Helper to combine headers
const combineHeaders = (token?: string) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json'
  };
  
  const authToken = token || getToken();
  if (authToken) {
    headers['Authorization'] = `Bearer ${authToken}`;
  }
  
  return headers;
};

export async function login(email: string, password: string) {
  try {
    console.log("Sending login request to:", `${API_BASE_URL}/auth/login`);
    
    const response = await fetch(`${API_BASE_URL}/auth/login`, {
      method: "POST",
      ...getCommonOptions(),
      body: JSON.stringify({ email, password })
    });
    
    if (!response.ok) {
      const errorText = await response.text();
      console.error("Login failed:", errorText);
      throw new Error("Authentication failed");
    }
    
    const data = await response.json();
    console.log("Login response:", { ...data, token: data.token ? "TOKEN_EXISTS" : "NO_TOKEN" });
    
    // Ensure we have a token
    if (!data.token) {
      console.error("No token in response");
      throw new Error("No authentication token received");
    }
    
    return data;
  } catch (error) {
    console.error("API error during login:", error);
    throw error;
  }
}

export async function getBlogs() {
  const response = await fetch(`${API_BASE_URL}/blogs`, getCommonOptions());
  
  if (!response.ok) {
    throw new Error("Failed to fetch blogs");
  }
  
  return await response.json();
}

export async function getBlogBySlug(slug: string) {
  const response = await fetch(`${API_BASE_URL}/blogs/${slug}`, getCommonOptions());
  
  if (!response.ok) {
    throw new Error("Failed to fetch blog");
  }
  
  return await response.json();
}

export async function getAdminBlogs(token?: string) {
  const authToken = token || getToken();
  console.log("Fetching admin blogs with token:", authToken ? authToken.substring(0, 10) + "..." : "NO_TOKEN");
  
  try {
    const response = await fetch(`${API_BASE_URL}/blogs`, {
      ...getCommonOptions(),
      headers: combineHeaders(authToken || undefined)
    });
    
    if (!response.ok) {
      const errorText = await response.text();
      console.error("Admin blogs fetch failed:", errorText);
      throw new Error("Failed to fetch admin blogs");
    }
    
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("API error when fetching admin blogs:", error);
    throw error;
  }
}

export async function createBlog(token: string | null, blog: any) {
  // Use FormData for multipart/form-data to support file uploads
  const formData = new FormData();
  
  // Append text fields
  formData.append('title', blog.title);
  formData.append('content', blog.content);
  
  // Append image file if it exists
  if (blog.image) {
    formData.append('image', blog.image);
  }
  
  const response = await fetch(`${API_BASE_URL}/admin/blogs`, {
    method: "POST",
    // Don't include Content-Type for FormData - browser sets it automatically with boundary
    headers: {
      Authorization: `Bearer ${token}`
    },
    credentials: "include",
    body: formData
  });
  
  if (!response.ok) {
    const errorText = await response.text();
    console.error("Create blog failed:", errorText);
    throw new Error("Failed to create blog");
  }
  
  return await response.json();
}

export async function updateBlog(token: string | null, id: number, blog: any) {
  const response = await fetch(`${API_BASE_URL}/admin/blogs/${id}`, {
    method: "PATCH", // Use PATCH instead of PUT if your backend supports it
    ...getCommonOptions(),
    headers: combineHeaders(token || undefined),
    body: JSON.stringify(blog)
  });
  
  if (!response.ok) {
    const errorText = await response.text();
    console.error("Update failed:", errorText);
    throw new Error("Failed to update blog");
  }
  
  return await response.json();
}

export async function deleteBlog(token: string | null, id: number) {
  const response = await fetch(`${API_BASE_URL}/admin/blogs/${id}`, {
    method: "DELETE",
    ...getCommonOptions(),
    headers: combineHeaders(token || undefined)
  });
  
  if (!response.ok) {
    throw new Error("Failed to delete blog");
  }
  
  return await response.json();
}
