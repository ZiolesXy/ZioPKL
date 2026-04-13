import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

// Daftar rute berdasarkan grup
const ADMIN_ONLY = ['/dashboard', '/transactions'];
const USER_PROTECTED = ['/flight', '/my-bookings'];
const AUTH_PAGES = ['/login', '/register'];
const SHARED_PROTECTED = ['/update-profile'];

export function middleware(request: NextRequest) {
  const { nextUrl, cookies } = request;

  const token = cookies.get('access_token')?.value || cookies.get('token')?.value;
  const role = cookies.get('role')?.value;

  const isGuest = !token;
  const isAdmin = role === 'admin' || role === 'super_admin';
  const isUser = role === 'user'; 

  const path = nextUrl.pathname;

  // 1. FLOW GUEST: Bisa akses root (/) tapi tidak bisa booking/profile
  if (isGuest) {
    if (USER_PROTECTED.includes(path) || ADMIN_ONLY.includes(path) || SHARED_PROTECTED.includes(path)) {
      return NextResponse.redirect(new URL('/login', request.url));
    }
    return NextResponse.next();
  }

  // 2. PROTEKSI AUTH: Jika sudah login, jangan biarkan ke halaman login/regis
  if (AUTH_PAGES.includes(path)) {
    return NextResponse.redirect(new URL('/', request.url));
  }

  // 3. PROTEKSI ADMIN: User biasa dilarang masuk rute Admin
  if (isUser && ADMIN_ONLY.some(route => path.startsWith(route))) {
    return NextResponse.redirect(new URL('/', request.url));
  }

  // 4. AKSES ADMIN: Admin bisa akses semuanya (User routes & Shared)
  // Tidak perlu redirect, biarkan lanjut ke NextResponse.next()
  
  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};