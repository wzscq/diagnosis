import { redirect } from 'next/navigation'
import Chart from '@/components/Chart';
import StatusCard from '@/components/StatusCard';
import {getVechicleStatus} from "@/api";
import styles from './page.module.css';

const option = {
    title:{
      text:"能耗曲线",
      left:'center',
      top:20,
      textStyle:{
        color:'#606060',
        fontWeight:500,
        fontSize:20
      }
    },
    xAxis: {
      name:'Time（s）',
      nameLocation:'middle',
      nameGap:30
    },
    yAxis: {
      name:'能耗（kw）',
      nameLocation:'middle',
      nameGap:50
    },
    series: [
      {
        data: [
          [10, 50],
          [50, 100],
          [60, 20],
          [70, 60],
          [80, 30],
          [90, 100],
          [170, 160],
          [270, 260],
          [370, 560],
          [470, 160],
          [570, 260],
          [670, 460],
          [770, 60],
        ],
        type: 'line'
      }
    ]
  };

export default async function vehicleStatus({params}){
  const result=await getVechicleStatus({id:params.statusid})
  const row=result?.result?.list?.[0]??{};

  return (
    <div className={styles.statusGrid}>
      <div className={styles.statusChart}>
        <Chart chartOption={option} />
      </div>
      <StatusCard title={'最大车速'} value={(row['speed_max']??'')+'km/h'} />
      <StatusCard title={'平均车速'} value={(row['speed_avg']??'')+'km/h'} />
      <StatusCard title={'行驶里程'} value={(row['mileage']??'')+'km'} />
      <StatusCard title={'行驶时间'} value={(row['travel_time']??'')+'s'} />
      <StatusCard title={'平均能耗'} value={(row['ec_avg']??'')+'kw'} />
    </div>
  )
}