import { notFound } from "next/navigation";
import { getFlightById } from "@/lib/api/FlightApi";
import { BookingForm } from "@/components/BookingForm";
import { Plane, Clock, MapPin, Users } from "lucide-react";

export default async function DetailFlightPage({
  params,
}: {
  params: { id: string };
}) {
  const { id } = await params;

  let flight;
  try {
    flight = await getFlightById(id);
  } catch (error) {
    notFound();
  }

  if (!flight) {
    notFound();
  }

  const departureDate = new Date(flight.departure_time);
  const arrivalDate = new Date(flight.arrival_time);

  const formatTime = (date: Date) =>
    date.toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit" });

  const formatDate = (date: Date) =>
    date.toLocaleDateString("id-ID", {
      weekday: "long",
      day: "numeric",
      month: "long",
      year: "numeric",
    });

  // Hitung durasi penerbangan
  const durationMs = arrivalDate.getTime() - departureDate.getTime();
  const durationHours = Math.floor(durationMs / (1000 * 60 * 60));
  const durationMinutes = Math.floor((durationMs % (1000 * 60 * 60)) / (1000 * 60));

  return (
    <div className="container mx-auto py-10 px-4">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-2">
          <div className="h-10 w-10 rounded-full bg-indigo-100 flex items-center justify-center">
            <Plane className="size-5 text-indigo-600 rotate-45" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-slate-800">
              Detail Penerbangan
            </h1>
            <p className="text-sm text-slate-500">
              {flight.flight_number} · {flight.airline.name}
            </p>
          </div>
        </div>
      </div>

      {/* Flight Info Card */}
      <div className="bg-white rounded-2xl shadow-md border border-slate-100 p-6 mb-8">
        <div className="flex flex-col md:flex-row items-center gap-6 md:gap-10">
          {/* Origin */}
          <div className="text-center md:text-left flex-1">
            <p className="text-3xl font-bold text-slate-800">
              {flight.origin.code}
            </p>
            <p className="text-sm text-slate-600 font-medium">
              {flight.origin.city}
            </p>
            <p className="text-xs text-slate-400">{flight.origin.name}</p>
            <p className="text-lg font-semibold text-indigo-600 mt-2">
              {formatTime(departureDate)}
            </p>
          </div>

          {/* Flight Path */}
          <div className="flex flex-col items-center gap-1 flex-shrink-0">
            <p className="text-xs text-slate-400 font-medium">
              {durationHours}j {durationMinutes}m
            </p>
            <div className="flex items-center gap-2">
              <div className="h-2 w-2 rounded-full bg-indigo-600" />
              <div className="h-[2px] w-24 bg-linear-to-r from-indigo-600 to-violet-500" />
              <Plane className="size-4 text-indigo-600 rotate-90" />
              <div className="h-[2px] w-24 bg-linear-to-r from-violet-500 to-indigo-600" />
              <div className="h-2 w-2 rounded-full bg-indigo-600" />
            </div>
            <p className="text-[10px] text-slate-400 uppercase tracking-wider font-medium">
              Direct Flight
            </p>
          </div>

          {/* Destination */}
          <div className="text-center md:text-right flex-1">
            <p className="text-3xl font-bold text-slate-800">
              {flight.destination.code}
            </p>
            <p className="text-sm text-slate-600 font-medium">
              {flight.destination.city}
            </p>
            <p className="text-xs text-slate-400">{flight.destination.name}</p>
            <p className="text-lg font-semibold text-indigo-600 mt-2">
              {formatTime(arrivalDate)}
            </p>
          </div>
        </div>

        {/* Detail Badges */}
        <div className="mt-6 pt-6 border-t border-slate-100 flex flex-wrap gap-4">
          <div className="flex items-center gap-2 bg-slate-50 px-4 py-2 rounded-xl">
            <Clock className="size-4 text-slate-400" />
            <span className="text-sm text-slate-600">
              {formatDate(departureDate)}
            </span>
          </div>
          <div className="flex items-center gap-2 bg-slate-50 px-4 py-2 rounded-xl">
            <Users className="size-4 text-slate-400" />
            <span className="text-sm text-slate-600">
              {flight.available_seats} kursi tersedia
            </span>
          </div>
          <div className="flex items-center gap-2 bg-slate-50 px-4 py-2 rounded-xl">
            <MapPin className="size-4 text-slate-400" />
            <span className="text-sm text-slate-600">
              {flight.airline.name}
            </span>
          </div>
        </div>

        {/* Class Pricing */}
        <div className="mt-6 pt-6 border-t border-slate-100">
          <h3 className="text-sm font-semibold text-slate-500 uppercase tracking-wider mb-3">
            Kelas & Harga
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            {flight.classes.map((flightClass) => {
              const classColors: Record<string, string> = {
                first: "border-amber-400 bg-amber-50 text-amber-700",
                business: "border-blue-400 bg-blue-50 text-blue-700",
                economy: "border-indigo-400 bg-indigo-50 text-indigo-700",
                economi: "border-indigo-400 bg-indigo-50 text-indigo-700",
              };
              const colorClass =
                classColors[flightClass.class_type.toLowerCase()] ||
                "border-slate-300 bg-slate-50 text-slate-700";

              return (
                <div
                  key={flightClass.id}
                  className={`border-2 rounded-xl p-4 ${colorClass}`}
                >
                  <p className="text-xs font-bold uppercase tracking-wider">
                    {flightClass.class_type}
                  </p>
                  <p className="text-xl font-bold mt-1">
                    Rp {flightClass.price.toLocaleString("id-ID")}
                  </p>
                  <p className="text-xs mt-1 opacity-70">
                    {flightClass.total_seats} kursi total
                  </p>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      {/* Booking Form + Seat Grid */}
      <div>
        <h2 className="text-xl font-bold text-slate-800 mb-2">
          Pesan Tiket Penerbangan
        </h2>
        <p className="text-sm text-slate-500 mb-6">
          Isi data penumpang dan pilih kursi untuk melanjutkan pemesanan.
        </p>
        <BookingForm flight={flight} />
      </div>
    </div>
  );
}