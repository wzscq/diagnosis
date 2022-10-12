import { useEffect,useRef } from 'react';
import { useResizeDetector } from 'react-resize-detector';
import * as echarts from 'echarts';

import { setAnalysisChart } from '../../../redux/reportSlice';
import { useDispatch } from 'react-redux';

export default function SignalChart({keyIndex,chartIndex,signal}){
    const dispatch=useDispatch();
    const refChart=useRef();
    const { width,ref } = useResizeDetector();
    const {SignalCoordinateValue,SignalName}=signal;
    
    console.log("width:",width,SignalCoordinateValue);

    useEffect(()=>{

        const option = {
            title: {
                text: SignalName,
                left: '1%'
            },
            tooltip: {
                trigger: 'axis'
            },
            grid: {
                left: '10%',
                right: '5%',
                bottom: '15%',
                top:40
            },
            xAxis: {
                data: SignalCoordinateValue.map(function (item) {
                    return item.Coordinate_X;
                })
            },
            yAxis: {},
            series: {
            type: 'line',
            data: SignalCoordinateValue.map(function (item) {
                return parseInt(item.Coordinate_Y);
            }),
            }
        };

        if(refChart&&refChart.current){
            console.log("chart width:",width);
            let chart=echarts.getInstanceByDom(refChart.current);        
            if(chart){
                chart.resize({width:width});
            } else {
                chart=echarts.init(refChart.current,'dark',{
                    width: width,
                    height: 150
                });
                chart.setOption(option);
            }
            dispatch(setAnalysisChart({keyIndex:keyIndex,chartIndex:chartIndex,chart:chart.getDataURL({type:'png'})}));
        }
    },[SignalCoordinateValue,SignalName,refChart,width,keyIndex,chartIndex,dispatch]);

    return (
        <>
            <div ref={refChart} style={{marginTop:5,width:width,marginBottom:20,marginLeft:0,height:150}}/>
            <div ref={ref}>{}</div>
        </>
    )
}