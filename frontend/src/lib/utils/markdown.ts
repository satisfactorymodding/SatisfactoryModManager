import DOMPurify from 'dompurify';
import { marked } from 'marked';

export const markdown = (md: string): string => {
  return DOMPurify.sanitize(marked(md));
};