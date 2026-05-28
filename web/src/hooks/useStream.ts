import { useState, useCallback, useRef } from 'react';

interface ComplianceIssue {
  issue_type: string;
  description: string;
  severity: string;
  location: string;
  suggestion: string;
}

interface StreamOptions {
  onChunk?: (chunk: string) => void;
  onIssue?: (issue: ComplianceIssue) => void;
  onComplete?: () => void;
  onError?: (error: Error) => void;
}

export function useStream(url: string, options: StreamOptions = {}) {
  const [content, setContent] = useState('');
  const [issues, setIssues] = useState<ComplianceIssue[]>([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const abortRef = useRef<AbortController | null>(null);

  const start = useCallback(async (body: Record<string, unknown>) => {
    abortRef.current = new AbortController();
    setIsStreaming(true);
    setContent('');
    setIssues([]);

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
        signal: abortRef.current.signal,
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body!.getReader();
      const decoder = new TextDecoder();

      while (true) {
        const { done, value } = await reader.read();
        if (done) break;

        const text = decoder.decode(value, { stream: true });
        const lines = text.split('\n').filter(l => l.startsWith('data: '));

        for (const line of lines) {
          const data = JSON.parse(line.slice(6));

          if (data.type === 'content') {
            setContent(prev => prev + data.chunk_text);
            options.onChunk?.(data.chunk_text);
          } else if (data.type === 'compliance_warning') {
            const newIssues = Array.isArray(data.issues) ? data.issues : [data.issues];
            setIssues(prev => [...prev, ...newIssues]);
            newIssues.forEach((issue: ComplianceIssue) => options.onIssue?.(issue));
          }
        }
      }

      options.onComplete?.();
    } catch (err) {
      if (err instanceof Error && err.name !== 'AbortError') {
        options.onError?.(err);
      }
    } finally {
      setIsStreaming(false);
    }
  }, [url, options]);

  const stop = useCallback(() => {
    abortRef.current?.abort();
    setIsStreaming(false);
  }, []);

  return { content, issues, isStreaming, start, stop };
}
