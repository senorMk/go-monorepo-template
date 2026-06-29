import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import express from 'express'
import compression from 'compression'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const isProduction = process.env.NODE_ENV === 'production'
const port = process.env.PORT || 3000
const base = process.env.BASE || '/'

const app = express()

let vite: Awaited<ReturnType<typeof import('vite')['createServer']>> | undefined

if (!isProduction) {
  const { createServer } = await import('vite')
  vite = await createServer({
    server: { middlewareMode: true },
    appType: 'custom',
    base,
  })
  app.use(vite.middlewares)
} else {
  app.use(compression())
  app.use(
    base,
    express.static(path.resolve(__dirname, 'dist/client'), { index: false }),
  )
}

app.use('*', async (req, res) => {
  try {
    const url = req.originalUrl.replace(base, '')

    let template: string
    let render: (url: string) => string | Promise<string>

    if (!isProduction && vite) {
      template = fs.readFileSync(path.resolve(__dirname, 'index.html'), 'utf-8')
      template = await vite.transformIndexHtml(url, template)
      render = (await vite.ssrLoadModule('/src/entry-server.tsx')).render
    } else {
      template = fs.readFileSync(
        path.resolve(__dirname, 'dist/client/index.html'),
        'utf-8',
      )
      render = (await import(new URL('./dist/server/entry-server.js', import.meta.url).href)).render
    }

    const appHtml = await render(url)
    const html = template.replace('<!--app-html-->', appHtml)

    res.status(200).set({ 'Content-Type': 'text/html' }).send(html)
  } catch (e: unknown) {
    if (!isProduction && vite) {
      vite.ssrFixStacktrace(e as Error)
    }
    console.error(e)
    res.status(500).end((e as Error).message)
  }
})

app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`)
})
