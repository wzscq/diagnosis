import {Row,Col} from 'antd';
import LogisticsInfo from './LogisticsInfo';
import SpecialLogisticsInfo from './SpecialLogisticsInfo';

import {setECUList} from '../../../redux/reportSlice';
import { useDispatch } from 'react-redux';
import { useEffect } from 'react';

const styleTitle={
    textAlign:"left",
    paddingLeft:20,
    fontSize:18,
    borderBottom:"1px solid #888"
}

const styleContent={
    textAlign:"left",
    paddingLeft:20,
    borderRight:"1px solid #888",
    borderBottom:"1px solid #888",
    wordBreak:"break-all"
}

const styleContentSN={
    textAlign:"center",
    borderRight:"1px solid #888",
    borderBottom:"1px solid #888",
    borderLeft:"1px solid #888"
}

function RecordsHeader({title}){
    return (
        <>
        <Row>
            <Col span={2} />
            <Col span={1} style={styleContentSN} >序号</Col>
            <Col span={3} style={styleContent} >故障代码</Col>
            <Col span={10} style={styleContent} >故障内容</Col>
            <Col span={3} style={styleContent} >故障时刻</Col>
            <Col span={3} style={styleContent} >车辆里程</Col>
            <Col span={2} />
        </Row>    
        </>
    );
}

const converTime=(time)=>{
    if(time.length<14){
        return time;
    }
    return time.substr(0,4)+'-'+time.substr(4,2)+'-'+time.substr(6,2)+' '+time.substr(8,2)+':'+time.substr(10,2)+':'+time.substr(12,2);
}

const getRecord=(rec,index)=>{
    return (
    <Row>
        <Col span={2} />
        <Col span={1} style={styleContentSN} >{index+1}</Col>
        <Col span={3} style={{...styleContent,color:rec.DtcId_State==="0"?"blue":"red"}} >{rec.DtcId}</Col>
        <Col span={10} style={styleContent} >{rec.DtcDescription}</Col>
        <Col span={3} style={styleContent} >{converTime(rec.Time)}</Col>
        <Col span={3} style={styleContent} >{rec.Mileage}</Col>
        <Col span={2} />
    </Row>);
}

/*const getEmptyRecord=(index)=>{
    return (
    <Row>
        <Col span={2} />
        <Col span={1} style={styleContentSN} >{index+1}</Col>
        <Col span={3} style={styleContent} ></Col>
        <Col span={10} style={styleContent} ></Col>
        <Col span={3} style={styleContent} ></Col>
        <Col span={3} style={styleContent} ></Col>
        <Col span={2} />
    </Row>);
}*/

export default function Records({report,vin}){
    const dispatch=useDispatch();

    const {Records,LogisticsInfo:logisticsInfo,SpecialLogisticsInfo:specialLogisticsInfo}=report;
    const ecuList={}
    Records.forEach(rec=>{
        if(!ecuList[rec.Ecu]){
            ecuList[rec.Ecu]={name:rec.Ecu,items:[],specialLogisticsInfo:specialLogisticsInfo};
        }
        ecuList[rec.Ecu].items.push(rec);
    });

    logisticsInfo.forEach(rec=>{
        if(!ecuList[rec.EcuName]){
            ecuList[rec.EcuName]={name:rec.EcuName,items:[]};
        }
        ecuList[rec.EcuName].logistics=rec;
    });

    console.log("ecuList",ecuList);

    const ecuControlList=[];
    for(let key in ecuList) {
        const ecuItem = ecuList[key];
        const ecuRecords=ecuItem.items.map((rec,index)=>getRecord(rec,index));
        ecuControlList.push(
            <>
                <LogisticsInfo title={ecuItem.name} logistics={ecuItem.logistics} vin={vin} />
                <SpecialLogisticsInfo title={ecuItem.name} logistics={ecuItem.specialLogisticsInfo} vin={vin} />
                <RecordsHeader title={ecuItem.name}/>
                {ecuRecords}
                <div style={{height:20,width:"100%"}}/>
            </>
        );
    }

    useEffect(()=>{
        dispatch(setECUList(ecuList));
    });

    //const escRecords=Records.filter(rec=>rec.Ecu==="ESC").map((rec,index)=>getRecord(rec,index));
    //const epsRecords=Records.filter(rec=>rec.Ecu==="EPS").map((rec,index)=>getRecord(rec,index));
    //const eBoosterRecords=Records.filter(rec=>rec.Ecu==="eBooster").map((rec,index)=>getRecord(rec,index));

    return (
        <div style={{width:"100%",border:"0px solid #EEEEEE",marginTop:0}}>
            <Row>
                <Col span={2} />
                <Col span={20} style={styleTitle}>2、故障码概览</Col>
                <Col span={2} />
            </Row>
            {ecuControlList}
        </div>
    );
}