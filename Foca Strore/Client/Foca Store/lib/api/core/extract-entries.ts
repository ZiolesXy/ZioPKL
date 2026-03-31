export function extractEntries<T>(payload: unknown, label: string): T[] {
    if (
        typeof payload !== "object" || payload === null || !("data" in payload)
    ){
        throw new Error(`Invalid ${label} response`)
    }

const data = (payload as any).data

if(!Array.isArray(data?.entries)){
    throw new Error(`Invalid ${label} entries`)
}

return data.entries as T[]

}