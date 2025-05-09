"use client";

import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { isAuthenticated, logout } from '@/lib/auth';

interface AdminLayoutProps {
  children: React.ReactNode;
}

export default function AdminLayout({ children }: AdminLayoutProps) {
  const router = useRouter();
  const [authenticated, setAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check authentication on client-side
    if (!isAuthenticated()) {
      router.push('/admin/login');
    } else {
      setAuthenticated(true);
    }
    setIsLoading(false);
  }, [router]);

  const handleLogout = () => {
    logout();
    router.push('/admin/login');
  };

  // Show nothing while checking authentication
  if (isLoading || !authenticated) {
    return null;
  }

  return (
    <div className="min-h-screen bg-blue-50">
      <nav className="bg-gradient-to-r from-blue-600 to-blue-800 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <Link href="/admin/dashboard" className="text-xl font-bold">
                  BlogKu Admin
                </Link>
              </div>
              <div className="ml-10 flex items-baseline space-x-4">
                <Link
                  href="/admin/dashboard"
                  className="px-3 py-2 rounded-md text-sm font-medium hover:bg-blue-500"
                >
                  Dashboard
                </Link>
                <Link
                  href="/admin/blogs"
                  className="px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700"
                >
                  Manage Blogs
                </Link>
                <Link
                  href="/"
                  className="px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-700"
                  target="_blank"
                >
                  View Site
                </Link>
              </div>
            </div>
            <button
              onClick={handleLogout}
              className="px-3 py-2 rounded-md text-sm font-medium hover:bg-blue-500"
            >
              Logout
            </button>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {children}
      </main>
    </div>
  );
}
