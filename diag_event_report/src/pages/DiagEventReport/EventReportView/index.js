import {useSelector} from 'react-redux';
import MultiSubChart from './MultiSubChart';
import MultiYAxisChart from './MultiYAxisChart';

export default function EventReportView({chartType}){
    const dataList=useSelector(state=>state.reportList.list.data);
    console.log(dataList);

    //数据收敛
    const signalList=[];
    dataList.forEach(dataItem => {
        dataItem.Records.forEach(recItem=>{
            signalList.push(...recItem.SignalList);
        });
    });

    if(dataList.length<=0){
        return null;
    }

    return (
        <div style={{width:'100%',height:'400px',overflow:"auto"}}>
            {chartType===0?
            <MultiYAxisChart signalList={signalList} />:
            <MultiSubChart signalList={signalList}/>
            }
        </div>
    );
}