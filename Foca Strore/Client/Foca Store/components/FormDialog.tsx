"use client"

import { useMemo, useState, useTransition, type FormEvent } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"
import { Field, FieldGroup } from "@/components/ui/field"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export type DialogFieldType = "text" | "number"

export type DialogField = {
  name: string
  label: string
  type?: DialogFieldType
  placeholder?: string
  defaultValue?: string | number
  required?: boolean
}

export default function FormDialog<TValues extends Record<string, any>>({
  triggerLabel,
  title,
  description,
  fields,
  submitLabel = "Save",
  onSubmit,
  onSuccess,
}: {
  triggerLabel: string
  title: string
  description?: string
  fields: DialogField[]
  submitLabel?: string
  onSubmit: (values: TValues) => Promise<void>
  onSuccess?: () => void
}) {
  const [open, setOpen] = useState(false)
  const [isPending, startTransition] = useTransition()
  const router = useRouter()

  const initialValues = useMemo(() => {
    const v: Record<string, any> = {}
    fields.forEach((f) => {
      v[f.name] = f.defaultValue ?? ""
    })
    return v as TValues
  }, [fields])

  const [values, setValues] = useState<TValues>(initialValues)

  const handleOpenChange = (next: boolean) => {
    setOpen(next)
    if (next) {
      setValues(initialValues)
    }
  }

  const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    startTransition(async () => {
      try {
        await onSubmit(values)
        onSuccess?.()
        setOpen(false)
        router.refresh()
      } catch (err) {
        const message = err instanceof Error ? err.message : "Insert gagal"
        alert(message)
      }
    })
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button variant="outline">{triggerLabel}</Button>
      </DialogTrigger>

      <DialogContent className="sm:max-w-sm">
        <form onSubmit={handleSubmit}>
          <DialogHeader>
            <DialogTitle>{title}</DialogTitle>
            {description ? <DialogDescription>{description}</DialogDescription> : null}
          </DialogHeader>

          <FieldGroup className="mt-4">
            {fields.map((f) => {
              const inputType = f.type ?? "text"

              return (
                <Field key={f.name}>
                  <Label htmlFor={f.name}>{f.label}</Label>
                  <Input
                    id={f.name}
                    name={f.name}
                    type={inputType}
                    placeholder={f.placeholder}
                    required={f.required}
                    value={values[f.name] ?? ""}
                    onChange={(e) => {
                      const raw = e.target.value
                      const nextValue = inputType === "number" ? (raw === "" ? "" : Number(raw)) : raw
                      setValues((prev) => ({ ...prev, [f.name]: nextValue }))
                    }}
                  />
                </Field>
              )
            })}
          </FieldGroup>

          <DialogFooter className="mt-6">
            <DialogClose asChild>
              <Button type="button" variant="outline" disabled={isPending}>
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit" disabled={isPending}>
              {isPending ? "Saving..." : submitLabel}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
