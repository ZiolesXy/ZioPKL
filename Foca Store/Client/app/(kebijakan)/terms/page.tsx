import LegalPage from '@/components/shared/LegalPage'

function page() {
  return (
    <LegalPage
      title="Terms of Service"
      updatedAt="26 Februari 2026"
      intro="Ini adalah placeholder Syarat & Ketentuan (dummy). Konten ini akan diganti dengan versi final."
      sections={[
        {
          title: "1. Penggunaan Layanan",
          paragraphs: [
            "Placeholder: Dengan menggunakan layanan ini, kamu menyetujui syarat yang berlaku.",
            "Placeholder: Detail aturan penggunaan akan diisi nanti.",
          ],
        },
        {
          title: "2. Akun & Keamanan",
          paragraphs: [
            "Placeholder: Kamu bertanggung jawab menjaga kerahasiaan akun dan aktivitas yang terjadi di akunmu.",
          ],
        },
        {
          title: "3. Pembayaran",
          paragraphs: [
            "Placeholder: Ketentuan harga, pembayaran, dan pajak akan mengikuti kebijakan yang berlaku.",
          ],
        },
        {
          title: "4. Perubahan Ketentuan",
          paragraphs: [
            "Placeholder: Kami dapat memperbarui syarat dan ketentuan dari waktu ke waktu.",
          ],
        },
      ]}
    />
  )
}

export default page
