"use client";

import { useEffect, useRef } from "react";
import { EditorView, basicSetup } from "codemirror";
import { EditorState } from "@codemirror/state";
import { oneDark } from "@codemirror/theme-one-dark";
import { javascript } from "@codemirror/lang-javascript";

interface CodeEditorProps {
  value: string;
  onChange: (value: string) => void;
}

export function CodeEditor({ value, onChange }: CodeEditorProps) {
  const editorRef = useRef<HTMLDivElement>(null);
  const viewRef = useRef<EditorView | null>(null);

  useEffect(() => {
    if (!editorRef.current) return;

    const state = EditorState.create({
      doc: value,
      extensions: [
        basicSetup,
        javascript(),
        oneDark,
        EditorView.updateListener.of((update) => {
          if (update.docChanged) {
            onChange(update.state.doc.toString());
          }
        }),
        EditorView.theme({
          "&": {
            height: "100%",
          },
          ".cm-scroller": {
            fontFamily:
              "ui-monospace, SFMono-Regular, 'SF Mono', Consolas, 'Liberation Mono', Menlo, monospace",
          },
          ".cm-focused": {
            outline: "none",
          },
        }),
      ],
    });

    const view = new EditorView({
      state,
      parent: editorRef.current,
    });

    viewRef.current = view;

    return () => {
      view.destroy();
    };
  }, []);

  useEffect(() => {
    if (viewRef.current && viewRef.current.state.doc.toString() !== value) {
      viewRef.current.dispatch({
        changes: {
          from: 0,
          to: viewRef.current.state.doc.length,
          insert: value,
        },
      });
    }
  }, [value]);

  return (
    <div className="h-full flex flex-col">
      <div className="border-b border-border bg-muted/50 px-4 py-2">
        <div className="flex items-center justify-between">
          <span className="text-sm font-medium">Contract.sol</span>
          <div className="flex items-center space-x-2 text-xs text-muted-foreground">
            <span>Solidity</span>
            <span>•</span>
            <span>UTF-8</span>
          </div>
        </div>
      </div>
      <div ref={editorRef} className="flex-1 overflow-hidden" />
    </div>
  );
}
