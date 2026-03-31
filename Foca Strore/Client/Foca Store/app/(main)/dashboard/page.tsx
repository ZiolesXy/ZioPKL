import { ChartLineStep } from "@/components/LineChart"
import { ChartBarLabelCustom} from "@/components/BarChart"
import CardTable from "@/components/CardTable"
export default async function UsersPage() {

  return (
    <>
      <div className="w-full p-5">
        <div className="grid min-w-0 grid-cols-1 gap-5 md:grid-cols-2">
          <div className="min-w-0">
            <ChartBarLabelCustom />
          </div>
          <div className="min-w-0">
            <ChartLineStep />
          </div>
          <div className="min-w-0 md:col-span-2">
            <CardTable />
          </div>
        </div>
      </div>
   </>
  )
}
