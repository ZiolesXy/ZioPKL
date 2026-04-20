import Image from "next/image"

export default function AuthHeroPanel() {
  return (
    <div className="hidden lg:flex lg:items-center lg:justify-center bg-muted relative min-h-150 p-12 rounded-r-3xl overflow-hidden">
      <Image
        src="/loginimage.avif"
        alt="Login Background"
        fill
        className="object-cover"
        priority
      />

      <div className="relative z-10 w-full max-w-xs rounded-2xl border border-white/20 bg-white/10 p-6 shadow-2xl backdrop-blur-md min-h-100 flex flex-col">
        <div className="space-y-2">
          <h3 className="text-xl font-semibold text-white text-center">
            Voca Store Here
          </h3>
          <p className="text-lg text-slate-200 leading-relaxed text-center">
            E-commerce Nomor 1 di Jagat Sawit yang sudah dipercaya oleh Prabowo
            Subianto.
          </p>
        </div>

        <div className="mt-auto pt-4">
          <p className="text-lg font-medium text-white leading-relaxed">
            Moto Kami :{" "}
            <span className="italic text-white">Nyawit Nomor 1</span>
          </p>
        </div>
      </div>
    </div>
  )
}
