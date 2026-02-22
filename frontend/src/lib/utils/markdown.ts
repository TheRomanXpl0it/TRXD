import { marked } from 'marked';
import DOMPurify from 'dompurify';

marked.setOptions({
	breaks: true,
	gfm: true
});

export function renderMarkdown(markdown: string): string {
	if (!markdown) return '';

	const rawHtml = marked.parse(markdown, { async: false }) as string;

	return DOMPurify.sanitize(rawHtml, {
		ALLOWED_TAGS: [
			'p',
			'br',
			'strong',
			'em',
			'u',
			's',
			'del',
			'ins',
			'h1',
			'h2',
			'h3',
			'h4',
			'h5',
			'h6',
			'ul',
			'ol',
			'li',
			'blockquote',
			'pre',
			'code',
			'a',
			'img',
			'table',
			'thead',
			'tbody',
			'tr',
			'th',
			'td',
			'hr',
			'sup',
			'sub'
		],
		ALLOWED_ATTR: ['href', 'src', 'alt', 'title', 'target', 'rel'],
		ALLOW_DATA_ATTR: false
	});
}

export function renderMarkdownInline(markdown: string): string {
	if (!markdown) return '';

	const html = renderMarkdown(markdown);

	// this might not be safe
	return html.replace(/^<p>([\s\S]*)<\/p>$/, '$1');
}
