import { useRef,useEffect,useMemo } from 'react';
import { useResizeDetector } from 'react-resize-detector';
import * as echarts from 'echarts';
import { useDispatch } from 'react-redux';
import { getColor } from '../../../utils/colorPalette';
import { setAnalysisChart } from '../../../redux/reportSlice';

export default function SignalChartSplit({report,itmeIndex,signalList}){
    const dispatch=useDispatch();
    const refChart=useRef();
    const { width,height,ref } = useResizeDetector();

    const {seriesData,yAxis,grid,xAxis}=useMemo(()=>{
        let maxX=-1;
        let minX=-1;
        const seriesData= signalList.map((item,index)=>{        
            return ({
                name:item.SignalName+(item.SignalUint?'('+item.SignalUint+')':''),
                type: 'line',
                smoth:true,
                yAxisIndex:index,
                xAxisIndex:index,
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

        const grid=signalList.map((item,index,list)=>{
            return {
                top:(100*index+40)+'px',
                height:'60px',
                left: '40px',
                right: '10px',
            }
        });

        const {xLimitValue}=report;
        if(xLimitValue){
            const xLimit=parseFloat(xLimitValue);
            if(maxX<xLimit){
                maxX=xLimit;
            }
        }

        const xAxis=signalList.map((item,index)=>{
            return {
                type: 'value',
                min:Math.floor(minX),
                max:Math.ceil(maxX),
                gridIndex: index,
                position:'bottom',
                offset:0
            }
        });

        const yAxis=signalList.map((item,index)=>{
            return {
                type: 'value',
                name: item.SignalName+(item.SignalUint?'('+item.SignalUint+')':''),
                nameGap:5,
                nameTextStyle:{
                    align:'left'
                },
                position: 'left',
                //offset: index*40,
                alignTicks: true,
                gridIndex: index,
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
        return {seriesData,yAxis,grid,xAxis};
    },[signalList]);

    const option = useMemo(()=>{
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

        return ({
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
        toolbox: {
            feature: {
                dataZoom: {
                    yAxisIndex: 'none'
                },
                saveAsImage: {}
            }
        },
        series: seriesData
    });},[seriesData,yAxis,grid,xAxis]);

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
    },[refChart,width,height,option,itmeIndex,dispatch]);

    return (
        <div style={{width:'100%',height:signalList.length*100+40,overflow:"hidden"}}>
            <div ref={refChart} />
            <div ref={ref} style={{height:signalList.length*100+40}}>{}</div>
        </div>
    );
}