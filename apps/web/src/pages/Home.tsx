export default function Home() {
  return (
    <section className="px-6 py-28 max-w-4xl mx-auto text-center">
      <h1 className="text-5xl font-bold tracking-tight leading-tight mb-6">
        Welcome to <span className="text-brand-600">APP_DISPLAY_NAME</span>
      </h1>
      <p className="text-xl text-gray-500 max-w-2xl mx-auto mb-10">
        Your app description goes here.
      </p>
      <div id="download" className="flex flex-col sm:flex-row gap-4 justify-center">
        <a
          href="#"
          className="inline-flex items-center justify-center gap-3 bg-gray-900 text-white px-6 py-3.5 rounded-xl font-medium hover:bg-gray-700 transition-colors"
        >
          <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
            <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.8-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
          </svg>
          <div className="text-left">
            <div className="text-xs opacity-75">Download on the</div>
            <div className="text-sm font-semibold">App Store</div>
          </div>
        </a>
        <a
          href="#"
          className="inline-flex items-center justify-center gap-3 bg-gray-900 text-white px-6 py-3.5 rounded-xl font-medium hover:bg-gray-700 transition-colors"
        >
          <svg className="w-6 h-6" viewBox="0 0 24 24" fill="currentColor">
            <path d="M3.18 23.76c.3.17.64.24.99.2L15.39 12 3.18 0c-.35-.04-.69.03-.99.2-.59.34-.94.96-.94 1.65v20.26c0 .69.35 1.31.94 1.65zM16.54 10.85l2.83-2.83-9.9-5.71 7.07 8.54zM19.37 15.98l-2.83-2.83-7.07 8.54 9.9-5.71zM21.49 9.89l-2.85-1.64-2.93 2.93 2.93 2.93 2.85-1.64c.81-.47 1.3-1.32 1.3-2.29s-.49-1.82-1.3-2.29z"/>
          </svg>
          <div className="text-left">
            <div className="text-xs opacity-75">Get it on</div>
            <div className="text-sm font-semibold">Google Play</div>
          </div>
        </a>
      </div>
    </section>
  )
}
