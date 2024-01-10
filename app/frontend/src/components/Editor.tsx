"use client";

import React, { useEffect, useState } from "react";
import CodeMirror, { ViewUpdate } from "@uiw/react-codemirror";
import { markdown, markdownLanguage } from "@codemirror/lang-markdown";
import { languages } from "@codemirror/language-data";
import { useEditor, usePreference, useStore } from "@/lib/store";
import { toast } from "./ui/use-toast";
import { ToastAction } from "./ui/toast";
import { Skeleton } from "./ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";

import Markdown from "react-markdown";
import remarkGfm from "remark-gfm";
import Mdx from "./mdx/mdx";

const getFileContent = async (
  filename: string,
  token: string
): Promise<string> => {
  const content: string | undefined = await fetch(
    `/api/v1/file/${btoa(filename)}`,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  )
    .then((res) => {
      if (!res.ok) {
        if (res.status === 401) {
          throw new Error("token 填写错误");
        }
        return undefined;
      }
      return res.text();
    })
    .catch((err) => {
      toast({
        variant: "destructive",
        title: err.message,
        action: <ToastAction altText="cancel">Cancel</ToastAction>,
      });
      return undefined;
    });
  if (!content) {
    return "";
  }
  return content;
};

export function Preview() {
  const content = useStore(useEditor, (state) => state.content);
  const Foo = () => <p> foo bar foo bar</p>;
  return (
    <>
      {/* <Markdown remarkPlugins={[]}>{content}</Markdown> */}
      <Mdx content={content}></Mdx>
    </>
  );
}

export function DocEditor() {
  const file = useStore(useEditor, (state) => state.file);
  const token = usePreference((s) => s.token);
  const [doc, setDoc] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(true);
  const setEditor = useEditor((state) => state.Set);
  const getContent = useEditor((s) => s.GetContent);

  useEffect(() => {
    if (file === undefined) return;
    const str = getContent();
    if (str.length !== 0) {
      setLoading(false);
      setDoc(str);
      return;
    }
    getFileContent(file ? file : "", token).then((str) => {
      setLoading(false);
      setDoc(str);
    });
  }, [file]);

  const onChange = (doc: string, viewUpdate: ViewUpdate) => {
    setEditor({ content: doc });
  };

  const _codeMirrorEditor: React.FC<{
    initDoc: string;
    className?: string;
  }> = ({ initDoc, className }) => {
    return (
      <CodeMirror
        value={initDoc}
        extensions={[
          markdown({ base: markdownLanguage, codeLanguages: languages }),
        ]}
        onChange={onChange}
        className={cn("w-full h-screen", className)}
      />
    );
  };
  const CodeMirrorEditor = React.memo(_codeMirrorEditor);

  return (
    <>
      {loading ? (
        <div className="w-full">
          <Skeleton className="w-full h-screen " />
        </div>
      ) : (
        <div className="w-full ">
          <CodeMirrorEditor
            initDoc={doc}
            className="rounded-md h-[calc(100vh-10rem)] overflow-y-scroll"
          />
        </div>
      )}
    </>
  );
}

export default function Editor() {
  return (
    <>
      <Tabs defaultValue="editor" className="w-full px-2">
        <TabsList className="grid grid-cols-2">
          <TabsTrigger value="editor">Editor</TabsTrigger>
          <TabsTrigger value="preview">Preview</TabsTrigger>
        </TabsList>
        <TabsContent value="editor">
          <DocEditor />
        </TabsContent>
        <TabsContent value="preview">
          <Preview />
        </TabsContent>
      </Tabs>
    </>
  );
}
