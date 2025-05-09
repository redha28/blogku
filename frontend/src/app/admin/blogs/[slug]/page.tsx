"use client";

import AdminLayout from '@/components/admin/AdminLayout';
import LoadingSpinner from '@/components/admin/LoadingSpinner';
import { getAdminBlogs, updateBlog } from '@/lib/api';
import { getToken } from '@/lib/auth';
import { Blog } from '@/types/blog';
import { useParams, useRouter } from 'next/navigation';
import { FormEvent, useEffect, useState } from 'react';
import toast from 'react-hot-toast';

// Correct interface for Next.js App Router page component
// type PageProps = {
//   params: { slug: string };
//   searchParams?: { [key: string]: string | string[] | undefined };
// }

export default function EditBlogPage() {
  const slug = useParams().slug;
  const router = useRouter();
  const [blog, setBlog] = useState<Blog | null>(null);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    async function fetchBlog() {
      try {
        setLoading(true);
        const token = getToken();
        if (!token) {
          router.push('/admin/login');
          return;
        }

        const blogs = await getAdminBlogs(token);
        const blogData = Array.isArray(blogs.blogs) ? blogs.blogs : blogs;
        const foundBlog = blogData.find((b: Blog) => b.slug === slug);
        
        if (!foundBlog) {
          setError('Blog post not found');
          return;
        }

        setBlog(foundBlog);
        setTitle(foundBlog.title);
        setContent(foundBlog.content);
      } catch (err) {
        setError('Failed to load blog post');
        console.error(err);
      } finally {
        setLoading(false);
      }
    }

    fetchBlog();
  }, [slug, router]);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError('');

    try {
      const token = getToken();
      if (!token) throw new Error('Not authenticated');
      
      if (!blog || !blog.id) {
        throw new Error('Blog post not found or missing ID');
      }

      // Use the blog ID for the update operation, not the slug
      await updateBlog(token, blog.id, {
        title,
        content,
      });

      // Show success message
      toast.success('Blog post updated successfully');
      
      // Redirect to manage blogs page
      router.push('/admin/blogs');
    } catch (err) {
      toast.error('Failed to update blog post');
      console.error(err);
      setError('Failed to update blog post');
    } finally {
      setSaving(false);
    }
  }

  return (
    <AdminLayout>
      <div className="px-4 py-6 sm:px-0">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-semibold text-gray-900">Edit Blog Post</h1>
          <button
            onClick={() => router.push('/admin/blogs')}
            className="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <svg className="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            Back to blogs
          </button>
        </div>

        {loading ? (
          <div className="bg-white shadow overflow-hidden sm:rounded-lg p-8">
            <div className="flex flex-col items-center justify-center h-64">
              <LoadingSpinner size="large" />
              <p className="mt-4 text-gray-500">Loading blog post...</p>
            </div>
          </div>
        ) : error && !blog ? (
          <div className="bg-white shadow overflow-hidden sm:rounded-lg p-8">
            <div className="text-center py-10">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-red-100 mb-4">
                <svg className="h-8 w-8 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 className="text-lg font-medium text-gray-900">Error</h3>
              <p className="mt-2 text-gray-500">{error}</p>
              <button
                onClick={() => router.push('/admin/blogs')}
                className="mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                Back to blogs
              </button>
            </div>
          </div>
        ) : (
          <div className="bg-white shadow overflow-hidden sm:rounded-lg">
            <form onSubmit={handleSubmit} className="px-4 py-5 sm:p-6 space-y-6">
              {blog && (
                <div className="bg-blue-50 rounded-md p-4">
                  <div className="flex">
                    <div className="flex-shrink-0">
                      <svg className="h-5 w-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                    </div>
                    <div className="ml-3 flex flex-col sm:flex-row sm:justify-between w-full">
                      <div>
                        <p className="text-sm text-blue-700">
                          <span className="font-medium">Slug:</span> {blog.slug}
                        </p>
                        <p className="text-sm text-blue-700 mt-1">
                          <span className="font-medium">Published:</span> {new Date(blog.published_at).toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
                        </p>
                      </div>
                      <div className="mt-2 sm:mt-0">
                        <a 
                          href={`/blog/${blog.slug}`}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="inline-flex items-center px-3 py-1 border border-transparent text-sm leading-4 font-medium rounded-md text-blue-700 bg-blue-100 hover:bg-blue-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                        >
                          <svg className="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                          </svg>
                          Preview
                        </a>
                      </div>
                    </div>
                  </div>
                </div>
              )}
              
              <div className="space-y-1">
                <label htmlFor="title" className="block text-sm font-medium text-gray-700">
                  Title
                </label>
                <input
                  type="text"
                  id="title"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  required
                  className="block text-black w-full shadow-sm sm:text-sm border-gray-300 rounded-md p-2 border focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter blog title"
                />
              </div>

              <div className="space-y-1">
                <label htmlFor="content" className="block text-sm font-medium text-gray-700">
                  Content
                </label>
                <textarea
                  id="content"
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  required
                  rows={15}
                  className="block w-full text-black shadow-sm sm:text-sm border-gray-300 rounded-md p-2 border focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Write your blog content here..."
                />
              </div>

              <div className="flex justify-end">
                <button
                  type="button"
                  onClick={() => router.push('/admin/blogs')}
                  className="mr-3 inline-flex justify-center py-2 px-4 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={saving}
                  className={`inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 ${
                    saving ? "opacity-50 cursor-not-allowed" : ""
                  }`}
                >
                  {saving ? (
                    <>
                      <LoadingSpinner size="small" color="white" className="mr-2" />
                      Saving...
                    </>
                  ) : (
                    "Save Changes"
                  )}
                </button>
              </div>
            </form>
          </div>
        )}
      </div>
    </AdminLayout>
  );
}
