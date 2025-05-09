"use client";

import AdminLayout from '@/components/admin/AdminLayout';
import { getAdminBlogs } from '@/lib/api';
import { getToken, isAuthenticated } from '@/lib/auth';
import { Blog } from '@/types/blog';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

export default function AdminDashboardPage() {
  const router = useRouter();
  const [blogs, setBlogs] = useState<Blog[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [authChecked, setAuthChecked] = useState(false);

  useEffect(() => {
    // First, check if user is authenticated
    const checkAuth = () => {
      if (!isAuthenticated()) {
        console.log("Not authenticated, redirecting to login");
        // Uncomment this line to enable redirection
        router.replace('/admin/login');
        return false;
      }
      return true;
    };

    if (!authChecked) {
      const isAuth = checkAuth();
      setAuthChecked(isAuth);
      if (!isAuth) return;
    }

    async function fetchBlogs() {
      try {
        const token = getToken();
        console.log("Token found:", token ? "Yes" : "No");
        
        if (!token) {
          console.error("No token available");
          setError('Authentication error: No token available. Please log in again.');
          return;
        }
        
        const data = await getAdminBlogs(token);
        console.log("Blogs fetched successfully");
        setBlogs(data.blogs || data); // Handle either response format
      } catch (err) {
        console.error("Error fetching blogs:", err);
        setError('Failed to load blogs. Please refresh or try logging in again.');
      } finally {
        setLoading(false);
      }
    }

    if (authChecked) {
      fetchBlogs();
    }
  }, [authChecked]);

  // If still checking auth or not authenticated, show loading
  if (!authChecked) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-blue-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500 mx-auto"></div>
          <p className="mt-4 text-gray-600">Checking authentication...</p>
        </div>
      </div>
    );
  }

  return (
    <AdminLayout>
      <div className="px-4 py-6 sm:px-0">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-semibold text-gray-900">Admin Dashboard</h1>
          <Link
            href="/admin/blogs/new"
            className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
          >
            Create New Blog
          </Link>
        </div>

        <div className="bg-white shadow overflow-hidden sm:rounded-lg mb-6">
          <div className="px-4 py-5 sm:px-6">
            <h2 className="text-lg font-medium text-gray-900">Overview</h2>
          </div>
          <div className="border-t border-gray-200 px-4 py-5 sm:px-6">
            <dl className="grid grid-cols-1 gap-x-4 gap-y-8 sm:grid-cols-3">
              <div className="sm:col-span-1">
                <dt className="text-sm font-medium text-gray-500">Total Blogs</dt>
                <dd className="mt-1 text-3xl font-semibold text-gray-900">
                  {loading ? '...' : blogs.length}
                </dd>
              </div>
              <div className="sm:col-span-1">
                <dt className="text-sm font-medium text-gray-500">Published Blogs</dt>
                <dd className="mt-1 text-3xl font-semibold text-gray-900">
                  {loading ? '...' : blogs.length}
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <div className="bg-white shadow overflow-hidden sm:rounded-lg">
          <div className="px-4 py-5 sm:px-6">
            <h2 className="text-lg font-medium text-gray-900">Recent Blogs</h2>
          </div>

          {loading ? (
            <div className="px-4 py-5 sm:px-6">Loading...</div>
          ) : error ? (
            <div className="px-4 py-5 sm:px-6 text-red-500">{error}</div>
          ) : blogs.length === 0 ? (
            <div className="px-4 py-5 sm:px-6">No blogs found</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200">
                <thead>
                  <tr>
                    <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Title
                    </th>
                    <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Slug
                    </th>
                    <th className="px-6 py-3 bg-gray-50 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Published Date
                    </th>
                    <th className="px-6 py-3 bg-gray-50 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="bg-white divide-y divide-gray-200">
                  {blogs.slice(0, 5).map((blog) => (
                    <tr key={blog.id}>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                        {blog.title}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {blog.slug}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                        {blog.published_at}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                        <Link
                          href={`/admin/blogs/${blog.slug}`}
                          className="text-blue-600 hover:text-blue-900 mr-4"
                        >
                          Edit
                        </Link>
                        <Link
                          href={`/blog/${blog.slug}`}
                          target="_blank"
                          className="text-green-600 hover:text-green-900"
                        >
                          View
                        </Link>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {blogs.length > 5 && (
            <div className="px-4 py-3 bg-gray-50 text-right sm:px-6">
              <Link
                href="/admin/blogs"
                className="text-sm font-medium text-blue-600 hover:text-blue-500"
              >
                View all blogs →
              </Link>
            </div>
          )}
        </div>
      </div>
    </AdminLayout>
  );
}
