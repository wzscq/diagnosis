'use client'
import Chart from '../Chart';
import StatusCard from '../StatusCard';
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

export default function vehicleStatus(){
    return (
    <div className={styles.statusGrid}>
        <div style={{gridColumn:'1 / span 5;',gridRow:'1 / span 2;'}}>
          <Chart chartOption={option} />
        </div>
        <StatusCard title={'最大车速'} value={'53.5 km/h'} />
        <StatusCard title={'平均车速'} value={'23.5 km/h'} />
        <StatusCard title={'行驶里程'} value={'3.5 km'} />
        <StatusCard title={'行驶时间'} value={'35min26s'} />
        <StatusCard title={'平均能耗'} value={'17.5 kw'} />
    </div>
    )
}