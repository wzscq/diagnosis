import { useEffect,useRef,useMemo } from 'react';
import { useResizeDetector } from 'react-resize-detector';
import * as echarts from 'echarts';

export default function ProjectCarChart({data}){
    const refChart=useRef();
    const { width,height,ref } = useResizeDetector();
    console.log('ProjectCarChart',data);    
    const option = useMemo(()=>{
        return {
            backgroundColor:'',
            title: {
                text: '项目车辆数量',
                left:'center',
                textStyle:{
                    fontStyle:'italic'
                }
            },
            tooltip: {
                trigger: 'axis'
            },
            grid:{
                left:25,
                top:30,
                right:70,
                bottom:25,
            }, 
            xAxis: [
                {
                    type: 'category',
                    axisTick: { show: false },
                    data: ['车辆数量'],
                    axisLabel:{
                        show:false
                    }
                }
            ],
            yAxis: {
                type:'value',
                show:true,
                axisLine:{
                    show:true
                }
            },
            legend: {
                orient:'vertical',
                right:0,
                top:0,
                data: data.map(item=>{
                    return item.ProjectNum
                })
            },
            series: data.map(item=>({
                type:'bar',
                barGap:'100%',
                label:{
                    show:true,
                    position:'bottom',
                    formatter:item.ProjectNum
                },

                name:item.ProjectNum,
                data:[item.count],
            })),
        }
    },[data]);

    useEffect(()=>{
        if(refChart&&refChart.current){
            let chart=echarts.getInstanceByDom(refChart.current);        
            if(chart){
                chart.resize({width:width,height:height});
                chart.setOption(option);
            } else {
                chart=echarts.init(refChart.current,'dark',{
                    width: width,
                    height: height
                });
                chart.setOption(option);
            }
        }
    },[option,refChart,width,height]);

    return (
        <>
            <div ref={refChart} style={{marginTop:5,width:width,marginBottom:20,marginLeft:0,height:'100%'}}/>
            <div style={{width:'100%',height:'100%'}} ref={ref}>{}</div>
        </>
    )
}