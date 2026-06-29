import { Link } from 'react-router-dom'

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-white text-gray-900 flex flex-col">
      <nav className="border-b border-gray-100 px-6 py-4 flex items-center justify-between">
        <Link to="/" className="font-bold text-gray-900 text-lg hover:opacity-80 transition-opacity">
          APP_DISPLAY_NAME
        </Link>
        <div className="flex gap-6 text-sm text-gray-500">
          <a href="/#features" className="hover:text-gray-900 transition-colors">Features</a>
          <a href="/#download" className="hover:text-gray-900 transition-colors">Download</a>
          <a href="/#contact" className="hover:text-gray-900 transition-colors">Contact</a>
        </div>
      </nav>

      <main className="flex-1">{children}</main>

      <footer id="contact" className="border-t border-gray-100 px-6 py-12">
        <div className="max-w-5xl mx-auto flex flex-col md:flex-row justify-between gap-8">
          <div>
            <p className="font-semibold text-gray-900 mb-1">APP_DISPLAY_NAME</p>
            <p className="text-sm text-gray-400 mt-2">
              <a href="mailto:info@APP_DOMAIN" className="hover:text-gray-600 transition-colors">
                info@APP_DOMAIN
              </a>
            </p>
          </div>
          <div className="text-sm text-gray-400 md:text-right">
            <div className="flex gap-4 md:justify-end text-xs">
              <Link to="/terms-and-conditions" className="hover:text-gray-600 transition-colors">Terms &amp; Conditions</Link>
              <Link to="/privacy-policy" className="hover:text-gray-600 transition-colors">Privacy Policy</Link>
            </div>
            <p className="mt-3">&copy; {new Date().getFullYear()} APP_DISPLAY_NAME</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
