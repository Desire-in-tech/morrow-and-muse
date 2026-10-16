import { Route, Routes } from 'react-router'
import { Shell } from './components/Shell.tsx'
import { CollectionPage } from './pages/CollectionPage.tsx'
import { HomePage } from './pages/HomePage.tsx'
import { NotFoundPage } from './pages/NotFoundPage.tsx'
import { ProductPage } from './pages/ProductPage.tsx'
import { ShopPage } from './pages/ShopPage.tsx'

export default function App() {
  return (
    <Routes>
      <Route element={<Shell />}>
        <Route index element={<HomePage />} />
        <Route path="shop" element={<ShopPage />} />
        <Route path="collections/:slug" element={<CollectionPage />} />
        <Route path="products/:slug" element={<ProductPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  )
}
