import { useSelector } from "react-redux"; 
import {PageHeader} from 'antd';
import CarInfo from './CarInfo';
import Records from './Records';
import Analysis from './Analysis';
import DownloadButton from './DownloadButton';

import { 
    OP_TYPE,
    OPEN_LOCATION,
    FRAME_MESSAGE_TYPE } from "../../../utils/constant";

export default function ReportView({sendMessageToParent}){
    const {item}=useSelector(state=>state.frame);
    const report=useSelector(state=>state.reportList.list&&state.reportList.list.data?state.reportList.list.data[0]:undefined);
    const data=useSelector(state=>state.data.data&&state.data.data.list?state.data.data.list[0]:undefined);
    
    const goBack=()=>{
        const operation={
            type:OP_TYPE.CLOSE,
            params:{
                location:OPEN_LOCATION.MODAL
            }
        }
        const message={
            type:FRAME_MESSAGE_TYPE.DO_OPERATION,
            data:{
                operationItem:{
                    ...operation
                }
            }
        };
        sendMessageToParent(message);
    }

    if(data&&report){
        const id=data.diag_report;
        const vin=data.id;
        return (
            <div>
                <PageHeader 
                    title="返回"
                    onBack={goBack}
                    ghost={false}
                    extra={[
                        <DownloadButton/>
                    ]}
                />
                <div style={{paddingBottom:'20px', marginLeft:"auto",marginRight:"auto"}}>
                    <div style={{width:'100%',textAlign:'center',paddingTop:20,letterSpacing:5,fontSize:21,fontWeight:600}}>
                        {"诊断报告"}
                    </div>
                    <div style={{width:'100%',textAlign:'center',padding:10}}>
                        {id}
                    </div>
                    <CarInfo report={report} vin={vin}/>
                    <Records report={report} vin={vin} />
                    <Analysis report={report} vin={vin}/>
                </div>
                <div style={{width:"100%",height:50}}/>
            </div>
        )
    } else {
        const dataID=item?.input?.selectedRowKeys[0];
        const reportID=(data&&data.diag_report)?data.diag_report:dataID;
        return (
        <div>
            <PageHeader 
                title="返回"
                onBack={goBack}
                ghost={false}
            />
            <div style={{paddingBottom:'20px', marginLeft:"auto",marginRight:"auto"}}>
                <div style={{width:'100%',textAlign:'center',paddingTop:20,letterSpacing:5,fontSize:21,fontWeight:600}}>
                    {"诊断报告"}
                </div>
                <div style={{width:'100%',textAlign:'center',padding:10}}>
                    {reportID}
                </div>
            </div>
        </div>
        );       
    }
}