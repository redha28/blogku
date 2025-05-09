"use client";

import { getBlogBySlug } from "@/lib/api";
import Link from "next/link";
import Image from "next/image";
import { notFound, useParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";
import { formatDate } from "@/lib/utils";
import { Button } from "@/components/ui/Button";
import { Card, CardContent } from "@/components/ui/card";
import { ArrowLeft } from "lucide-react";

interface Blog {
  id: number | string;
  title: string;
  content: string;
  slug: string;
  published_at: string | Date;
  image_path?: string;
}

function BlogSkeleton() {
  return (
    <div className="space-y-4 animate-pulse">
      <div className="h-8 w-3/4 mb-4 bg-gray-200 dark:bg-gray-700 rounded"></div>
      <div className="h-4 w-1/4 mb-6 bg-gray-200 dark:bg-gray-700 rounded"></div>
      <div className="h-64 w-full mb-6 bg-gray-200 dark:bg-gray-700 rounded-lg"></div>
      <div className="space-y-3">
        <div className="h-4 w-full bg-gray-200 dark:bg-gray-700 rounded"></div>
        <div className="h-4 w-full bg-gray-200 dark:bg-gray-700 rounded"></div>
        <div className="h-4 w-full bg-gray-200 dark:bg-gray-700 rounded"></div>
        <div className="h-4 w-5/6 bg-gray-200 dark:bg-gray-700 rounded"></div>
      </div>
    </div>
  );
}

function BlogContent({ blog }: { blog: Blog }) {
  return (
    <article className="prose prose-lg max-w-none dark:prose-invert">
      <h1 className="text-3xl font-bold mb-4 text-blue-800 dark:text-blue-400">
        {blog.title}
      </h1>

      <p className="text-blue-600 dark:text-blue-400 mb-6 font-medium">
        Published on{" "}
        {formatDate(blog.published_at, "long")}
      </p>

      {blog.image_path && (
        <div className="relative h-80 w-full mb-8 overflow-hidden rounded-lg shadow-lg">
          <Image
            src={`http://localhost:8080/public/uploads/${blog.image_path}`}
            alt={blog.title}
            fill
            className="object-cover"
            priority
          />
        </div>
      )}

      <Card className="border-l-4 border-blue-500">
        <CardContent className="p-6 prose dark:prose-invert">
          {blog.content.split("\n").map((paragraph: string, index: number) =>
            paragraph.trim() ? (
              <p key={index} className="mb-4 text-gray-700 dark:text-gray-300 leading-relaxed">
                {paragraph}
              </p>
            ) : (
              <div key={index} className="h-4" />
            )
          )}
        </CardContent>
      </Card>
    </article>
  );
}

export default function BlogPage() {
  const params = useParams();
  const slug = params.slug as string;
  const [blog, setBlog] = useState<Blog | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function loadBlog() {
      try {
        setLoading(true);
        const blogData = await getBlogBySlug(slug);
        setBlog(blogData);
      } catch (error) {
        console.error(`Error fetching blog with slug ${slug}:`, error);
        setError("Failed to load blog post");
      } finally {
        setLoading(false);
      }
    }
    
    loadBlog();
  }, [slug]);

  if (error) {
    notFound();
  }

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-900">
      <header className="bg-gradient-to-r from-blue-600 to-blue-800 text-white shadow-md">
        <div className="container mx-auto py-6 px-4 sm:px-6 lg:px-8 flex justify-between items-center">
          <Link
            href="/"
            className="text-3xl font-bold hover:text-blue-100 transition-colors"
          >
            BlogKu
          </Link>

          <Button variant="secondary" asChild>
            <Link href="/">
              Back to Home
            </Link>
          </Button>
        </div>
      </header>

      <main className="container mx-auto py-10 px-4 sm:px-6 lg:px-8 max-w-4xl">
        {loading ? (
          <BlogSkeleton />
        ) : blog ? (
          <BlogContent blog={blog} />
        ) : null}

        <div className="mt-10 flex justify-center">
          <Button variant="outline" asChild>
            <Link href="/" className="flex items-center text-blue-600 dark:text-blue-400">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Back to all articles
            </Link>
          </Button>
        </div>
      </main>

      <footer className="bg-blue-800 text-white py-8 mt-12">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="text-blue-200">© {new Date().getFullYear()} BlogKu. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}
