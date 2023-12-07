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
                text: '项目故障数分布',
                left:'center',
                textStyle:{
                    color:'#0F0F0F',
                }
            },
            tooltip: {
                trigger: 'axis'
            },
            grid:{
                left:50,
                top:25,
                right:5,
                bottom:40,
            }, 
            dataZoom: [
                {
                  type: 'slider',
                  startValue: 0,
                  endValue: 1,
                  height:15,
                  bottom:5
                }
            ],
            xAxis: [
                {
                    type: 'category',
                    axisTick: { show: true },
                    data: data?.map(item=>item.ProjectNum),
                    axisLabel:{
                        show:true,
                        color:'#000',
                        interval:0,
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
            series:[{
                type:'bar',
                barGap:'100%',
                label:{
                    show:true,
                    position:'insideTop',
                    formatter: '{c}',
                    color:'#000',
                    fontSize:12,
                    rotate:0
                },
                data:data?.map(item=>item.count),
            }],
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
        <div style={{width:'100%',height:'200px'}} ref={ref}>
            <div ref={refChart} style={{marginTop:5,width:width,marginBottom:20,marginLeft:0,height:'200px'}}/>
        </div>
    )
}