import { RegisterForm } from '@/components/auth/register-form'

export default function RegisterPage() {
  return (
    <div className="flex min-h-screen w-full items-center justify-center bg-slate-50 p-4 md:p-10">
      <div className="w-full max-w-4xl overflow-hidden rounded-2xl bg-white shadow-[0_20px_50px_-12px_rgba(0,0,0,0.1)] animate-in fade-in zoom-in-95 duration-500">
        <RegisterForm />
      </div>
    </div>
  )
}
