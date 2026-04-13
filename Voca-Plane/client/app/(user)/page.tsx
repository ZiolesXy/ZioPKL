import {
  Card,
  CardContent,
} from "@/components/ui/card"
import Image from "next/image"
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group"
import { Plane, Search, ArrowRight } from "lucide-react"
import { CardTicket } from "@/components/CardTicket"
import Link from 'next/link'
import { getFlights } from "@/lib/api/FlightApi"
import { ApiResponse, Flight } from "@/lib/type/flight";

async function SearchPage() {
  const flightsData = await getFlights()

  return (
    <div className="min-h-screen bg-background">
      {/* Premium Hero Section */}
      <section className="relative pt-40 pb-32 overflow-hidden">
        {/* Simplified Background Mesh */}
        <div className="absolute top-0 left-0 w-full h-full -z-10 bg-[radial-gradient(circle_at_50%_0%,rgba(var(--primary),0.05),transparent_70%)]" />
        <div className="absolute -top-40 -right-40 w-96 h-96 bg-primary/10 rounded-full blur-[80px] -z-10" />
        <div className="absolute top-1/2 -left-20 w-80 h-80 bg-indigo-500/5 rounded-full blur-[60px] -z-10" />
        
        <div className="container mx-auto px-8">
          <div className="max-w-4xl mx-auto text-center space-y-10 animate-in fade-in slide-in-from-bottom-8 duration-700">
            <div className="inline-flex items-center gap-3 px-5 py-2.5 rounded-full bg-slate-900/5 dark:bg-white/5 border border-primary/10">
              <span className="flex size-2 rounded-full bg-primary" />
              <span className="text-[11px] font-black uppercase tracking-[0.3em] text-primary/80">Premium Jet Experience</span>
            </div>

            <h1 className="text-7xl lg:text-[10rem] font-black text-slate-900 dark:text-white tracking-tighter leading-[0.85]">
              Sky&apos;s no <br />
              <span className="text-primary italic">longer the limit.</span>
            </h1>

            <p className="max-w-2xl mx-auto text-slate-500 dark:text-slate-400 text-xl font-medium leading-relaxed">
              Experience the future of seamless flight booking. <br className="hidden md:block" />
              High-tier service for world-class travelers.
            </p>
            
            {/* Integrated Search Bar - Optimized & Functional */}
            <div className="w-full max-w-4xl mx-auto group">
              <form action="/flight" method="GET" className="bg-white/90 dark:bg-slate-900/90 backdrop-blur-xl p-3 rounded-[3rem] shadow-xl border border-white/40 dark:border-white/5 flex flex-col md:flex-row gap-3 transition-all">
                <div className="flex-1 flex items-center px-6 gap-3">
                  <div className="flex flex-col items-start min-w-[120px]">
                    <span className="text-[9px] font-black uppercase tracking-widest text-primary/60 ml-3 mb-0.5">Origin</span>
                    <input 
                      type="text" 
                      name="origin"
                      placeholder="Jakarta?" 
                      required
                      className="w-full h-10 bg-transparent border-none focus:ring-0 text-slate-900 dark:text-white font-bold text-lg placeholder:text-slate-300 dark:placeholder:text-slate-700 p-3"
                    />
                  </div>
                  
                  <div className="h-8 w-px bg-slate-200 dark:bg-slate-700 hidden md:block" />
                  
                  <div className="flex flex-col items-start min-w-[120px]">
                    <span className="text-[9px] font-black uppercase tracking-widest text-primary/60 ml-3 mb-0.5">Destination</span>
                    <input 
                      type="text" 
                      name="destination"
                      placeholder="Bali?" 
                      required
                      className="w-full h-10 bg-transparent border-none focus:ring-0 text-slate-900 dark:text-white font-bold text-lg placeholder:text-slate-300 dark:placeholder:text-slate-700 p-3"
                    />
                  </div>
                </div>
                
                <button type="submit" className="h-16 px-12 bg-primary text-white font-black text-lg rounded-[2.5rem] shadow-lg shadow-primary/20 hover:scale-[1.01] active:scale-95 transition-transform duration-200 flex items-center justify-center gap-2">
                  <Search className="size-5" />
                  Book Now
                </button>
              </form>
            </div>
          </div>
        </div>
      </section>

      {/* Stats/Quick Glance Bar */}
      <section className="py-12 border-y border-slate-100 bg-white/50">
        <div className="container mx-auto px-6">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            {[
              { label: "Daily Flights", val: "120+" },
              { label: "Destinations", val: "45" },
              { label: "Happy Clients", val: "12k+" },
              { label: "Partner Airlines", val: "15" },
            ].map((stat, i) => (
              <div key={i} className="text-center md:text-left space-y-1">
                <p className="text-[10px] font-black uppercase tracking-widest text-slate-400">{stat.label}</p>
                <p className="text-3xl font-black text-slate-900 tracking-tighter">{stat.val}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Popular Destinations Section */}
      <section className="py-24">
        <div className="container mx-auto px-6">
          <div className="flex flex-col md:flex-row justify-between items-end mb-16 gap-4">
            <div className="max-w-lg">
              <span className="text-primary font-black tracking-[0.3em] uppercase text-[10px] mb-4 block">Recommended for you</span>
              <h2 className="text-5xl font-black text-slate-900 tracking-tight">Top Destinations</h2>
            </div>
            <Link href="/flight" className="group flex items-center gap-2 text-slate-400 font-bold text-sm hover:text-primary transition-colors">
              Explores all deals <ArrowRight className="size-4 group-hover:translate-x-1 transition-transform" />
            </Link>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-y-16 gap-x-8">
            {flightsData.data.slice(0, 4).map((flight: Flight) => (
              <CardTicket key={flight.id} flight={flight} />
            ))}
          </div>
        </div>
      </section>

      {/* Features Section (Why Choose Voca Plane?) */}
      <section className="py-24 bg-slate-50/50">
        <div className="container mx-auto px-4">
          <div className="flex flex-col lg:flex-row items-center gap-20">
            {/* Left: Image with Floating Cards */}
            <div className="w-full lg:w-1/2">
              <div className="relative">
                <div className="relative h-[600px] w-full max-w-md mx-auto aspect-[3/4] rounded-[3rem] overflow-hidden shadow-xl">
                  <Image
                    src="/destination.jpg"
                    alt="Why Voca Plane"
                    fill
                    className="object-cover"
                  />
                  <div className="absolute inset-0 bg-primary/10 mix-blend-overlay" />
                </div>

                {/* Simplified Floating Elements */}
                <Card className="absolute top-10 -left-10 md:-left-20 p-6 rounded-3xl shadow-xl border-none max-w-[200px]">
                  <div className="flex items-center gap-4">
                    <div className="bg-yellow-100 p-3 rounded-2xl">
                      <span className="text-yellow-600 text-2xl font-bold">⭐</span>
                    </div>
                    <div>
                      <p className="text-lg font-black text-slate-900 leading-none">4.9/5</p>
                      <p className="text-[10px] text-slate-500 font-bold uppercase tracking-widest mt-1">Global Rating</p>
                    </div>
                  </div>
                </Card>

                <Card className="absolute bottom-20 -right-10 p-6 rounded-[2rem] shadow-xl border-none bg-white/95 backdrop-blur-md max-w-[240px]">
                  <div className="space-y-4">
                    <div className="flex -space-x-3">
                      {[1, 2, 3, 4].map((i) => (
                        <div key={i} className="w-10 h-10 rounded-full border-2 border-white bg-slate-200" />
                      ))}
                      <div className="w-10 h-10 rounded-full border-2 border-white bg-primary flex items-center justify-center text-[10px] text-white font-black">
                        +1K
                      </div>
                    </div>
                    <p className="text-xs font-bold text-slate-700 leading-relaxed">
                      "Best platform to book my business flights with zero hassle."
                    </p>
                  </div>
                </Card>
              </div>
            </div>

            <div className="w-full lg:w-1/2 space-y-12">
              <div>
                <h2 className="text-4xl md:text-5xl font-black text-slate-900 leading-tight mb-6">
                  Mengapa Memilih<br />
                  <span className="text-primary">Voca Plane?</span>
                </h2>
                <p className="text-slate-500 text-xl font-medium leading-relaxed">
                  Kami mengutamakan kenyamanan, transparansi, dan keamanan dalam setiap perjalanan udara Anda.
                </p>
              </div>

              <div className="space-y-10">
                {[
                  {
                    num: "01",
                    title: "Temukan Destinasi Terbaik",
                    desc: "Pilih rute dari berbagai bandara internasional utama dengan kode IATA resmi untuk akurasi maksimal."
                  },
                  {
                    num: "02",
                    title: "Transparansi Harga & Kelas",
                    desc: "Pilihan fleksibel mulai dari Economy, Business, hingga First Class dengan detail fasilitas yang jelas."
                  },
                  {
                    num: "03",
                    title: "Keamanan Transaksi Terjamin",
                    desc: "Sistem booking modern dengan enkripsi tingkat tinggi untuk memastikan setiap transaksi Anda aman."
                  }
                ].map((feature, i) => (
                  <div key={i} className="flex items-start gap-6 group">
                    <div className="flex-shrink-0 w-16 h-16 flex items-center justify-center rounded-3xl bg-white shadow-xl text-primary font-black text-2xl transition-all group-hover:bg-primary group-hover:text-white group-hover:-translate-y-1">
                      {feature.num}
                    </div>
                    <div>
                      <h3 className="text-2xl font-bold text-slate-800 mb-2">{feature.title}</h3>
                      <p className="text-slate-500 text-lg leading-relaxed">{feature.desc}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Explore More Section */}
      <section className="relative py-32 px-4 overflow-hidden">
        <div className="absolute inset-0 bg-slate-950 -z-10" />
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[1000px] h-[1000px] bg-primary/5 rounded-full blur-[100px] -z-10" />

        <div className="container mx-auto relative z-10">
          <div className="flex flex-col md:flex-row justify-between items-center md:items-end mb-20 gap-8">
            <div className="text-center md:text-left">
              <span className="text-primary font-black tracking-[0.3em] uppercase text-xs mb-4 block">Discover More</span>
              <h2 className="text-5xl md:text-7xl font-black text-primary leading-none">Explore All<br />Available Flights</h2>
            </div>
            <Link
              href="/flight"
              className="flex items-center gap-3 px-10 py-5 bg-primary text-white hover:bg-primary/90 rounded-2xl font-bold transition-all shadow-xl shadow-primary/20 group scale-110 md:scale-100"
            >
              See All Flights
              <ArrowRight className="w-5 h-5 transition-transform group-hover:translate-x-1" />
            </Link>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-y-16 gap-x-8">
            {flightsData.data.slice(4, 8).map((flight: Flight) => (
              <CardTicket key={flight.id} flight={flight} />
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-32 px-4">
        <div className="container mx-auto max-w-6xl">
          <div className="relative bg-primary rounded-[4rem] p-12 md:p-24 text-center overflow-hidden shadow-xl shadow-primary/20">
            {/* Decorative element */}
            <div className="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full blur-xl -translate-y-1/2 translate-x-1/2" />
            
            <div className="relative z-10 space-y-10">
              <h2 className="text-4xl md:text-7xl font-black text-white tracking-tighter leading-tight">
                Pesan dan Pergi<br />Sekarang Juga
              </h2>
              <p className="text-white/80 text-xl font-medium max-w-2xl mx-auto">
                Jelajahi destinasi impian Anda dengan standar layanan premium dan harga terbaik hanya di Voca Plane.
              </p>
              <div className="flex flex-col sm:flex-row justify-center gap-6">
                <Link href="/flight">
                  <button className="w-full sm:w-auto px-12 py-5 bg-white text-primary font-black rounded-2xl transition-all shadow-xl hover:scale-105 active:scale-95">
                    Pesan Tiket Sekarang
                  </button>
                </Link>
                <Link href="/register">
                  <button className="w-full sm:w-auto px-12 py-5 bg-transparent border-2 border-white/30 text-white font-black rounded-2xl hover:bg-white/10 transition-all">
                    Daftar Member
                  </button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}

export default SearchPage
