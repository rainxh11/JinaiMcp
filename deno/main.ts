import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { PuppeteerService, formatSnapshot, ScrappingOptions } from './services/puppeteer.ts';
import { saveScreenshot, ensureStorageDir, fileExists } from './services/storage.ts';

const app = new Hono();

// Middleware
app.use('*', cors());

// Serve static screenshots
app.get('/instant-screenshots/:filename', async (c) => {
  const filename = c.req.param('filename');
  const filepath = `/app/local-storage/instant-screenshots/${filename}`;

  const exists = await fileExists(filepath);
  if (!exists) {
    return c.text('Screenshot not found', 404);
  }

  const file = await Deno.readFile(filepath);
  return new Response(file, {
    headers: {
      'Content-Type': 'image/png'
    }
  });
});

// Main reader endpoint
app.all('/*', async (c) => {
  const path = c.req.path();
  const respondWith = c.req.header('X-Respond-With') || 'markdown';

  // Skip if it's the screenshots endpoint
  if (path.startsWith('/instant-screenshots/')) {
    return c.notFound();
  }

  // Extract URL from path
  const url = path.slice(1);
  if (!url || url === 'favicon.ico') {
    return c.text('Reader MCP Server - Deno + Hono + Puppeteer\n\nEndpoints:\n- GET /<url> - Fetch and convert URL\n- X-Respond-With: markdown|html|text|screenshot|pageshot');
  }

  // Validate URL
  let targetUrl: URL;
  try {
    targetUrl = new URL(url.startsWith('http') ? url : `http://${url}`);
    if (!['http:', 'https:'].includes(targetUrl.protocol)) {
      throw new Error('Invalid protocol');
    }
  } catch {
    return c.text('Invalid URL', 400);
  }

  // Initialize puppeteer if needed
  const puppeteer = new PuppeteerService();
  await puppeteer.init();

  // Parse options
  const options: ScrappingOptions = {
    favorScreenshot: ['screenshot', 'pageshot'].includes(respondWith),
    timeoutMs: parseInt(c.req.header('X-Timeout') || '30000'),
    overrideUserAgent: c.req.header('X-User-Agent') || c.req.header('User-Agent'),
    waitForSelector: c.req.header('X-Wait-For-Selector')
  };

  // Scrap the URL
  const snapshot = await puppeteer.scrap(targetUrl.toString(), options);

  // Save screenshots if needed
  if (options.favorScreenshot && snapshot.screenshot) {
    const filename = `${crypto.randomUUID()}.png`;
    await saveScreenshot(filename, snapshot.screenshot);
    (snapshot as any).screenshotUrl = `http://${c.req.header('host')}/instant-screenshots/${filename}`;
  }

  if (options.favorScreenshot && snapshot.pageshot) {
    const filename = `${crypto.randomUUID()}.png`;
    await saveScreenshot(filename, snapshot.pageshot);
    (snapshot as any).pageshotUrl = `http://${c.req.header('host')}/instant-screenshots/${filename}`;
  }

  // Format response
  const response = formatSnapshot(snapshot, respondWith);

  // Handle screenshot/pageshot redirects
  if (respondWith === 'screenshot' && (snapshot as any).screenshotUrl) {
    return c.redirect((snapshot as any).screenshotUrl);
  }
  if (respondWith === 'pageshot' && (snapshot as any).pageshotUrl) {
    return c.redirect((snapshot as any).pageshotUrl);
  }

  // Return formatted content
  let content = '';
  if (response.content) content = response.content;
  else if (response.html) content = response.html;
  else if (response.text) content = response.text;
  else content = JSON.stringify(response, null, 2);

  return c.text(content, 200, {
    'Content-Type': respondWith === 'html' ? 'text/html' : 'text/plain; charset=utf-8'
  });
});

// Start server
const port = parseInt(Deno.env.get('PORT') || '3000');

async function startServer() {
  await ensureStorageDir();
  console.log(`Reader service starting on port ${port}...`);
  await Deno.serve({ port }, app.fetch);
}

startServer().catch(console.error);
