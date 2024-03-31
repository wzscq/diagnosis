import { redirect } from 'next/navigation'
import {getVechicleStatus} from "@/api";

export default async function vehicleStatus(){
    const data = await getVechicleStatus();
    redirect(`/vehiclestatus/${data?.result?.list?.[0]?.id}`)
}