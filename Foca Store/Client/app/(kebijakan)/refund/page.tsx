import LegalPage from '@/components/shared/LegalPage'

function page() {
  return (
    <LegalPage
      title="Refund Policy"
      updatedAt="26 Februari 2026"
      intro="Ini adalah placeholder kebijakan refund/pengembalian dana (dummy). Konten ini akan diganti dengan versi final."
      sections={[
        {
          title: "1. Syarat Pengajuan Refund",
          paragraphs: [
            "Placeholder: Refund dapat diajukan dalam jangka waktu tertentu setelah pesanan diterima.",
            "Placeholder: Syarat detail (produk rusak/salah kirim/dll) akan ditentukan nanti.",
          ],
        },
        {
          title: "2. Proses Refund",
          paragraphs: [
            "Placeholder: Ajukan melalui halaman bantuan/CS dan siapkan bukti pendukung.",
            "Placeholder: Waktu proses bersifat estimasi.",
          ],
        },
        {
          title: "3. Metode Pengembalian Dana",
          paragraphs: [
            "Placeholder: Pengembalian dana dapat dilakukan ke metode pembayaran awal atau metode lain sesuai kebijakan.",
          ],
        },
        {
          title: "4. Pengecualian",
          paragraphs: [
            "Placeholder: Beberapa jenis produk mungkin tidak memenuhi syarat refund.",
          ],
        },
      ]}
    />
  )
}

export default page
