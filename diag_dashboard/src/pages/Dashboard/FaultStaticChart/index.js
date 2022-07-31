import { useEffect,useRef,useMemo } from 'react';
import { useResizeDetector } from 'react-resize-detector';
import * as echarts from 'echarts';

export default function FaultStaticChart({data}){
    const refChart=useRef();
    const { width,height,ref } = useResizeDetector();

    console.log(data)

    const tmpData=useMemo(()=>{
        return [
        {type:"EPS",count:data.epsCount},
        {type:"ESC",count:data.escCount},
        {type:"IBS",count:data.ibsCount},
    ]},[data]);

    console.log(tmpData)
    
    const option = useMemo(()=>{
        return {
            backgroundColor:'',
            title: {
                text: '故障类型',
                left:'center',
                top:0,
                textStyle:{
                    fontStyle:'italic'
                }
            },
            tooltip: {
                trigger: 'item'
            },
            legend: {
                orient:'vertical',
                right:5,
                top:'center',
                data: tmpData.map(item=>{
                    return item.type
                })
            },
            series:{
                type: 'pie',
                radius: ['50%', '90%'],
                left:0,
                top:30,
                right:80,
                bottom:15,
                label: {
                    alignTo: 'edge',
                    formatter: '{b}:{c}',
                    minMargin: 5,
                    edgeDistance: 10,
                    lineHeight: 15,
                    rich: {
                    time: {
                        fontSize: 10,
                        color: '#999'
                    }
                    }
                },
                labelLine: {
                    length: 15,
                    length2: 0,
                    maxSurfaceAngle: 80
                },
                data: tmpData.map(item=>{
                    return {value:item.count,name:item.type};
                })
            },
        }
    },[tmpData]);

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