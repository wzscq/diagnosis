import { getDrivingStyle } from "@/api"
import Chart from "@/components/Chart"
import { Card} from "antd"
import DrivingStyleData from "@/components/DrivingStyleData"
import Toolbar from './Toolbar'
import styles from './page.module.css'

const waterfallOption = {
    title:{
        text:'温和型',
        right:20,
        textStyle:{
            color:'#1677FF'
        }
    },
    grid: {
      left: '20',
      right: '20',
      top:40,
      bottom: '100',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      axisLabel:{
         rotate:45,
         color:"#000"
      },
      data: ['车速','油门深度','油门变化率','刹车深度','刹车变化率','转向速度','侧向加速度','纵向加速度','','']
    },
    yAxis: {
      type: 'value',
      interval:5,
      min:0,
      max:40
    },
    series: [
      {
        type: 'bar',
        stack: 'Total',
        silent: true,
        itemStyle: {
          borderColor: 'transparent',
          color: 'transparent'
        },
        emphasis: {
          itemStyle: {
            borderColor: 'transparent',
            color: 'transparent'
          }
        },
        data: [0, 0, 0, 0, 0, 0,0,0]
      },
      {
        name: 'Income',
        type: 'bar',
        stack: 'Total',
        label: {
          show: true,
          position: 'top'
        },
        data: [0, 0, 0, 0, 0, 0,0,0],
        markArea:{
            silent:true,
            data:[
                [
                    {
                        name: '保守   \n\n0≤∑≤21',
                        itemStyle:{
                            
                        },
                        label:{
                            position:'insideRight',
                            color:'#000'
                        },
                        yAxis:0
                    },
                    {
                        yAxis: 21
                    }
                ],
                [
                    {
                        name: '温和   \n\n22≤∑≤35',
                        itemStyle:{
                            color:'#1677FF',
                            opacity:0.5
                        },
                        label:{
                            position:'insideRight',
                            color:'#000'
                        },
                        yAxis: 21
                    },
                    {
                        yAxis: 35
                    }
                ],
                [
                    {
                        name: '激进   \n\n36≤∑≤40',
                        itemStyle:{
                            color:'#124D9F',
                            opacity:0.5
                        },
                        label:{
                            position:'insideRight',
                            color:'#000'
                        },
                        yAxis: 35
                    },
                    {
                        yAxis: 40
                    }
                ]
            ]
        }
      },
    ]
}


const radarOption = {
    radar: {
      // shape: 'circle',
      indicator: [
        { name: '车速', max: 5,color:'#000' },
        { name: '油门深度', max: 5,color:'#000' },
        { name: '油门变化率', max: 5,color:'#000' },
        { name: '刹车深度', max: 5,color:'#000' },
        { name: '刹车变化率', max: 5,color:'#000' },
        { name: '转向速度', max: 5,color:'#000' },
        { name: '侧向加速度', max: 5,color:'#000' },
        { name: '纵向加速度', max: 5,color:'#000' }
      ]
    },
    series: [
      {
        name: '',
        type: 'radar',
        label:{
          show:true,
          color:'#000'
        },
        data: [
          {
            value: [0, 0, 0, 0, 0, 0,0,0]
          }
        ]
      }
    ]
  };

export default async function DrivingStylePage({params}){
    const result=await getDrivingStyle(params.styleid)
    const row=result?.result?.list?.[0]??{};

    if(row){
        const value=[
            row['speed_max_leve']??0,
            row['mtpd_ratio_level']??0,
            row['mtc_ratio_level']??0,
            row['mbpd_ratio_level']??0,
            row['mbc_ratio_level']??0,
            row['mtwav_level']??0,
            row['lateral_acc_max_level']??0,
            row['long_acc_max_level']??0
        ]
        radarOption.series[0].data=[{
            value:[...value]
        }]
        const gradeValue=[...value]
        let total=0;
        for(let i=0;i<gradeValue.length;i++){
            gradeValue[i]=parseInt(gradeValue[i])+total
            total=gradeValue[i]
        }
        waterfallOption.series[0].data=[...gradeValue]
        waterfallOption.series[1].data=[...value]
    }

    if(params.direction!=='h'){
      waterfallOption.xAxis.axisLabel.rotate=0
    }

    return (
        <div className={styles.main}>
            <div className={styles.header}>
                <div className={styles.title}>驾驶风格分析</div>
                <div className={styles.subTitle}>车辆编号：{row?.['vehicle_code']??''}</div>
                <div className={styles.toolbar}>
                    <Toolbar styleid={params.styleid}  direction={params.direction}/>
                </div>
            </div>
            <div className={params.direction==='h'?styles.contentH:styles.contentV}>
                <Card size='small' title="数据统计" style={{height:'100%'}}>
                    <DrivingStyleData data={row}/>
                </Card>
                <Card size='small' title="驾驶指数" style={{height:'100%'}}>
                    <Chart chartOption={radarOption} />
                </Card>
                <Card size='small' title="驾驶风格" style={{height:'100%'}}>
                    <Chart chartOption={waterfallOption} />
                </Card>
            </div>
        </div>
    )
}