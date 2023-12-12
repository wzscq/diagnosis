import {useSelector} from 'react-redux';
import MultiSubChart from './MultiSubChart';
import MultiYAxisChart from './MultiYAxisChart';

export default function EventReportView({chartType}){
    const dataList=useSelector(state=>state.reportList.list.data);
    console.log(dataList);

    //数据收敛
    const signalList=[];
    const remark=[];
    dataList.forEach(dataItem => {
        dataItem.Records.forEach(recItem=>{
            signalList.push(...recItem.SignalList);
            if(recItem.remark!==undefined&&recItem.remark!==""){
                remark.push(recItem.remark);
            }
        });
    });

    if(dataList.length<=0){
        return null;
    }

    return (
        <div style={{width:'100%',height:'400px',overflow:"auto"}}>
            {remark.length>0?<div style={{position:'absolute',left:'23px',top:'50px',textAlign:'left'}}>{remark[0]}</div>:null}
            {chartType===0?
            <MultiYAxisChart signalList={signalList} />:
            <MultiSubChart signalList={signalList}/>
            }
        </div>
    );
}