import { useRef,useEffect, useMemo } from 'react';
import * as echarts from 'echarts';
import { useResizeDetector } from 'react-resize-detector';
import { getColor } from '../../../utils/colorPalette';

export default function MultiYAxisChart({signalList}){
    const refChart=useRef();
    const { width,height,ref } = useResizeDetector();
    
    const legendData=signalList.map(item=>item.SignalName+(item.SignalUint?'('+item.SignalUint+')':''));
    
    let maxX=-1;
    let minX=-1;

    let mapAxisValue={};

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
                //mapAxisValue[x]=x;
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
    /*console.log(mapAxisValue);
    const mapLabelValues={};
    Object.keys(mapAxisValue).forEach(axisValue=>{
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
        mapLabelValues[axisValue]=values;
    });*/

    const grid= {
        left: (signalList.length*40+'px'),
        right: '10px'
    };
                
    const xAxis={
        type: 'value',
        min:Math.floor(minX),
        max:Math.ceil(maxX)
    };

    const yAxis=signalList.map((_,index)=>{
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
        //const values =mapLabelValues[axisValue];

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
    }

    useEffect(()=>{
        if(refChart&&refChart.current){
            
            let chart=echarts.getInstanceByDom(refChart.current);        
            if(chart){
                console.log("resize");
                chart.resize({width:width,height:height});     
            } else {
                chart=echarts.init(refChart.current,'',{
                    width: width,
                    height: height
                });
                console.log("setOption");
                chart.setOption(option);   
            }
        }
    });

    return (
        <div style={{width:'100%',height:'400px',overflow:"hidden"}}>
            <div ref={refChart} />
            <div ref={ref} style={{height:'400px'}}>{}</div>
        </div>
    );
}