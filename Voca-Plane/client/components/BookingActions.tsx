"use client";

import { useState } from "react";
import { Download, Loader2, Mail } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Transaction } from "@/lib/type/transaction";

type BookingActionsProps = {
  booking: Transaction;
  isSuccess: boolean;
};

function formatDateTime(value?: string) {
  if (!value) return "-";

  return new Date(value).toLocaleString("id-ID", {
    day: "numeric",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function formatCurrency(value?: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value || 0);
}

function escapeHtml(value?: string) {
  if (!value) return "";

  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

function buildETicketHtml(booking: Transaction) {
  const passengers =
    booking.transaction_items ||
    booking.transactions_passangers ||
    booking.transactions_passanger ||
    [];
  const airlineName = booking.flight?.airline?.name || "Maskapai";
  const originCode = booking.flight?.origin?.code || "---";
  const originCity = booking.flight?.origin?.city || "-";
  const destinationCode = booking.flight?.destination?.code || "---";
  const destinationCity = booking.flight?.destination?.city || "-";
  const departureTime = formatDateTime(booking.flight?.departure_time);
  const arrivalTime = formatDateTime(booking.flight?.arrival_time);
  const totalBeforeDiscount = (booking.total_price || 0) + (booking.discount || 0);

  const passengerRows = passengers
    .map(
      (passenger, index) => `
        <tr>
          <td>${index + 1}</td>
          <td>${escapeHtml(passenger.passenger_name)}</td>
          <td>${escapeHtml(passenger.nationality)}</td>
          <td>${escapeHtml(passenger.passport_no)}</td>
          <td>${escapeHtml(passenger.class_name)}</td>
          <td>${escapeHtml(passenger.seat_number)}</td>
        </tr>
      `
    )
    .join("");

  return `<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>E-Ticket ${escapeHtml(booking.code)}</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f8fafc;
      --card: #ffffff;
      --line: #dbe4f0;
      --text: #0f172a;
      --muted: #475569;
      --accent: #2563eb;
      --soft: #eff6ff;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      padding: 32px;
      background: var(--bg);
      color: var(--text);
      font-family: Arial, Helvetica, sans-serif;
    }
    .ticket {
      max-width: 920px;
      margin: 0 auto;
      background: var(--card);
      border: 1px solid var(--line);
      border-radius: 24px;
      overflow: hidden;
    }
    .hero {
      padding: 28px 32px;
      background: linear-gradient(135deg, #1d4ed8, #0f172a);
      color: #fff;
    }
    .hero h1 {
      margin: 0 0 6px;
      font-size: 30px;
    }
    .hero p {
      margin: 0;
      opacity: 0.85;
    }
    .section {
      padding: 24px 32px;
      border-top: 1px solid var(--line);
    }
    .grid {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 16px;
    }
    .info {
      padding: 16px;
      background: #f8fafc;
      border: 1px solid var(--line);
      border-radius: 16px;
    }
    .label {
      display: block;
      margin-bottom: 8px;
      color: var(--muted);
      font-size: 12px;
      font-weight: 700;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }
    .route {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 16px;
      margin-top: 8px;
      padding: 20px;
      background: var(--soft);
      border-radius: 20px;
    }
    .route strong {
      display: block;
      font-size: 34px;
    }
    .route span {
      color: var(--muted);
    }
    table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 8px;
    }
    th, td {
      padding: 12px;
      border-bottom: 1px solid var(--line);
      text-align: left;
      font-size: 14px;
    }
    th {
      color: var(--muted);
      font-size: 12px;
      text-transform: uppercase;
      letter-spacing: 0.06em;
    }
    .summary {
      width: 100%;
      max-width: 360px;
      margin-left: auto;
    }
    .summary-row {
      display: flex;
      justify-content: space-between;
      gap: 16px;
      padding: 10px 0;
      border-bottom: 1px solid var(--line);
    }
    .summary-row.total {
      font-weight: 700;
      font-size: 18px;
      border-bottom: none;
    }
    .footer {
      color: var(--muted);
      font-size: 12px;
      line-height: 1.6;
    }
    @media print {
      body {
        padding: 0;
        background: #fff;
      }
      .ticket {
        border: none;
        border-radius: 0;
      }
    }
    @media (max-width: 640px) {
      body { padding: 12px; }
      .hero, .section { padding: 20px; }
      .grid { grid-template-columns: 1fr; }
      .route { flex-direction: column; align-items: flex-start; }
    }
  </style>
</head>
<body>
  <main class="ticket">
    <section class="hero">
      <h1>E-Ticket</h1>
      <p>${escapeHtml(airlineName)} • Booking ${escapeHtml(booking.code)}</p>
    </section>

    <section class="section">
      <span class="label">Rute Penerbangan</span>
      <div class="route">
        <div>
          <strong>${escapeHtml(originCode)}</strong>
          <span>${escapeHtml(originCity)}</span>
        </div>
        <div>
          <span>${escapeHtml(booking.flight?.flight_number || booking.code)}</span>
        </div>
        <div style="text-align:right;">
          <strong>${escapeHtml(destinationCode)}</strong>
          <span>${escapeHtml(destinationCity)}</span>
        </div>
      </div>
    </section>

    <section class="section">
      <div class="grid">
        <div class="info">
          <span class="label">Maskapai</span>
          <div>${escapeHtml(airlineName)}</div>
        </div>
        <div class="info">
          <span class="label">Status Pembayaran</span>
          <div>${escapeHtml(booking.payment_status)}</div>
        </div>
        <div class="info">
          <span class="label">Berangkat</span>
          <div>${escapeHtml(departureTime)}</div>
        </div>
        <div class="info">
          <span class="label">Tiba</span>
          <div>${escapeHtml(arrivalTime)}</div>
        </div>
      </div>
    </section>

    <section class="section">
      <span class="label">Daftar Penumpang</span>
      <table>
        <thead>
          <tr>
            <th>No</th>
            <th>Nama</th>
            <th>Kewarganegaraan</th>
            <th>Passport</th>
            <th>Kelas</th>
            <th>Kursi</th>
          </tr>
        </thead>
        <tbody>${passengerRows}</tbody>
      </table>
    </section>

    <section class="section">
      <span class="label">Rincian Harga</span>
      <div class="summary">
        <div class="summary-row">
          <span>Harga tiket</span>
          <span>${escapeHtml(formatCurrency(totalBeforeDiscount))}</span>
        </div>
        <div class="summary-row">
          <span>Diskon</span>
          <span>${escapeHtml(formatCurrency(booking.discount || 0))}</span>
        </div>
        <div class="summary-row total">
          <span>Total bayar</span>
          <span>${escapeHtml(formatCurrency(booking.total_price || 0))}</span>
        </div>
      </div>
    </section>

    <section class="section footer">
      Dokumen ini dibuat otomatis dari halaman booking aplikasi. Simpan file ini atau cetak ke PDF bila diperlukan saat check-in.
    </section>
  </main>
</body>
</html>`;
}

export function BookingActions({ booking, isSuccess }: BookingActionsProps) {
  const [isDownloading, setIsDownloading] = useState(false);

  const handleDownload = async () => {
    if (!isSuccess || isDownloading) return;

    setIsDownloading(true);

    try {
      const blob = new Blob([buildETicketHtml(booking)], {
        type: "text/html;charset=utf-8",
      });
      const url = URL.createObjectURL(blob);
      const anchor = document.createElement("a");

      anchor.href = url;
      anchor.download = `e-ticket-${booking.code}.html`;
      document.body.appendChild(anchor);
      anchor.click();
      anchor.remove();

      URL.revokeObjectURL(url);
    } finally {
      setIsDownloading(false);
    }
  };

  return (
    <>
      <Button
        variant="outline"
        className="h-14 rounded-2xl font-bold bg-white border-slate-200"
        onClick={handleDownload}
        disabled={!isSuccess || isDownloading}
      >
        {isDownloading ? (
          <Loader2 className="mr-2 size-5 animate-spin" />
        ) : (
          <Download className="mr-2 size-5" />
        )}
        {isDownloading ? "Menyiapkan E-Ticket..." : "Download E-Ticket"}
      </Button>
      <Button
        className="h-14 rounded-2xl font-bold bg-blue-600"
        type="button"
        disabled
      >
        <Mail className="mr-2 size-5" /> Kirim Email
      </Button>
    </>
  );
}
