import ProductCard from '@/components/product/ProductCard'
import { getProducts } from '@/lib/api/product'

function shuffleArray<T>(array: T[]): T[] {
    const shuffled = [...array];
    for (let i = shuffled.length - 1; i > 0; i--) {
        const j = Math.floor(Math.random() * (i + 1));
        [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]]
    }
    return shuffled
}
async function ProductPage() {
    const rawProducts = await getProducts()
    const products = shuffleArray(rawProducts || [])
    return (
        <div>
            <div className="p-6">
                <h2 className="text-2xl font-bold mb-6">Semua Produk</h2>
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                    {products.map((item) => (
                        <ProductCard
                            id={item.id}
                            slug={item.slug}
                            key={item.id}
                            name={item.name}
                            price={item.price}
                            category={item.category}
                            stock={item.stock}
                            image={item.image_url}
                        />
                    ))}
                </div>
            </div>
        </div>
    )
}

export default ProductPage

