"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { putData } from "@/lib/api/resources/put";

type PutButtonProps = {
  id: number | string;
  endpoint: string;
  defaultValues: Record<string, any>;
  onSuccess?: () => void;
};

export default function PutButton({
  id,
  endpoint,
  defaultValues,
  onSuccess,
}: PutButtonProps) {
  const [form, setForm] = useState(defaultValues);
  const [loading, setLoading] = useState(false);
  const [open, setOpen] = useState(false);
  const router = useRouter();

  const handleOpenChange = (next: boolean) => {
    setOpen(next);
    if (next) {
      setForm(defaultValues);
    }
  };

  const handleChange = (key: string, value: string) => {
    const isNumeric = key.toLowerCase().endsWith("_id");
    setForm({
      ...form,
      [key]: isNumeric ? (value === "" ? "" : Number(value)) : value,
    });
  };

  const handleSubmit = async () => {
    try {
      setLoading(true);
      await putData(endpoint, id, form);
      onSuccess?.();
      setOpen(false);
      router.refresh();
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button variant="outline">Edit</Button>
      </DialogTrigger>

      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Data</DialogTitle>
        </DialogHeader>

        <div className="space-y-3">
          {Object.keys(form).map((key) => (
            <Input
              key={key}
              placeholder={key}
              type={key.toLowerCase().endsWith("_id") ? "number" : "text"}
              value={form[key] ?? ""}
              onChange={(e) => handleChange(key, e.target.value)}
            />
          ))}

          <Button onClick={handleSubmit} disabled={loading}>
            {loading ? "Saving..." : "Save"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}