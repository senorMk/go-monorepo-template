export default function Terms() {
  return (
    <div className="px-6 py-20 max-w-3xl mx-auto">
      <div className="mb-12">
        <h1 className="text-4xl font-bold tracking-tight mb-3">Terms &amp; Conditions</h1>
        <p className="text-sm text-gray-400">Last updated {new Date().toLocaleDateString('en-GB', { day: 'numeric', month: 'long', year: 'numeric' })}</p>
      </div>
      <div className="space-y-10 text-gray-600 leading-relaxed">
        <section>
          <h2 className="text-xl font-semibold text-gray-900 mb-3">1. Acceptance of Terms</h2>
          <p>By using APP_DISPLAY_NAME, you agree to be bound by these Terms &amp; Conditions.</p>
        </section>
        <section>
          <h2 className="text-xl font-semibold text-gray-900 mb-3">2. Contact Us</h2>
          <p>
            For any questions, please contact us at{' '}
            <a href="mailto:info@APP_DOMAIN" className="text-brand-600 hover:underline">info@APP_DOMAIN</a>.
          </p>
        </section>
      </div>
    </div>
  )
}
