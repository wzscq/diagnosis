import {useRef,useMemo,useEffect} from 'react';
import * as echarts from 'echarts';
import { useDispatch } from 'react-redux';
import { useResizeDetector } from 'react-resize-detector';
import { getColor } from '../../../utils/colorPalette';
import { setAnalysisChart } from '../../../redux/reportSlice';

export default function SignalChart({itmeIndex,signalList}){
    const dispatch=useDispatch();
    const refChart=useRef();
    const { width,height,ref } = useResizeDetector();
    
    const option=useMemo(()=>{
        //数据收敛
        const legendData=signalList.map(item=>item.SignalName+(item.SignalUint?'('+item.SignalUint+')':''));
    
        let maxX=-1;
        let minX=-1;
        const seriesData= signalList.map((item,index)=>{     
            return ({
                name:item.SignalName+(item.SignalUint?'('+item.SignalUint+')':''),
                type: 'line',
                smoth:true,
                yAxisIndex:index,
                symbolSize:1,
                lineStyle:{
                    width:1
                },
                data:item.SignalCoordinateValue.map(valueItem=>{
                    const x=parseFloat(valueItem.Coordinate_X);
                    if(minX===-1||minX>x){
                        minX=x;
                    }
                    if(maxX===-1||maxX<x){
                        maxX=x;
                    }
                    const y=parseFloat(valueItem.Coordinate_Y);     
                    return ([x,y]);
                })
            });
        });

        const grid= {
            left: (signalList.length*40+'px'),
            right: '10px'
        };
                
        const xAxis={
            type: 'value',
            min:Math.floor(minX),
            max:Math.ceil(maxX)
        };

        const yAxis=signalList.map((item,index)=>{
            return {
                type: 'value',
                name: "",
                position: 'left',
                offset: index*40,
                alignTicks: true,
                axisLine: {
                    show: true,
                    lineStyle: {
                        color: getColor(index)
                    }
                },
                axisLabel: {
                    formatter: '{value}'
                }
            }
        });

        const getLabel=(axisValue)=>{
            const values = seriesData.map(element => {
                return {
                    label:element.name,
                    value:element.data.reduce((prev,current)=>{
                            if(prev&&current[0]>axisValue){
                                return prev;
                            }
                            return current[1];
                        },false)
                    }
            });

            return values.reduce((prev,current,index)=>{
                return `${prev}
                    <div style="display:inline-block;width:15px;height:15px;background:${getColor(index)};margin-right:5px;"></div>
                    <span style="line-height: 15px;font-size: 12px;vertical-align:top;margin-right:5px;">
                        ${current.label}:
                    </span>
                    <span style="line-height: 15px;font-size: 14px;vertical-align:top; color:blue" >${current.value}</span>
                    <br/>`;
            },'');
        }

        const option = {
            xAxis: xAxis,
            grid: grid,
            yAxis: yAxis,
            tooltip: {
                trigger: 'axis',
                axisPointer: {
                    type: 'cross'
                },
                formatter:(params)=>{
                    const label=getLabel(params[0].axisValue);
                    return label;
                },
            },
            legend: {
                data: legendData,
                right:100,
                top:5
            },
            toolbox: {
                feature: {
                    dataZoom: {
                        yAxisIndex: 'none'
                    },
                    saveAsImage: {}
                }
            },
            dataZoom: [
                {
                    type: 'inside',
                    start: 0,
                    end: 100,
                },
                {
                    start: 0,
                    end: 100,
                }
            ],
            series: seriesData
        };
        return option;
    },[signalList]);

    useEffect(()=>{
        if(refChart&&refChart.current){
            let chart=echarts.getInstanceByDom(refChart.current);        
            if(chart){
                chart.resize({width:width,height:height});     
            } else {
                chart=echarts.init(refChart.current,'',{
                    width: width,
                    height: height
                });
                chart.setOption(option);   
            }
            if(width>0&&height>0){
                dispatch(setAnalysisChart({itmeIndex:itmeIndex,chartIndex:0,chart:chart.getDataURL({type:'png'})}));
            }
        }
    },[dispatch,refChart,width,height,option,itmeIndex]);
    
    return (
        <div style={{width:'100%',height:'400px',overflow:"hidden"}}>
            <div ref={refChart} />
            <div ref={ref} style={{height:'400px'}}>{}</div>
        </div>
    );
}