"use client";

import Cookies from "js-cookie";

// Cookie configuration
const AUTH_COOKIE_NAME = "authToken";
const COOKIE_OPTIONS = {
  expires: 7,
  path: "/",
  secure: process.env.NODE_ENV === "production",
  sameSite: "Lax" as const,
};

// Helper to determine if we're on the client side
const isClient = typeof window !== "undefined";

export async function login(email: string, password: string) {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1"}/auth/login`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      }
    );

    if (!response.ok) {
      throw new Error("Authentication failed");
    }

    const data = await response.json();

    // Only set cookie on client-side
    if (isClient) {
      Cookies.set(AUTH_COOKIE_NAME, data.token, COOKIE_OPTIONS || undefined);
      localStorage.setItem(AUTH_COOKIE_NAME, data.token);
    }

    return data;
  } catch (error) {
    console.error("Login error:", error);
    throw error;
  }
}

export function logout() {
  if (isClient) {
    Cookies.remove(AUTH_COOKIE_NAME, { path: "/" });
    localStorage.removeItem(AUTH_COOKIE_NAME);
  }
}

export function getAuthData() {
  if (isClient) {
    const token = getToken();
    if (token) {
      return { token };
    }
  }
  return null;
}

export function isAuthenticated() {
  return getToken() !== null;
}

export function getToken() {
  if (!isClient) return null;

  // Try cookie first, then localStorage
  const token =
    Cookies.get(AUTH_COOKIE_NAME) || localStorage.getItem(AUTH_COOKIE_NAME);
  return token || null;
}
