import Image from "next/image";
import { getBlogs } from "@/lib/api";
import Link from "next/link";
import { Blog } from "@/types/blog";
import { formatDate } from "@/lib/utils";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface BlogsApiResponse {
  total: number;
  blogs: Blog[];
  meta: {
    page: number;
    limit: number;
    totalPage: number;
    totalItems: number;
  };
}

async function getPageData() {
  try {
    const response = (await getBlogs()) as BlogsApiResponse;
    return { blogs: response.blogs || [] };
  } catch (error) {
    console.error("Error fetching blogs:", error);
    return { blogs: [] };
  }
}

export default async function Home() {
  const { blogs }: { blogs: Blog[] } = await getPageData();

  return (
    <div className="min-h-screen bg-slate-50 dark:bg-slate-900">
      <header className="bg-white dark:bg-slate-800 border-b border-slate-200 dark:border-slate-700 shadow-sm">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8 flex justify-between items-center">
          <h1 className="text-3xl font-bold text-slate-900 dark:text-slate-50">
            BlogKu
          </h1>
          <Link href="/admin/login" className="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300 font-medium">
            Admin Login
          </Link>
        </div>
      </header>
      
      <main className="max-w-7xl mx-auto py-10 px-4 sm:px-6 lg:px-8">
        <h2 className="text-3xl font-bold tracking-tight mb-8 text-slate-900 dark:text-slate-50">
          Latest Blog Posts
        </h2>
        
        {blogs.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-slate-500 dark:text-slate-400 text-lg">
              No blog posts found. Check back later for updates!
            </p>
          </div>
        ) : (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {blogs.map((blog: Blog) => (
              <Link href={`/blog/${blog.slug}`} key={blog.id} className="block transition-all hover:scale-[1.02] duration-200">
                <Card className="h-full overflow-hidden hover:shadow-md transition-shadow">
                  {blog.image_path && (
                    <div className="aspect-video w-full relative overflow-hidden">
                      <img 
                        src={`http://localhost:8080/public/uploads/${blog.image_path}`} 
                        alt={blog.title}
                        className="object-cover w-full h-full"
                      />
                    </div>
                  )}
                  <CardHeader className="pb-2">
                    <CardTitle className="line-clamp-2 text-xl">{blog.title}</CardTitle>
                    <CardDescription>
                      <time dateTime={blog.published_at}>
                        {formatDate(blog.published_at, "short")}
                      </time>
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <p className="text-slate-600 dark:text-slate-300 line-clamp-3">
                      {blog.content.substring(0, 150)}
                      {blog.content.length > 150 ? "..." : ""}
                    </p>
                  </CardContent>
                  <CardFooter>
                    <Badge variant="outline" className="text-blue-600 dark:text-blue-400">Read more</Badge>
                  </CardFooter>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </main>
      
      <footer className="bg-slate-800 text-slate-200 py-8 mt-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="text-slate-300">© {new Date().getFullYear()} BlogKu. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
}
