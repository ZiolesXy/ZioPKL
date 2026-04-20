"use client"

import { useEffect, useMemo, useState } from "react"

import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { Calendar } from "@/components/ui/calendar"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"
import { Calendar as CalendarIcon, ImageIcon } from "lucide-react"
import Image from "next/image"
import { cn } from "@/lib/utils"


export type GenericEditFieldType = "text" | "number" | "textarea" | "file" | "select" | "date" 


// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type GenericEditField<TValues extends Record<string, any>> = {
    name: keyof TValues & string
    label: string
    type: GenericEditFieldType
    placeholder?: string
    required?: boolean
    accept?: string
    hint?: string
    min?: number
    options?: Array<{ value: string; label: string }>
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
interface GenericEditDialogProps<TValues extends Record<string, any>> {
    open: boolean
    onOpenChange: (open: boolean) => void
    title: string
    description?: string
    submitLabel?: string
    initialValues: TValues
    fields: Array<GenericEditField<TValues>>
    onSubmit: (values: TValues) => Promise<{ success: true } | { success: false; message: string }>
    onSuccess?: () => void
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function GenericEditDialog<TValues extends Record<string, any>>({
    open,
    onOpenChange,
    title,
    description,
    submitLabel = "Simpan Perubahan",
    initialValues,
    fields,
    onSubmit,
    onSuccess,
}: GenericEditDialogProps<TValues>) {
    const [values, setValues] = useState<TValues>(initialValues)
    const [fileValues, setFileValues] = useState<Record<string, File | null>>({})
    const [previews, setPreviews] = useState<Record<string, string>>( {})
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [error, setError] = useState("")

    const fileFields = useMemo(() => fields.filter((f) => f.type === "file"), [fields])

    useEffect(() => {
        setValues(initialValues)
        setError("")
        setIsSubmitting(false)
        const nextPreviews: Record<string, string> = {}
        for (const f of fileFields) {
            const key = f.name
            const maybeUrl = (initialValues as Record<string, unknown>)[key]
            if (typeof maybeUrl === "string") nextPreviews[key] = maybeUrl
        }
        setPreviews(nextPreviews)
        const nextFiles: Record<string, File | null> = {}
        for (const f of fileFields) nextFiles[f.name] = null
        setFileValues(nextFiles)
    }, [initialValues, fileFields])

    const setFieldValue = (name: string, value: unknown) => {
        setValues((prev) => ({ ...prev, [name]: value }))
    }

    const handleFileChange = (name: string, e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0]
        if (!file) return
        setFileValues((prev) => ({ ...prev, [name]: file }))
        setPreviews((prev) => ({ ...prev, [name]: URL.createObjectURL(file) }))
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        setIsSubmitting(true)
        setError("")
        try {
            const submitValues = { ...values } as Record<string, unknown>
            for (const [k, v] of Object.entries(fileValues)) {
                submitValues[k] = v
            }
            const result = await onSubmit(submitValues as unknown as TValues)
            if (!result.success) {
                setError(result.message)
                return
            }
            onSuccess?.()
            onOpenChange(false)
        } catch {
            setError("Terjadi kesalahan tidak dikenal.")
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-137.5 max-h-[90vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>{title}</DialogTitle>
                    {description ? (
                        <DialogDescription>{description}</DialogDescription>
                    ) : null}
                </DialogHeader>

                <form onSubmit={handleSubmit} className="space-y-5 mt-2">
                    {fields.map((field) => {
                        const fieldName = field.name
                        const id = `edit-${fieldName}`
                        const value = (values as Record<string, unknown>)[fieldName]
                        const stringValue = typeof value === "string" || typeof value === "number" ? value : (value ?? "") as string

                        if (field.type === "textarea") {
                            return (
                                <div key={fieldName} className="space-y-2">
                                    <Label htmlFor={id}>{field.label}</Label>
                                    <Textarea
                                        id={id}
                                        value={stringValue}
                                        onChange={(e) => setFieldValue(fieldName, e.target.value)}
                                        placeholder={field.placeholder}
                                        className="min-h-25"
                                        required={field.required}
                                    />
                                </div>
                            )
                        }

                        if (field.type === "file") {
                            const preview = previews[fieldName]
                            return (
                                <div key={fieldName} className="space-y-2">
                                    <Label>{field.label}</Label>
                                    <div className="flex items-start gap-4">
                                        <div className="h-24 w-24 rounded-md border-2 border-dashed flex items-center justify-center overflow-hidden bg-muted/50 shrink-0">
                                            {preview ? (
                                                <Image
                                                    src={preview}
                                                    alt="Preview"
                                                    width={96}
                                                    height={96}
                                                    className="object-cover w-full h-full"
                                                />
                                            ) : (
                                                <ImageIcon className="h-8 w-8 text-muted-foreground" />
                                            )}
                                        </div>
                                        <div className="flex-1 space-y-1">
                                            <Input
                                                id={id}
                                                type="file"
                                                accept={field.accept}
                                                onChange={(e) => handleFileChange(fieldName, e)}
                                                className="cursor-pointer"
                                                required={field.required}
                                            />
                                            {field.hint ? (
                                                <p className="text-xs text-muted-foreground">{field.hint}</p>
                                            ) : null}
                                        </div>
                                    </div>
                                </div>
                            )
                        }

                        if (field.type === "select") {
                            return (
                                <div key={fieldName} className="space-y-2">
                                    <Label htmlFor={id}>{field.label}</Label>
                                    <select
                                        id={id}
                                        className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background"
                                        value={stringValue}
                                        onChange={(e) => setFieldValue(fieldName, e.target.value)}
                                        required={field.required}
                                    >
                                        {field.placeholder ? (
                                            <option value="" disabled>
                                                {field.placeholder}
                                            </option>
                                        ) : null}
                                        {(field.options ?? []).map((opt) => (
                                            <option key={opt.value} value={opt.value}>
                                                {opt.label}
                                            </option>
                                        ))}
                                    </select>
                                </div>
                            )
                        }

                        if (field.type === "date") {
                            const selectedDate = typeof value === "string" && value ? new Date(value) : undefined
                            const hasValidDate = selectedDate instanceof Date && !Number.isNaN(selectedDate.getTime())
                            return (
                                <div key={fieldName} className="space-y-2">
                                    <Label>{field.label}</Label>
                                    <Popover>
                                        <PopoverTrigger asChild>
                                            <Button
                                                type="button"
                                                variant="outline"
                                                className={cn(
                                                    "w-full justify-start text-left font-normal",
                                                    !hasValidDate && "text-muted-foreground"
                                                )}
                                            >
                                                <CalendarIcon className="mr-2 h-4 w-4" />
                                                {hasValidDate ? selectedDate.toLocaleDateString("id-ID") : "Pilih tanggal"}
                                            </Button>
                                        </PopoverTrigger>
                                        <PopoverContent className="w-auto p-0" align="start">
                                            <Calendar
                                                mode="single"
                                                selected={hasValidDate ? selectedDate : undefined}
                                                onSelect={(d) => {
                                                    if (!d) return
                                                    setFieldValue(fieldName, d.toISOString())
                                                }}
                                                initialFocus
                                            />
                                        </PopoverContent>
                                    </Popover>
                                </div>
                            )
                        }

                        const isNumber = field.type === "number"
                        return (
                            <div key={fieldName} className="space-y-2">
                                <Label htmlFor={id}>{field.label}</Label>
                                <Input
                                    id={id}
                                    type={isNumber ? "number" : "text"}
                                    value={stringValue}
                                    onChange={(e) => setFieldValue(fieldName, e.target.value)}
                                    placeholder={field.placeholder}
                                    required={field.required}
                                    min={field.min}
                                />
                            </div>
                        )
                    })}

                    {error && (
                        <div className="p-3 text-sm font-medium text-destructive bg-destructive/10 rounded-md">
                            {error}
                        </div>
                    )}

                    <DialogFooter className="gap-2 pt-2">
                        <Button
                            type="button"
                            variant="outline"
                            onClick={() => onOpenChange(false)}
                            disabled={isSubmitting}
                        >
                            Batal
                        </Button>
                        <Button
                            type="submit"
                            className="bg-teal-600 hover:bg-teal-700"
                            disabled={isSubmitting}
                        >
                            {isSubmitting ? "Menyimpan..." : submitLabel}
                        </Button>
                    </DialogFooter>
                </form>
            </DialogContent>
        </Dialog>
    )
}
