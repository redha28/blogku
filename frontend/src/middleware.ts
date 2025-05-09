import { NextRequest, NextResponse } from 'next/server';

export function middleware(request: NextRequest) {
  console.log(`[Middleware] Checking path: ${request.nextUrl.pathname}`);
  
  if (request.nextUrl.pathname.startsWith('/admin') && 
      !request.nextUrl.pathname.startsWith('/admin/login')) {
    
    // Check for authentication from cookies or Authorization header
    const authCookie = request.cookies.get('authToken');
    const authHeader = request.headers.get('authorization');
    
    console.log(`[Middleware] Auth check: Cookie=${authCookie ? 'yes' : 'no'}, Header=${authHeader ? 'yes' : 'no'}`);
    
    // If no authentication is found, redirect to login
    if (!authCookie && !authHeader) {
      console.log(`[Middleware] No authentication found, redirecting to login`);
      return NextResponse.redirect(new URL('/admin/login', request.url));
    }
    
    // Basic validation of token format
    if (authCookie && authCookie.value.length < 10) {
      console.log(`[Middleware] Invalid token format, redirecting to login`);
      return NextResponse.redirect(new URL('/admin/login', request.url));
    }
  }
  
  return NextResponse.next();
}

export const config = {
  matcher: ['/admin/:path*']
};
