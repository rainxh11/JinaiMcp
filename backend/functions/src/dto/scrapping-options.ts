import { marshalErrorLike, RPCHost, RPCReflection, HashManager, AssertionFailureError, ParamValidationError } from 'civkit';
import _ from 'lodash';

export class CrawlerOptions {
    respondWith?: 'markdown' | 'html' | 'text' | 'screenshot' | 'pageshot' = 'markdown';
    withGeneratedAlt?: boolean;
    withLinksSummary?: boolean;
    withImagesSummary?: boolean;
    keepImgDataUrl?: boolean;
    cacheTolerance?: string | number;
    userAgent?: string;
    timeout?: number;
    html?: string;
    proxyUrl?: string;
    removeSelector?: string | string[];
    targetSelector?: string | string[];
    waitForSelector?: string | string[];
    withIframe?: boolean;

    static from(query: any, req?: any): CrawlerOptions {
        const options = new CrawlerOptions();

        options.respondWith = query['X-Respond-With'] || query.respondWith || 'markdown';
        options.withGeneratedAlt = query['X-With-Generated-Alt'] === 'true' || query.withGeneratedAlt === 'true';
        options.withLinksSummary = query['X-With-Links-Summary'] === 'true' || query.withLinksSummary === 'true';
        options.withImagesSummary = query['X-With-Images-Summary'] === 'true' || query.withImagesSummary === 'true';
        options.keepImgDataUrl = query['X-Keep-Img-Data-Url'] === 'true' || query.keepImgDataUrl === 'true';
        options.cacheTolerance = query['X-Cache-Tolerance'] || query.cacheTolerance;
        options.userAgent = query['X-User-Agent'] || query['User-Agent'] || query.userAgent;
        options.timeout = parseInt(query['X-Timeout'] || query.timeout) || undefined;
        options.html = query.html;
        options.proxyUrl = query['X-Proxy-Url'] || query.proxyUrl;
        options.removeSelector = query['X-Remove-Selector'] || query.removeSelector;
        options.targetSelector = query['X-Target-Selector'] || query.targetSelector;
        options.waitForSelector = query['X-Wait-For-Selector'] || query.waitForSelector;
        options.withIframe = query['X-With-Iframe'] === 'true' || query.withIframe === 'true';

        return options;
    }
}

export class CrawlerOptionsHeaderOnly {
    respondWith?: 'markdown' | 'html' | 'text' | 'screenshot' | 'pageshot' = 'markdown';
    withGeneratedAlt?: boolean;
    withLinksSummary?: boolean;
    withImagesSummary?: boolean;
    keepImgDataUrl?: boolean;
    cacheTolerance?: string | number;
    userAgent?: string;

    static from(req: any): CrawlerOptionsHeaderOnly {
        const options = new CrawlerOptionsHeaderOnly();

        options.respondWith = req.headers['x-respond-with'] || 'markdown';
        options.withGeneratedAlt = req.headers['x-with-generated-alt'] === 'true';
        options.withLinksSummary = req.headers['x-with-links-summary'] === 'true';
        options.withImagesSummary = req.headers['x-with-images-summary'] === 'true';
        options.keepImgDataUrl = req.headers['x-keep-img-data-url'] === 'true';
        options.cacheTolerance = req.headers['x-cache-tolerance'];
        options.userAgent = req.headers['x-user-agent'] || req.headers['user-agent'];

        return options;
    }
}
