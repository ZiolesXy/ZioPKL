"use client"

import { useState, useEffect } from "react"
import { toast } from "sonner"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { format } from "date-fns"
import { cn } from "@/lib/utils"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { Calendar } from "@/components/ui/calendar"
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area"
import { PlusCircle, Edit2, Loader2, Calendar as CalendarIcon, Trash2, Plus } from "lucide-react"

export interface FormField {
  name: string
  label: string
  type?: "text" | "number" | "url" | "select" | "textarea" | "file" | "dynamic-list" | "datetime"
  placeholder?: string
  required?: boolean
  options?: { label: string; value: string | number }[] // For select type
}

interface UpsertFormProps {
  title: string
  description: string
  fields: FormField[]
  initialData?: any // If provided, it's Update mode
  triggerLabel?: string
  triggerVariant?: "default" | "outline" | "ghost" | "icon"
  triggerIcon?: React.ReactNode
  maxWidth?: string
  columns?: 1 | 2
  onSubmit: (data: any) => Promise<void>
  onSuccess?: () => void
}

export function UpsertForm({
  title,
  description,
  fields,
  initialData,
  triggerLabel,
  triggerVariant = "default",
  triggerIcon,
  maxWidth = "sm:max-w-[500px]",
  columns = 1,
  onSubmit,
  onSuccess,
}: UpsertFormProps) {
  const [open, setOpen] = useState(false)
  const [loading, setLoading] = useState(false)
  const [formData, setFormData] = useState<any>({})
  const [preview, setPreview] = useState<string | null>(null)

  // Reset form when initialData changes or modal opens
  useEffect(() => {
    if (open) {
      if (initialData) {
        setFormData(initialData)
        // If there's a file field (e.g., logo), we might set preview if it's a URL
        // However, FlightsTableData doesn't use file field for now.
        // For general use, if a field is named 'logo_url' or similar, we could preview it.
      } else {
        setFormData({})
      }
    } else {
      setPreview(null)
    }
  }, [open, initialData])

  const handleChange = (name: string, value: any, type?: string) => {
    if (type === "file") {
      const file = value.target.files?.[0]
      if (file) {
        setFormData((prev: any) => ({ ...prev, [name]: file }))
        setPreview(URL.createObjectURL(file)) // Buat URL sementara untuk preview
      }
    } else {
      setFormData((prev: any) => ({ ...prev, [name]: value }))
    }
  }

  const handleDynamicChange = (fieldName: string, index: number, subField: string, value: any) => {
    setFormData((prev: any) => {
      const list = [...(prev[fieldName] || [])];
      list[index] = { ...list[index], [subField]: value };
      return { ...prev, [fieldName]: list };
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await onSubmit(formData)
      toast.success(`${title} berhasil disimpan`)
      setOpen(false)
      if (onSuccess) onSuccess()
    } catch (error: any) {
      console.error("Gagal menyimpan data:", error)
      toast.error(`Gagal menyimpan ${title}: ${error.message || "Terjadi kesalahan"}`)
    } finally {
      setLoading(false)
    }
  }

  const isUpdate = !!initialData

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {triggerVariant === "icon" ? (
          <Button variant="outline" size="icon" className="h-8 w-8">
            {triggerIcon || <Edit2 className="h-4 w-4" />}
          </Button>
        ) : (
          <Button variant={triggerVariant} className="gap-2">
            {triggerIcon || (isUpdate ? <Edit2 className="h-4 w-4" /> : <PlusCircle className="h-4 w-4" />)}
            {triggerLabel || (isUpdate ? "Edit" : "Tambah Baru")}
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className={cn("bg-card border-none shadow-2xl overflow-hidden p-0", maxWidth)}>
        <div className="bg-primary/5 p-6 pb-2">
          <DialogHeader>
            <DialogTitle className="text-2xl font-bold tracking-tight">
              {isUpdate ? "Update" : "Create New"} {title}
            </DialogTitle>
            <DialogDescription className="text-muted-foreground mt-1">
              {description}
            </DialogDescription>
          </DialogHeader>
        </div>

        <ScrollArea className="max-h-[85vh]">
          <form onSubmit={handleSubmit} className="p-6 space-y-6">
            <div className={cn("grid gap-5", columns === 2 ? "grid-cols-1 md:grid-cols-2" : "grid-cols-1")}>
              {fields.map((field) => (
                <div 
                  key={field.name} 
                  className={cn(
                    "grid gap-2",
                    (field.type === "textarea" || field.type === "dynamic-list") && columns === 2 ? "md:col-span-2" : ""
                  )}
                >
                  <Label htmlFor={field.name} className="text-sm font-semibold flex gap-1">
                  {field.label}
                  {field.required && <span className="text-destructive">*</span>}
                </Label>
                {field.type === "dynamic-list" ? (
                  <div className="space-y-3 p-4 border rounded-lg bg-muted/20">
                    <div className="flex items-center justify-between">
                      <Label className="text-base">{field.label}</Label>
                      <Button 
                        type="button" 
                        variant="outline" 
                        size="sm" 
                        className="h-8 gap-1"
                        onClick={() => {
                          const list = formData[field.name] || [];
                          setFormData({ ...formData, [field.name]: [...list, { class_type: "", price: 0 }] });
                        }}
                      >
                        <Plus className="h-4 w-4" />
                        Tambah Class
                      </Button>
                    </div>
                    {(formData[field.name] || []).map((item: any, idx: number) => (
                      <div key={idx} className="flex items-center gap-2 bg-background p-2 rounded-md border">
                        <select
                          className="w-32 h-9 rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                          value={item.class_type}
                          onChange={(e) => {
                            const list = [...formData[field.name]];
                            list[idx] = { ...list[idx], class_type: e.target.value };
                            setFormData({ ...formData, [field.name]: list });
                          }}
                        >
                          <option value="">Pilih Kelas</option>
                          <option value="First">First</option>
                          <option value="Business">Business</option>
                          <option value="Economy">Economy</option>
                        </select>
                        <Input
                          type="number"
                          placeholder="Harga"
                          className="flex-1 h-9"
                          value={item.price || ""}
                          onChange={(e) => {
                            const list = [...formData[field.name]];
                            list[idx] = { ...list[idx], price: Number(e.target.value) };
                            setFormData({ ...formData, [field.name]: list });
                          }}
                        />
                        <Button 
                          type="button" 
                          variant="ghost" 
                          size="icon" 
                          className="h-9 w-9 text-destructive hover:text-destructive/80 hover:bg-destructive/10"
                          onClick={() => {
                            const list = [...formData[field.name]];
                            list.splice(idx, 1);
                            setFormData({ ...formData, [field.name]: list });
                          }}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    ))}
                  </div>
                ) :
                  field.type === "datetime" ? (
                    <div className="flex flex-wrap gap-2">
                       <Popover>
                        <PopoverTrigger asChild>
                          <Button
                            variant={"outline"}
                            className={cn(
                              "w-full sm:w-[240px] h-11 justify-start text-left font-normal",
                              !formData[field.name] && "text-muted-foreground"
                            )}
                          >
                            <CalendarIcon className="mr-2 h-4 w-4" />
                            {formData[field.name] ? (
                              format(new Date(formData[field.name]), "PPP")
                            ) : (
                              <span>Pilih Tanggal</span>
                            )}
                          </Button>
                        </PopoverTrigger>
                        <PopoverContent className="w-auto p-0" align="start">
                          <Calendar
                            mode="single"
                            selected={formData[field.name] ? new Date(formData[field.name]) : undefined}
                            onSelect={(date) => {
                              if (!date) return;
                              const current = formData[field.name] ? new Date(formData[field.name]) : new Date();
                              date.setHours(current.getHours());
                              date.setMinutes(current.getMinutes());
                              setFormData({ ...formData, [field.name]: date.toISOString() });
                            }}
                            initialFocus
                          />
                        </PopoverContent>
                      </Popover>
                      <Input
                        type="time"
                        className="w-full sm:w-[120px] h-11"
                        value={formData[field.name] ? format(new Date(formData[field.name]), "HH:mm") : ""}
                        onChange={(e) => {
                          const [hours, minutes] = e.target.value.split(":").map(Number);
                          const current = formData[field.name] ? new Date(formData[field.name]) : new Date();
                          current.setHours(hours);
                          current.setMinutes(minutes);
                          setFormData({ ...formData, [field.name]: current.toISOString() });
                        }}
                      />
                    </div>
                  ) :
                  field.type === "file" ? (
                    <div className="space-y-3">
                      <Input
                        id={field.name}
                        type="file"
                        accept="image/*"
                        className="h-11 cursor-pointer pt-2"
                        onChange={(e) => handleChange(field.name, e, "file")}
                        required={field.required && !isUpdate} // Jika update, tidak wajib isi ulang
                      />
                      {/* Preview Gambar */}
                      {(preview || (isUpdate && formData[field.name])) && (
                        <div className="relative h-24 w-24 rounded-lg border overflow-hidden bg-muted">
                          <img
                            src={preview || formData[field.name]}
                            alt="Preview"
                            className="h-full w-full object-contain"
                          />
                        </div>
                      )}
                    </div>
                  ) : field.type === "textarea" ? (
                    <textarea
                      id={field.name}
                      className="flex min-h-[100px] w-full rounded-lg border border-input bg-transparent px-3 py-2 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                      placeholder={field.placeholder}
                      value={formData[field.name] || ""}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                      required={field.required}
                    />
                  ) : field.type === "select" ? (
                    <select
                      id={field.name}
                      className="flex h-10 w-full rounded-lg border border-input bg-transparent px-3 py-2 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 appearance-none"
                      value={formData[field.name] || ""}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                      required={field.required}
                    >
                      <option value="" disabled>{field.placeholder || `Pilih ${field.label}`}</option>
                      {field.options?.map((opt) => (
                        <option key={opt.value} value={opt.value}>{opt.label}</option>
                      ))}
                    </select>
                  ) : (
                    <Input
                      id={field.name}
                      type={field.type || "text"}
                      placeholder={field.placeholder}
                      className="h-11 rounded-lg border-muted-foreground/20 focus:border-primary transition-all shadow-none"
                      value={formData[field.name] || ""}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                      required={field.required}
                    />
                  )}
              </div>
            ))}
          </div>

          <DialogFooter className="pt-4 border-t gap-3 sm:gap-0 mt-6">
            <Button
              type="button"
              variant="ghost"
              onClick={() => setOpen(false)}
              disabled={loading}
              className="h-11 px-6 font-medium"
            >
              Batal
            </Button>
            <Button
              type="submit"
              disabled={loading}
              className="h-11 px-8 font-bold shadow-lg shadow-primary/20 transition-all hover:scale-[1.02] active:scale-[0.98]"
            >
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Menyimpan...
                </>
              ) : (
                "Simpan Data"
              )}
            </Button>
          </DialogFooter>
          </form>
          <ScrollBar orientation="vertical" />
        </ScrollArea>
      </DialogContent>
    </Dialog>
  )
}

