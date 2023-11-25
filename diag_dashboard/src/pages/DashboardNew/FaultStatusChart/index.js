import { useEffect,useRef,useMemo } from 'react';
import { useResizeDetector } from 'react-resize-detector';
import * as echarts from 'echarts';

export default function FaultStatusChart({data}){
    const refChart=useRef();
    const { width,height,ref } = useResizeDetector();

    const tmpData=useMemo(()=>{
        return [
        {type:"Open",count:data.openCount},
        {type:"Closed",count:data.closedCount},
    ]},[data]);
    
    const option = useMemo(()=>{
        return {
            color:['#ee6666', '#73c0de', '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc','#5470c6','#fac858'],
            backgroundColor:'',
            title: {
                text: '故障状态',
                left:'center',
                textStyle:{
                    color:'#0F0F0F',
                }
            },
            tooltip: {
                trigger: 'item'
            },
            series:{
                type: 'pie',
                radius: ['30%', '90%'],
                left:5,
                top:35,
                right:5,
                bottom:20,
                label: {
                    alignTo: 'edge',
                    formatter: '{b}:{c}',
                    minMargin: 5,
                    edgeDistance: 10,
                    lineHeight: 12,
                    fontSize:12,
                    color:'#000'
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
        <div style={{width:'100%',height:'200px'}} ref={ref}>
            <div ref={refChart} style={{marginTop:5,width:width,marginBottom:20,marginLeft:0,height:'200px'}}/>
        </div>
    )
}