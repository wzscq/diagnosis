import {Row,Col} from 'antd';

import SignalChart from "./SignalChart";
import { setAnalysisItem } from '../../../redux/reportSlice';
import { useDispatch } from 'react-redux';
import { useEffect } from 'react';

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
    wordBreak:"break-all"
}

const styleLabelNoBorder={
    borderLeft:"1px solid #888",
    textAlign:"right",
    paddingRight:20,
    lineHeight:"45px",
    wordBreak:"break-all"
}

const styleContentNoBorder={
    borderRight:"1px solid #888",
    wordBreak:"break-all"
}

const styleOtherItemContent={
    borderRight:"1px solid #888",
    wordBreak:"break-all"
}

const styleOtherItemLabel={
    borderLeft:"1px solid #888",
    textAlign:"right",
    paddingRight:20,
    lineHeight:"45px",
    wordBreak:"break-all"
}

const styleOtherItemInner={
    width:"calc(100% - 10px)",
    borderRadius:5, 
    margin:5,
    padding:5,
    border:"1px Solid #EEEEEE",
    wordBreak:"break-all"
}

const styleSignalChart={
    width:"calc(100%)",
    borderRadius:5, 
    margin:0,
    border:"0px Solid #EEEEEE",
    overflow:"hidden"
}

const styleSignalAnalysis={
    width:"calc(100% - 10px)",
    borderRadius:5, 
    margin:5,
    padding:5,
    border:"1px Solid #EEEEEE",
    wordBreak:"break-all"
}

export default function AnalysisItem({itmeIndex,report,rec,vin}){
    const dispatch=useDispatch();

    const signalAnalysis=[(<p style={{textIndent:"20px",padding:"10px"}}>{rec.FurtherCheckedSignalsAnalysis}</p>)];

    let signalList=rec.SignalList;
   
    const signalChart=signalList.map((item,index)=>{
        return <SignalChart itmeIndex={itmeIndex} chartIndex={index+''} signal={item}/>
    });

    const index=0;
    const OtherItems=[
        (
            <Row>
                <Col span={2} />
                <Col span={4} style={styleOtherItemLabel}>{index===0?"故障原因:":""}</Col>
                <Col span={7} >
                    <div style={styleOtherItemInner}>
                        {rec.PossibleCauses}
                    </div>
                </Col>
                <Col span={2} style={{...styleOtherItemLabel,borderLeft:0}}>{index===0?"修复建议:":""}</Col>
                <Col span={7} style={styleOtherItemContent}>
                    <div style={styleOtherItemInner}>
                        {rec.RecommendedRecovery}
                    </div>
                </Col>
                <Col span={2} />
            </Row>
        )];

    useEffect(()=>{
        dispatch(setAnalysisItem({itmeIndex:itmeIndex,item:rec}));
    });

    const converTime=(time)=>{
        if(time.length<14){
            return time;
        }
        return time.substr(0,4)+'-'+time.substr(4,2)+'-'+time.substr(6,2)+' '+time.substr(8,2)+':'+time.substr(10,2)+':'+time.substr(12,2);
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
                <Col span={4} style={styleLabelNoBorder}>故障内容:</Col>
                <Col span={16} style={styleContentNoBorder}>
                    <div style={{width:"calc(100% - 10px)",borderRadius:5, margin:5,padding:5,border:"1px Solid #EEEEEE"}}>
                    {rec?.DtcDescription}
                    </div>
                </Col>
                <Col span={2} />
            </Row>
            {OtherItems}
            <Row>
                <Col span={2} />
                <Col span={4} style={{...styleOtherItemLabel,}}>{"相关信号:"}</Col>
                <Col span={7} >
                    <div style={styleSignalChart}>
                        {signalChart}
                    </div>
                </Col>
                <Col span={2} style={{...styleOtherItemLabel,borderLeft:0}}>{"信号分析:"}</Col>
                <Col span={7} style={styleOtherItemContent}>
                    <div style={styleSignalAnalysis}>
                        {signalAnalysis}
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