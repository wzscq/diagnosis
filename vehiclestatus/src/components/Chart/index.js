'use client'
import React, {useEffect,useRef} from 'react';
import * as echarts from 'echarts';
import { useResizeDetector } from 'react-resize-detector';

import styles from './page.module.css'

export default function Chart({chartOption}){
    const refChart=useRef();
    const { width,ref,height } = useResizeDetector();

    useEffect(()=>{
        if(refChart&&refChart.current&&chartOption!==null){
            let chart=echarts.getInstanceByDom(refChart.current);        
            if(chart){
                chart.dispose();
            }
            chart=echarts.init(refChart.current,'light');
            chart.setOption(chartOption);
        }
    },
    [refChart,chartOption]);

    useEffect(()=>{
        if(refChart&&refChart.current){
            let chart=echarts.getInstanceByDom(refChart.current);        
            if(chart){
                chart.resize({width:width,height:height});
            }
        }
    },
    [refChart,width,height]);

    return (
        <div className={styles.wrapper}>
            <div ref={refChart} style={{width:'100%',height:'100%'}} />
            <div ref={ref} style={{width:'100%',height:'100%'}}>{}</div>
        </div>
    );
}