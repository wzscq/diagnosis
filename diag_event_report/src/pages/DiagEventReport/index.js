import { useEffect,useState } from 'react';
import {useDispatch,useSelector} from 'react-redux';
import { Button,PageHeader } from 'antd';
import {DatabaseFilled,ProfileFilled} from '@ant-design/icons';
import useFrame from '../../hook/useFrame';
import PageLoading from './PageLoading';
import EventReportView from './EventReportView';
import {getEventReport} from '../../api';
import {
    OP_TYPE,
    OPEN_LOCATION,
    FRAME_MESSAGE_TYPE} from '../../utils/constant';
import {
    createQueryDataMessage
} from '../../utils/normalOperations';

const activeButton={margin:2,border:"1px solid #40a9ff",color:"#40a9ff"};
const normalButton={margin:2};

export default function DiagEventReport(){
    const sendMessageToParent=useFrame();
    const dispatch=useDispatch();
    const {origin,item}=useSelector(state=>state.frame);
    const dataLoaded=useSelector(state=>state.data.loaded);  
    const dataList=useSelector(state=>state.data.data?.list);  
    const reportLoaded=useSelector(state=>state.reportList.loaded);
    const [chartType,setChartType]=useState(0);
    //console.log(dataList);
    //加载数据
    useEffect(()=>{
        if(dataLoaded===false&&origin&&item){
            console.log('sendMessageToParent',origin,item);
            //目前的表单页面仅支持单条数据的编辑和展示
            const dataID=item?.input?.selectedRowKeys[0];
            if(dataID){
                const frameParams={
                    frameType:item.frameType,
                    frameID:item.params.key,
                    origin:origin
                };
                const queryParams={
                    modelID:item.input.modelID,
                    filter:{id:dataID},
                    fields:[{field:"id"},{field:"event_name"},{field:"vin"},{field:"event_report"}],
                    pagination:{current:1,pageSize:1}
                };
                sendMessageToParent(createQueryDataMessage(frameParams,queryParams));
            }
        }
    },[dataLoaded,item,origin,sendMessageToParent]);

    //加载报告
    useEffect(()=>{
        if(dataLoaded===true&&reportLoaded===false){
            if(dataList.length>0){
                const id=dataList[0].event_report;
                const params=id.split('_');
                console.log(params);
                const searchParams={
                    page:0,
                    count:1,
                    collection:params[0],
                    dtc:params[1],
                    startDate:params[2],
                    endDate:params[2]
                }
                dispatch(getEventReport(searchParams));
            }
        }
    },[dataLoaded,reportLoaded,dataList,dispatch]);

    const closeReport=()=>{
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

    const subTitle=dataList?.length>0?
        '场景：'+dataList[0].event_name+'; VIN:'+dataList[0].vin:'';

    return (
        <PageHeader
            className="site-page-header-responsive"
            title="场景分析"
            subTitle={subTitle}
            extra={[
                <Button size='small' onClick={()=>setChartType(0)} style={chartType===0?activeButton:normalButton}><ProfileFilled /></Button>,
                <Button size='small' onClick={()=>setChartType(1)} style={chartType===0?normalButton:activeButton}><DatabaseFilled /></Button>,
                <Button type='primary' onClick={closeReport} key="close">关闭</Button>
            ]}
        >
            {
                (dataLoaded&&reportLoaded)?
                <EventReportView chartType={chartType} sendMessageToParent={sendMessageToParent}/>:
                <PageLoading/>
            }
        </PageHeader>
    );
}