import {Button,Row,Col} from 'antd';
import {DatabaseFilled,ProfileFilled} from '@ant-design/icons';
import SignalChart from "./SignalChartNew";
import SignalChartSplit from './SignalChartSplit';
import { setAnalysisItem } from '../../../redux/reportSlice';
import { useDispatch } from 'react-redux';
import { useEffect, useState } from 'react';

const styleLabel={
    borderRight:"1px solid #888",
    borderLeft:"1px solid #888",
    borderBottom:"1px solid #888",
    textAlign:"left",
    paddingLeft:20,
    backgroundColor:"#BBBBBB",
    wordBreak:"break-all"
}

const styleContent={
    borderRight:"1px solid #888",
    borderBottom:"1px solid #888",
    paddingLeft:20,
    wordBreak:"break-all",
    textAlign:"left",
}

const styleSignalChart={
    width:"calc(100%)",
    borderRadius:5, 
    margin:0,
    paddingLeft:"5px",
    paddingRight:"5px",
    border:"0px Solid #EEEEEE",
    overflow:"hidden"
}

const styleSignalTitle={
    borderLeft:"1px solid #888",
    textAlign:'left',
    paddingLeft:'20px'
}

const styleSignalToolbar={
    borderRight:"1px solid #888",
    textAlign:'left'
}

const styleSignalChartCell={
    borderLeft:"1px solid #888",
    borderRight:"1px solid #888"
}

const activeButton={margin:2,border:"1px solid #40a9ff",color:"#40a9ff"};
const normalButton={margin:2};

export default function AnalysisItem({itmeIndex,report,rec,vin}){
    const dispatch=useDispatch();
    const [chartType,setChartType]=useState(0);
    let signalList=rec.SignalList;
   
    useEffect(()=>{
        dispatch(setAnalysisItem({itmeIndex:itmeIndex,item:rec}));
    });

    const converTime=(time)=>{
        if(time.length<14){
            return time;
        }
        return time.substr(0,4)+'-'+time.substr(4,2)+'-'+time.substr(6,2)+' '+time.substr(8,2)+':'+time.substr(10,2)+':'+time.substr(12,2);
    }

    let possibleCauses=rec?.PossibleCauses
    console.log('wzstest',possibleCauses);
    if(possibleCauses.length>0){
        possibleCauses=possibleCauses.split(/\n/g)
    } else {
        possibleCauses=[];
    }

    return (
        <div style={{width:"100%",height:"auto",margin:"auto",marginBottom:30}}>
            <Row>
                <Col span={2} />
                <Col span={4} style={{...styleLabel,borderTop:"1px solid #888"}}>故障代码</Col>
                <Col span={16} style={{...styleContent,borderTop:"1px solid #888"}}>{rec?.DtcId}</Col>
                <Col span={2} />
            </Row>
            <Row>
                <Col span={2} />
                <Col span={4} style={{...styleLabel,borderTop:"0px solid #888"}}>控制器</Col>
                <Col span={16} style={{...styleContent,borderTop:"0px solid #888"}}>{rec?.Ecu}</Col>
                <Col span={2} />
            </Row>
            <Row>
                <Col span={2} />
                <Col span={4} style={styleLabel}>故障时刻</Col>
                <Col span={16} style={styleContent}>{converTime(rec.Time)}</Col>
                <Col span={2} />
            </Row>
            <Row>
                <Col span={2} />
                <Col span={4} style={styleLabel}>车辆里程</Col>
                <Col span={16} style={styleContent}>{rec?.Mileage}</Col>
                <Col span={2} />
            </Row>
            <Row>
                <Col span={2} />
                <Col span={4} style={styleLabel}>故障状态</Col>
                <Col span={16} style={styleContent}>{rec?.DtcId_State}</Col>
                <Col span={2} />
            </Row>
            <Row>
                <Col span={2} />
                <Col span={4} style={styleLabel}>故障内容</Col>
                <Col span={16} style={styleContent}>{
                    rec?.DtcDescription
                }</Col>
                <Col span={2} />
            </Row>
            <Row>
                <Col span={2} />
                <Col span={4} style={styleLabel}>故障原因</Col>
                <Col span={16} style={styleContent}>
                  {possibleCauses?.map(item=>(<>{item}<br/></>))}  
                </Col>
                <Col span={2} />
            </Row>
            <Row>
                <Col span={2} />
                <Col span={4} style={styleLabel}>修复建议</Col>
                <Col span={16} style={styleContent}>{rec?.RecommendedRecovery}</Col>
                <Col span={2} />
            </Row>
            {rec?.remark!==undefined&&rec?.remark!==""?
                (<Row>
                    <Col span={2} />
                    <Col span={4} style={styleLabel}>备注</Col>
                    <Col span={16} style={styleContent}>{rec?.remark}</Col>
                    <Col span={2} />
                </Row>):null
            }
            
            <Row>
                <Col span={2} />
                <Col span={4} style={styleSignalTitle}>{"相关信号:"}</Col>
                <Col span={16}  style={styleSignalToolbar}>
                    <Button size='small' onClick={()=>setChartType(0)} style={chartType===0?activeButton:normalButton}><ProfileFilled /></Button>
                    <Button size='small' onClick={()=>setChartType(1)} style={chartType===0?normalButton:activeButton}><DatabaseFilled /></Button>
                </Col>
                <Col span={2} />
            </Row>
            <Row>
                <Col span={2} />
                <Col span={20} style={styleSignalChartCell} >
                    <div style={styleSignalChart}>
                        {signalList.length>0?(chartType===0?<SignalChart report={report} itmeIndex={itmeIndex} signalList={signalList}/>:<SignalChartSplit report={report} itmeIndex={itmeIndex} signalList={signalList}/>):null}
                    </div>
                </Col>
                <Col span={2} />
            </Row>               
            <Row>
                <Col span={2} />
                <Col span={20} style={{height:20,border:"1px solid #888",borderTop:"0px"}} />
                <Col span={2} />
            </Row>                        
        </div>
    )
}