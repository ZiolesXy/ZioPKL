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
  options?: { label: string; value: string | number }[]
}

interface UpsertFormProps {
  title: string
  description: string
  fields: FormField[]
  initialData?: any
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
  maxWidth = "sm:max-w-[600px]",
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

  useEffect(() => {
    if (open) {
      if (initialData) {
        setFormData(initialData)
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
        setPreview(URL.createObjectURL(file))
      }
    } else {
      setFormData((prev: any) => ({ ...prev, [name]: value }))
    }
  }

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
      
      {/* Penyesuaian: max-h-screen dan rounded-none di mobile agar tidak terpotong */}
      <DialogContent className={cn(
        "bg-white dark:bg-slate-900 border-none shadow-2xl p-0 flex flex-col overflow-hidden",
        "w-full h-full sm:h-[92vh] sm:rounded-3xl rounded-none min-h-0",
        maxWidth
      )}>
        
        {/* Header: Padding lebih kecil di mobile */}
        <div className="bg-slate-50 dark:bg-white/5 p-5 md:p-6 pb-3 border-b border-slate-100 dark:border-white/5 relative overflow-hidden shrink-0">
          <div className="absolute top-0 right-0 p-10 opacity-10 rotate-12 hidden sm:block">
             <Plane className="size-24 text-primary" />
          </div>
          <DialogHeader className="relative z-10 text-left">
            <div className="bg-primary/20 size-10 rounded-xl flex items-center justify-center text-primary mb-3 shadow-inner">
               {isUpdate ? <Edit2 className="size-5" /> : <PlusCircle className="size-5" />}
            </div>
            <DialogTitle className="text-xl md:text-2xl font-black tracking-tighter text-slate-900 dark:text-white">
              {isUpdate ? "Modify" : "Provision"} {title}
            </DialogTitle>
            <DialogDescription className="text-slate-500 font-medium text-xs md:sm mt-1 max-w-md">
              {description}
            </DialogDescription>
          </DialogHeader>
        </div>

        {/* Area Scroll: min-h-0 penting untuk flex-1 agar scroll aktif */}
        <ScrollArea className="flex-1 w-full min-h-0">
          <form onSubmit={handleSubmit} id="upsert-form" className="p-5 md:p-8 space-y-6">
            <div className={cn("grid gap-5 md:gap-8", columns === 2 ? "grid-cols-1 md:grid-cols-2" : "grid-cols-1")}>
              {fields.map((field) => (
                <div 
                  key={field.name} 
                  className={cn(
                    "grid gap-2 md:gap-3",
                    (field.type === "textarea" || field.type === "dynamic-list") && columns === 2 ? "md:col-span-2" : ""
                  )}
                >
                  <Label htmlFor={field.name} className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400 ml-1">
                    {field.label}
                    {field.required && <span className="text-primary ml-1">•</span>}
                  </Label>

                {field.type === "dynamic-list" ? (
                  <div className="space-y-4 p-4 md:p-8 rounded-[1.5rem] md:rounded-[2.5rem] bg-slate-50 dark:bg-white/2 border border-slate-100 dark:border-white/5">
                    <div className="flex items-center justify-between mb-2">
                      <Label className="font-black text-slate-900 dark:text-white uppercase tracking-widest text-[10px]">Configuration Nodes</Label>
                      <Button 
                        type="button" 
                        variant="ghost" 
                        size="sm" 
                        className="h-8 md:h-10 px-3 md:px-4 rounded-xl gap-2 font-black text-[9px] md:text-[10px] uppercase tracking-widest hover:bg-white dark:hover:bg-white/5 transition-all"
                        onClick={() => {
                          const list = formData[field.name] || [];
                          setFormData({ ...formData, [field.name]: [...list, { class_type: "", price: 0 }] });
                        }}
                      >
                        <Plus className="h-3 w-3 md:h-4 md:w-4" />
                        Add Node
                      </Button>
                    </div>
                    {(formData[field.name] || []).map((item: any, idx: number) => (
                      <div key={idx} className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 bg-white dark:bg-white/5 p-3 md:p-4 rounded-2xl border border-slate-100 dark:border-white/5 shadow-sm group relative">
                        <select
                          className="flex-1 h-11 md:h-12 rounded-xl border-none bg-slate-50 dark:bg-white/5 px-4 text-sm font-bold appearance-none outline-none focus:ring-2 focus:ring-primary/20"
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
                           <span className="absolute left-4 top-1/2 -translate-y-1/2 font-black text-slate-300 text-[10px]">IDR</span>
                           <Input
                             type="number"
                             placeholder="0"
                             className="h-11 md:h-12 pl-12 bg-slate-50 dark:bg-white/5 border-none rounded-xl focus:ring-2 focus:ring-primary/20 font-black text-sm outline-none"
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
                          className="h-10 w-full sm:w-10 rounded-xl text-red-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10 sm:opacity-0 sm:group-hover:opacity-100 transition-all"
                          onClick={() => {
                            const list = [...formData[field.name]];
                            list.splice(idx, 1);
                            setFormData({ ...formData, [field.name]: list });
                          }}
                        >
                          <Trash2 className="h-4 w-4" />
                          <span className="sm:hidden ml-2 text-[10px] font-bold uppercase">Remove Node</span>
                        </Button>
                      </div>
                    ))}
                  </div>
                ) :
                  field.type === "datetime" ? (
                    <div className="flex flex-col sm:flex-row gap-3">
                        <Popover>
                        <PopoverTrigger asChild>
                          <Button
                            variant={"outline"}
                            className={cn(
                              "flex-1 h-12 md:h-14 justify-start text-left font-bold rounded-2xl bg-slate-50 dark:bg-white/5 border-none focus:ring-2 focus:ring-primary/20",
                              !formData[field.name] && "text-slate-300"
                            )}
                          >
                            <CalendarIcon className="mr-3 h-5 w-5 text-primary" />
                            {formData[field.name] ? format(new Date(formData[field.name]), "PPP") : <span>Deployment Date</span>}
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
                          />
                        </PopoverContent>
                      </Popover>
                      <Input
                        type="time"
                        className="w-full sm:w-32 h-12 md:h-14 bg-slate-50 dark:bg-white/5 border-none rounded-2xl focus:ring-2 focus:ring-primary/20 font-black text-sm outline-none"
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
                  ) : field.type === "file" ? (
                    <div className="space-y-4">
                      <div className="relative group">
                         <Button variant="outline" className="w-full h-24 md:h-32 border-dashed border-2 rounded-2xl md:rounded-3xl flex flex-col items-center justify-center gap-2 text-slate-400 bg-slate-50/50 dark:bg-white/2">
                           <Loader2 className="size-6 md:size-8 opacity-20" />
                           <span className="font-black text-[9px] md:text-[10px] uppercase tracking-widest">Upload Asset</span>
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
                      {(preview || (isUpdate && formData[field.name])) && (
                        <div className="relative h-32 md:h-40 w-full rounded-2xl border border-slate-100 dark:border-white/10 overflow-hidden bg-white dark:bg-slate-800 p-4">
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
                      className="flex min-h-[120px] md:min-h-[160px] w-full rounded-2xl border-none bg-slate-50 dark:bg-white/5 px-6 py-4 text-sm font-medium focus:ring-2 focus:ring-primary/20 outline-none"
                      placeholder={field.placeholder}
                      value={formData[field.name] || ""}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                      required={field.required}
                    />
                  ) : field.type === "select" ? (
                    <div className="relative">
                       <select
                         id={field.name}
                         className="flex h-12 md:h-14 w-full rounded-2xl border-none bg-slate-50 dark:bg-white/5 px-6 py-2 text-sm font-bold appearance-none outline-none focus:ring-2 focus:ring-primary/20"
                         value={formData[field.name] || ""}
                         onChange={(e) => handleChange(field.name, e.target.value)}
                         required={field.required}
                       >
                         <option value="" disabled>{field.placeholder || `Select ${field.label}`}</option>
                         {field.options?.map((opt) => (
                           <option key={opt.value} value={opt.value}>{opt.label}</option>
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
                      className="h-12 md:h-14 bg-slate-50 dark:bg-white/5 border-none rounded-2xl focus:ring-2 focus:ring-primary/20 font-bold text-sm outline-none"
                      value={formData[field.name] || ""}
                      onChange={(e) => handleChange(field.name, e.target.value)}
                      required={field.required}
                    />
                  )}
              </div>
              ))}
            </div>
          </form>
        </ScrollArea>

        {/* Footer: Disesuaikan agar tombol tidak bertumpuk di mobile */}
        <div className="p-5 md:p-6 border-t border-slate-100 dark:border-white/5 bg-slate-50/30 dark:bg-white/2 flex flex-col sm:flex-row gap-4 items-center justify-between shrink-0">
           <p className="text-[9px] md:text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">
             Clearance <span className="text-primary">Admin Alpha</span>
           </p>
           <div className="flex gap-3 w-full sm:w-auto">
             <Button
               type="button"
               variant="ghost"
               onClick={() => setOpen(false)}
               disabled={loading}
               className="flex-1 sm:flex-none h-12 md:h-14 px-6 md:px-8 font-black uppercase text-[10px] rounded-2xl"
             >
               Discard
             </Button>
             <Button
               type="submit"
               form="upsert-form"
               disabled={loading}
               className="flex-1 sm:flex-none h-12 md:h-14 px-8 md:px-12 font-black uppercase text-[10px] bg-primary text-white shadow-xl shadow-primary/20 rounded-2xl"
             >
               {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : "Execute"}
             </Button>
           </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}