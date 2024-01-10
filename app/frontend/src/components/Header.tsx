"use client";

import React, { useEffect } from "react";
import { Check, ChevronsUpDown } from "lucide-react";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

import { toast } from "./ui/use-toast";
import { ToastAction } from "./ui/toast";
import { useEditor, usePreference } from "@/lib/store";
import Preference from "@/components/Preference";
import { ScrollArea } from "./ui/scroll-area";
import { useStore } from "zustand";
import { Skeleton } from "./ui/skeleton";

const getFileList = async (token: string): Promise<string[]> => {
  const list: string[] | undefined = await fetch("/api/v1/file", {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  })
    .then((res) => {
      console.log(res.ok);
      if (!res.ok) {
        throw new Error("token 填写错误");
      }
      return res;
    })
    .then((res) => res.json())
    .then((res) => res.data)
    .catch((err) => {
      toast({
        variant: "destructive",
        title: err.message,
        action: <ToastAction altText="cancel">Cancel</ToastAction>,
      });
      console.log(err);
      return undefined;
    });
  if (list) {
    return list;
  }
  return [];
};

export function FileSelector() {
  const [loading, setLoading] = React.useState<boolean>(true);
  const [open, setOpen] = React.useState(false);
  const value = useStore(useEditor, (s) => s.file);
  const [files, setFiles] = React.useState<string[]>([]);
  const token = usePreference((s) => s.token);
  const setFile = useEditor((s) => s.Set);

  useEffect(() => {
    if (token) {
      getFileList(token)
        .then(setFiles)
        .then(() => setLoading(false));
    }
  }, [token]);

  if (loading) {
    return <Skeleton className="w-full h-10" />;
  }
  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="justify-between w-full"
        >
          {value ? files.find((file) => file === value) : "Select File..."}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full p-0">
        <Command>
          <CommandInput placeholder="Search file..." />
          <CommandEmpty>No framework found.</CommandEmpty>
          <CommandGroup className="pr-0">
            <ScrollArea className="h-72">
              {files.map((file) => (
                <CommandItem
                  key={file}
                  value={file}
                  onSelect={(currentValue) => {
                    setOpen(false);
                    setFile({ file: currentValue, content: "" });
                  }}
                >
                  <Check
                    className={cn(
                      "mr-2 h-4 w-4",
                      value === file ? "opacity-100" : "opacity-0"
                    )}
                  />
                  {file}
                </CommandItem>
              ))}
            </ScrollArea>
          </CommandGroup>
        </Command>
      </PopoverContent>
    </Popover>
  );
}

export default function Header() {
  return (
    <div className="w-full flex flex-row items-center gap-x-4 p-2 sticky top-0 z-50 bg-white border-b mb-2">
      <div className="flex-1">
        <FileSelector />
      </div>
      <div className="pr-2">
        <Preference />
      </div>
    </div>
  );
}
