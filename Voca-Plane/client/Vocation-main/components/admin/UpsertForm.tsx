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
import { PlusCircle, Edit2, Loader2, Calendar as CalendarIcon, Trash2, Plus, Plane, ChevronLeft } from "lucide-react"

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
  triggerVariant?: "default" | "outline" | "ghost" | "icon" | "hidden"
  triggerIcon?: React.ReactNode
  maxWidth?: string
  columns?: 1 | 2
  onSubmit: (data: any) => Promise<void>
  onSuccess?: () => void
  open?: boolean
  onOpenChange?: (open: boolean) => void
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
  open: externalOpen,
  onOpenChange: setExternalOpen,
}: UpsertFormProps) {
  const [internalOpen, setInternalOpen] = useState(false)
  const open = externalOpen !== undefined ? externalOpen : internalOpen
  const setOpen = setExternalOpen !== undefined ? setExternalOpen : setInternalOpen
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
        {triggerVariant === "hidden" ? null : 
         triggerVariant === "icon" ? (
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
      <DialogContent className={cn("bg-white dark:bg-slate-900 border-none shadow-2xl overflow-hidden p-0 rounded-[3rem]", maxWidth)}>
        <div className="bg-slate-50 dark:bg-white/5 p-10 pb-8 border-b border-slate-100 dark:border-white/5 relative overflow-hidden">
          <div className="absolute top-0 right-0 p-10 opacity-10 rotate-12">
             <Plane className="size-24 text-primary" />
          </div>
          <DialogHeader className="relative z-10">
            <div className="bg-primary/20 size-14 rounded-2xl flex items-center justify-center text-primary mb-6 shadow-inner">
               {isUpdate ? <Edit2 className="size-6" /> : <PlusCircle className="size-6" />}
            </div>
            <DialogTitle className="text-4xl font-black tracking-tighter text-slate-900 dark:text-white">
              {isUpdate ? "Modify" : "Provision"} {title}
            </DialogTitle>
            <DialogDescription className="text-slate-500 font-medium text-lg mt-2 max-w-md">
              {description}
            </DialogDescription>
          </DialogHeader>
        </div>

        <ScrollArea className="max-h-[70vh]">
          <form onSubmit={handleSubmit} id="upsert-form" className="p-10 space-y-10">
            <div className={cn("grid gap-8", columns === 2 ? "grid-cols-1 md:grid-cols-2" : "grid-cols-1")}>
              {fields.map((field) => (
                <div 
                  key={field.name} 
                  className={cn(
                    "grid gap-3",
                    (field.type === "textarea" || field.type === "dynamic-list") && columns === 2 ? "md:col-span-2" : ""
                  )}
                >
                  <Label htmlFor={field.name} className="text-[11px] font-black uppercase tracking-[0.2em] text-slate-400 ml-2">
                    {field.label}
                    {field.required && <span className="text-primary ml-1">•</span>}
                  </Label>
                {field.type === "dynamic-list" ? (
                  <div className="space-y-4 p-8 rounded-[2.5rem] bg-slate-50 dark:bg-white/2 border border-slate-100 dark:border-white/5">
                    <div className="flex items-center justify-between mb-2">
                      <Label className="font-black text-slate-900 dark:text-white uppercase tracking-widest text-xs">Configuration Nodes</Label>
                      <Button 
                        type="button" 
                        variant="ghost" 
                        size="sm" 
                        className="h-10 px-4 rounded-xl gap-2 font-black text-[10px] uppercase tracking-widest hover:bg-white dark:hover:bg-white/5 transition-all"
                        onClick={() => {
                          const list = formData[field.name] || [];
                          setFormData({ ...formData, [field.name]: [...list, { class_type: "", price: 0 }] });
                        }}
                      >
                        <Plus className="h-4 w-4" />
                        Add Node
                      </Button>
                    </div>
                    {(formData[field.name] || []).map((item: any, idx: number) => (
                      <div key={idx} className="flex items-center gap-4 bg-white dark:bg-white/5 p-4 rounded-2xl border border-slate-100 dark:border-white/5 shadow-sm group">
                        <select
                          className="flex-1 h-12 rounded-xl border-none bg-slate-50 dark:bg-white/5 px-4 text-sm font-bold shadow-inner transition-all focus:ring-4 focus:ring-primary/20 appearance-none outline-none"
                          value={item.class_type}
                          onChange={(e) => {
                            const list = [...formData[field.name]];
                            list[idx] = { ...list[idx], class_type: e.target.value };
                            setFormData({ ...formData, [field.name]: list });
                          }}
                        >
                          <option value="">Select Tier</option>
                          <option value="First">First Class</option>
                          <option value="Business">Business</option>
                          <option value="Economy">Economy</option>
                        </select>
                        <div className="relative flex-[2]">
                           <span className="absolute left-4 top-1/2 -translate-y-1/2 font-black text-slate-300 text-xs">IDR</span>
                           <Input
                             type="number"
                             placeholder="0"
                             className="h-12 pl-12 bg-slate-50 dark:bg-white/5 border-none rounded-xl focus:ring-4 focus:ring-primary/20 font-black text-sm outline-none"
                             value={item.price || ""}
                             onChange={(e) => {
                               const list = [...formData[field.name]];
                               list[idx] = { ...list[idx], price: Number(e.target.value) };
                               setFormData({ ...formData, [field.name]: list });
                             }}
                           />
                        </div>
                        <Button 
                          type="button" 
                          variant="ghost" 
                          size="icon" 
                          className="h-12 w-12 rounded-xl text-slate-300 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10 transition-all opacity-0 group-hover:opacity-100"
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
                    <div className="flex gap-4">
                       <Popover>
                        <PopoverTrigger asChild>
                          <Button
                            variant={"outline"}
                            className={cn(
                              "flex-1 h-14 justify-start text-left font-bold rounded-2xl bg-slate-50 dark:bg-white/5 border-none focus:ring-4 focus:ring-primary/20 transition-all",
                              !formData[field.name] && "text-slate-300"
                            )}
                          >
                            <CalendarIcon className="mr-3 h-5 w-5 text-primary" />
                            {formData[field.name] ? (
                              format(new Date(formData[field.name]), "PPP")
                            ) : (
                              <span>Deployment Date</span>
                            )}
                          </Button>
                        </PopoverTrigger>
                        <PopoverContent className="w-auto p-0 rounded-3xl border-none shadow-3xl overflow-hidden" align="start">
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
                        className="w-32 h-14 bg-slate-50 dark:bg-white/5 border-none rounded-2xl focus:ring-4 focus:ring-primary/20 font-black text-sm outline-none"
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
                    <div className="space-y-4">
                      <div className="relative group">
                         <Button variant="outline" className="w-full h-32 border-dashed border-2 rounded-3xl flex flex-col items-center justify-center gap-2 text-slate-400 group-hover:border-primary group-hover:text-primary transition-all bg-slate-50/50 dark:bg-white/2">
                           <Loader2 className="size-8 opacity-20 group-hover:opacity-100 transition-opacity" />
                           <span className="font-black text-[10px] uppercase tracking-widest">Upload Digital Asset</span>
                         </Button>
                         <Input
                           id={field.name}
                           type="file"
                           accept="image/*"
                           className="absolute inset-0 opacity-0 cursor-pointer"
                           onChange={(e) => handleChange(field.name, e, "file")}
                           required={field.required && !isUpdate}
                         />
                      </div>
                      {/* Preview Gambar */}
                      {(preview || (isUpdate && formData[field.name])) && (
                        <div className="relative h-40 w-full rounded-3xl border border-slate-100 dark:border-white/10 overflow-hidden bg-white dark:bg-slate-800 shadow-inner p-4">
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
                      className="flex min-h-[160px] w-full rounded-2xl border-none bg-slate-50 dark:bg-white/5 px-6 py-4 text-sm font-medium shadow-inner transition-all focus:ring-4 focus:ring-primary/20 outline-none disabled:cursor-not-allowed disabled:opacity-50"
                      placeholder={field.placeholder}
                      value={formData[field.name] || ""}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                      required={field.required}
                    />
                  ) : field.type === "select" ? (
                    <div className="relative">
                       <select
                         id={field.name}
                         className="flex h-14 w-full rounded-2xl border-none bg-slate-50 dark:bg-white/5 px-6 py-2 text-sm font-bold shadow-inner transition-all focus:ring-4 focus:ring-primary/20 appearance-none outline-none disabled:cursor-not-allowed disabled:opacity-50"
                         value={formData[field.name] || ""}
                         onChange={(e) => handleChange(field.name, e.target.value)}
                         required={field.required}
                       >
                         <option value="" disabled className="text-slate-400">{field.placeholder || `Select ${field.label}`}</option>
                         {field.options?.map((opt) => (
                           <option key={opt.value} value={opt.value} className="bg-white dark:bg-slate-900 text-slate-900 dark:text-white">{opt.label}</option>
                         ))}
                       </select>
                       <div className="absolute right-6 top-1/2 -translate-y-1/2 pointer-events-none opacity-40">
                          <ChevronLeft className="-rotate-90 size-4" />
                       </div>
                    </div>
                  ) : (
                    <Input
                      id={field.name}
                      type={field.type || "text"}
                      placeholder={field.placeholder}
                      className="h-14 bg-slate-50 dark:bg-white/5 border-none rounded-2xl animate-none focus:ring-4 focus:ring-primary/20 font-bold text-sm outline-none shadow-inner"
                      value={formData[field.name] || ""}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                      required={field.required}
                    />
                  )}
              </div>
              ))}
            </div>
          </form>
          <ScrollBar orientation="vertical" />
        </ScrollArea>

        <div className="p-10 border-t border-slate-100 dark:border-white/5 bg-slate-50/30 dark:bg-white/2 flex flex-col sm:flex-row gap-4 items-center justify-between">
           <p className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">
             Security Clearance Level <span className="text-primary italic">Admin Alpha</span>
           </p>
           <div className="flex gap-4 w-full sm:w-auto">
             <Button
               type="button"
               variant="ghost"
               onClick={() => setOpen(false)}
               disabled={loading}
               className="flex-1 sm:flex-none h-14 px-8 font-black uppercase tracking-widest text-xs rounded-2xl hover:bg-white dark:hover:bg-white/5 transition-all"
             >
               Discard
             </Button>
             <Button
               type="submit"
               form="upsert-form"
               disabled={loading}
               className="flex-1 sm:flex-none h-14 px-12 font-black uppercase tracking-widest text-xs bg-primary text-white shadow-2xl shadow-primary/30 rounded-2xl transition-all hover:scale-[1.02] active:scale-95"
             >
               {loading ? (
                 <Loader2 className="h-5 w-5 animate-spin" />
               ) : (
                 "Execute Deployment"
               )}
             </Button>
           </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}

