import puppeteer from 'puppeteer-core';
import type { Browser, Page } from 'puppeteer-core';
import { JSDOM } from 'jsdom';
import TurndownService from 'turndown';

export interface ScrappingOptions {
  favorScreenshot?: boolean;
  waitForSelector?: string | string[];
  timeoutMs?: number;
  overrideUserAgent?: string;
}

export interface PageSnapshot {
  title: string;
  href: string;
  html: string;
  text: string;
  parsed?: {
    title?: string;
    content?: string;
    textContent?: string;
    publishedTime?: string;
  };
  screenshot?: Buffer;
  pageshot?: Buffer;
  error?: string;
}

export interface ReaderResponse {
  title?: string;
  url?: string;
  content?: string;
  html?: string;
  text?: string;
  screenshotUrl?: string;
  pageshotUrl?: string;
  error?: string;
}

export class PuppeteerService {
  private browser?: Browser;
  private page?: Page;
  private chromePath: string;

  constructor() {
    this.chromePath = Deno.env.get('PUPPETEER_EXECUTABLE_PATH') || '/usr/bin/google-chrome-stable';
  }

  async init(): Promise<void> {
    if (this.browser) return;

    this.browser = await puppeteer.launch({
      executablePath: this.chromePath,
      args: [
        '--no-sandbox',
        '--disable-setuid-sandbox',
        '--disable-dev-shm-usage',
        '--single-process',
        '--disable-gpu'
      ],
      timeout: 10000
    });

    const context = await this.browser.createBrowserContext();
    this.page = await context.newPage();

    await this.page.setBypassCSP(true);
    await this.page.setViewport({ width: 1920, height: 1080 });
  }

  async scrap(url: string, options: ScrappingOptions = {}): Promise<PageSnapshot> {
    if (!this.page) {
      await this.init();
    }

    const page = this.page!;
    const timeout = options.timeoutMs || 30000;

    try {
      // Set custom user agent if provided
      if (options.overrideUserAgent) {
        await page.setUserAgent(options.overrideUserAgent);
      }

      // Navigate to URL
      await page.goto(url, {
        waitUntil: ['load', 'domcontentloaded'],
        timeout
      });

      // Wait for selector if specified
      if (options.waitForSelector) {
        const selectors = Array.isArray(options.waitForSelector)
          ? options.waitForSelector
          : [options.waitForSelector];

        for (const selector of selectors) {
          try {
            await page.waitForSelector(selector, { timeout: 5000 });
          } catch {
            // Continue if selector not found
          }
        }
      }

      // Get page content
      const result = await page.evaluate(() => {
        // Try to use Readability if available
        let parsed = null;
        if (typeof (window as any).Readability !== 'undefined') {
          try {
            const clone = document.cloneNode(true);
            parsed = new (window as any).Readability(clone).parse();
          } catch {
            // Readability failed, continue without it
          }
        }

        return {
          title: document.title,
          href: document.location.href,
          html: document.documentElement.outerHTML,
          text: document.body.innerText,
          parsed
        };
      });

      // Take screenshots
      let screenshot: Buffer | undefined;
      let pageshot: Buffer | undefined;

      if (options.favorScreenshot) {
        screenshot = Buffer.from(await page.screenshot() as Uint8Array);
        pageshot = Buffer.from(await page.screenshot({ fullPage: true }) as Uint8Array);
      }

      return {
        ...result,
        screenshot,
        pageshot
      };

    } catch (error: any) {
      return {
        title: 'Error',
        href: url,
        html: '',
        text: error.message || 'Unknown error',
        error: error.message
      };
    }
  }

  async close(): Promise<void> {
    if (this.browser) {
      await this.browser.close();
      this.browser = undefined;
      this.page = undefined;
    }
  }
}

export function formatSnapshot(snapshot: PageSnapshot, mode: string): ReaderResponse {
  const response: ReaderResponse = {
    title: snapshot.parsed?.title || snapshot.title,
    url: snapshot.href
  };

  if (mode === 'html') {
    response.html = snapshot.html;
  } else if (mode === 'text') {
    response.text = snapshot.text;
  } else if (mode === 'markdown') {
    const jsdom = new JSDOM(snapshot.html, { url: snapshot.href });
    const doc = jsdom.window.document;

    // Simple markdown conversion
    let markdown = `# ${snapshot.title}\n\n`;

    // Convert paragraphs
    doc.querySelectorAll('p').forEach(p => {
      markdown += `${p.textContent}\n\n`;
    });

    // Convert links
    doc.querySelectorAll('a').forEach(a => {
      const href = a.getAttribute('href');
      const text = a.textContent;
      if (href && text) {
        markdown += `[${text}](${href})`;
      }
    });

    // Convert headings
    doc.querySelectorAll('h1, h2, h3, h4, h5, h6').forEach(h => {
      const level = parseInt(h.tagName[1]);
      const prefix = '#'.repeat(level);
      markdown += `\n${prefix} ${h.textContent}\n\n`;
    });

    response.content = markdown;
  } else if (mode === 'screenshot') {
    response.screenshotUrl = `/instant-screenshots/${crypto.randomUUID()}.png`;
  } else if (mode === 'pageshot') {
    response.pageshotUrl = `/instant-screenshots/${crypto.randomUUID()}.png`;
  }

  return response;
}
