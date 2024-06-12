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
                        name: '温和   \n\n0≤∑<20',
                        itemStyle:{
                            
                        },
                        label:{
                            position:'insideRight',
                            color:'#000'
                        },
                        yAxis:0
                    },
                    {
                        yAxis: 20
                    }
                ],
                [
                    {
                        name: '标准   \n\n20≤∑≤35',
                        itemStyle:{
                            color:'#1677FF',
                            opacity:0.5
                        },
                        label:{
                            position:'insideRight',
                            color:'#000'
                        },
                        yAxis: 20
                    },
                    {
                        yAxis: 35
                    }
                ],
                [
                    {
                        name: '激进   \n\n35<∑≤40',
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
      radius:'70%',
      nameGap:10,
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
          parseInt(row['speed_max_level']??0),
          parseInt(row['mtpd_ratio_level']??0),
          parseInt(row['mtc_ratio_level']??0),
          parseInt(row['mbpd_ratio_level']??0),
          parseInt(row['mbc_ratio_level']??0),
          parseInt(row['mtwav_level']??0),
          parseInt(row['lateral_acc_max_level']??0),
          parseInt(row['long_acc_max_level']??0)
        ]
        radarOption.series[0].data=[{
            value:[...value]
        }]
        const gradeValue=[0, 0, 0, 0, 0, 0,0,0]
        let total=0;
        for(let i=0;i<gradeValue.length;i++){
            gradeValue[i]=total
            total=value[i]+total
        }
        waterfallOption.series[0].data=[...gradeValue]
        waterfallOption.series[1].data=[...value]

        console.log('gradeValue',gradeValue,"value",value)

        waterfallOption.title.text=row['driving_style']??''
    }

    if(params.direction!=='h'){
      waterfallOption.xAxis.axisLabel.rotate=0
    } else {
      waterfallOption.xAxis.axisLabel.rotate=45
    }

    return (
        <div className={styles.main}>
            <div className={styles.header}>
                <div className={styles.title}>驾驶风格分析</div>
                <div className={styles.subTitle}></div>
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