import { redirect } from 'next/navigation'
import Chart from '@/components/Chart';
import StatusCard from '@/components/StatusCard';
import {getVechicleStatus,getVechicleStatusDetail} from "@/api";
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
    tooltip:{
      trigger:'axis'
    },
    xAxis: {
      name:'Time（s）',
      nameLocation:'middle',
      nameGap:30
    },
    yAxis: {
      name:'能耗（kwh/100km）',
      nameLocation:'middle',
      nameGap:50
    },
    series: [
      {
        data: [],
        type: 'line'
      }
    ]
  };

export default async function vehicleStatus({params}){
  let filter = {};
  if(params.statusid !== '-1'){
    filter={id:params.statusid}
  }
  const result=await getVechicleStatus(1,filter)
  const row=result?.result?.list?.[0]??{};

  const detail=await getVechicleStatusDetail(params.statusid)
  let data=[]
  let yAxisMin=-1;
  let yAxisMax=-1;
  let xAxisMax=-1;
  if(detail?.length>0){
    const detailRow=detail[0];
    if(detailRow){
      data=detailRow.SignalCoordinateValue.map((item)=>{
        let val=parseFloat(item.Coordinate_Y['$numberDouble'])
        if(yAxisMin===-1||val<yAxisMin){
          yAxisMin=val
        }
        if(yAxisMax===-1||val>yAxisMax){
          yAxisMax=val
        }
        return [parseFloat(item.Coordinate_X),val]
      })
      xAxisMax=data[data.length-1][0];
    }
  }

  //yAxisMin+2向下取整
  yAxisMin=Math.floor(yAxisMin)-1
  //yAxisMax+2向上取整
  yAxisMax=Math.ceil(yAxisMax)+1
  //console.log('yAxisMin',yAxisMin)
  //console.log('yAxisMax',yAxisMax)
  option.yAxis.min=yAxisMin;
  option.yAxis.max=yAxisMax;
  option.xAxis.max=xAxisMax;

  option.series[0].data=data;

  const travel_time=row['travel_time']??0
  //console.log('travel_time',travel_time)
  const travelFloat = parseInt(travel_time)
  const travelMin = Math.floor(travelFloat / 60)
  let travelTime = travelMin>0?travelMin + 'min':''
  const travelSec = travelFloat % 60
  travelTime = travelTime +' '+travelSec + 's'

  return (
    <div className={styles.statusGrid}>
      <div className={styles.statusChart}>
        <Chart chartOption={{...option}} />
      </div>
      <StatusCard title={'NOP激活行驶里程'} value={(row['nop_mileage']??'')+'km'} />
      <StatusCard title={'NOP退出次数'} value={(row['nop_disabledTimes']??'')} />
      <StatusCard title={'TJA激活行驶里程'} value={(row['tja_mileage']??'')+'km'} />
      <StatusCard title={'ACC激活行驶里程'} value={(row['acc_mileage']??'')+'km'} />
      <StatusCard title={'最大车速'} value={(row['speed_max']??'')+'km/h'} />
      <StatusCard title={'平均车速'} value={(row['speed_avg']??'')+'km/h'} />
      <StatusCard title={'平均能耗'} value={(row['ec_avg']??'')+'kwh/100km'} />
      <StatusCard title={'行驶时间'} value={travelTime} />
      <StatusCard title={'行驶里程'} value={(row['mileage']??'')+'km'} />
    </div>
  )
}