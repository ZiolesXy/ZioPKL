import LegalPage from '@/components/shared/LegalPage'

function page() {
  return (
    <LegalPage
      title="Privacy Policy"
      updatedAt="26 Februari 2026"
      intro="Ini adalah placeholder kebijakan privasi (dummy). Konten ini akan diganti dengan versi final."
      sections={[
        {
          title: "1. Informasi yang Kami Kumpulkan",
          paragraphs: [
            "Placeholder: kami dapat mengumpulkan informasi akun, informasi transaksi, dan data penggunaan aplikasi.",
            "Placeholder: detail spesifik akan ditentukan kemudian.",
          ],
        },
        {
          title: "2. Cara Kami Menggunakan Informasi",
          paragraphs: [
            "Placeholder: untuk memproses pesanan, meningkatkan layanan, dan mengelola keamanan akun.",
          ],
        },
        {
          title: "3. Berbagi Data",
          paragraphs: [
            "Placeholder: kami dapat berbagi data dengan penyedia layanan pihak ketiga seperlunya untuk menjalankan layanan.",
          ],
        },
        {
          title: "4. Kontak",
          paragraphs: [
            "Placeholder: Jika ada pertanyaan, hubungi kami melalui kanal dukungan resmi.",
          ],
        },
      ]}
    />
  )
}

export default page
