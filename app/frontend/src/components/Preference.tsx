"use client";

import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { usePreference, useStore } from "@/lib/store";
import { Settings } from "lucide-react";
import { Label } from "./ui/label";
import { Input } from "./ui/input";
import { Button } from "./ui/button";
import { useState } from "react";

export default function Preference() {
  const token = useStore(usePreference, (state) => state.token);
  const setPreference = usePreference((state) => state.Set);

  type FormData = {
    token: string | undefined;
  };

  const [form, setForm] = useState<FormData>({ token: token });

  return (
    <Dialog>
      <DialogTrigger>
        <Settings className="h-6 w-6" />
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Preference</DialogTitle>
          <DialogDescription>设置</DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <Label htmlFor="token" className="text-right">
              token
            </Label>
            <Input
              id="token"
              className="col-span-3"
              value={form.token}
              onChange={(e) => {
                setForm({ ...form, token: e.target.value });
              }}
            />
          </div>
        </div>
        <DialogFooter>
          <Button onClick={() => setPreference(form)}>Save changes</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
