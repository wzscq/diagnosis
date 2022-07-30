import { useDispatch,useSelector } from 'react-redux';
import { useEffect } from 'react';
import {createQueryDataMessage} from '../../utils/normalOperations';
import PageLoading from './PageLoading';
import ReportView from './ReportView';
import { getReport } from '../../api';
import {setFileName} from '../../redux/reportSlice';
import useFrame from '../../hook/useFrame';

export default function DiagReport(){
    const sendMessageToParent=useFrame();
    const dispatch=useDispatch();
    const {origin,item}=useSelector(state=>state.frame);
    const dataLoaded=useSelector(state=>state.data.loaded);  
    const dataList=useSelector(state=>state.data.data?.list);  
    const reportLoaded=useSelector(state=>state.reportList.loaded);

    console.log('DiagReport',item,dataLoaded,reportLoaded)
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
                    fields:[{field:"id"},{field:"diag_report"},{field:"vin"}],
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
                const id=dataList[0].diag_report;
                const params=id.split('_');
                console.log(params);
                const searchParams={
                    page:0,
                    count:1,
                    collection:params[0],
                    startDate:params[1],
                    endDate:params[1]
                }
                dispatch(getReport(searchParams));
                dispatch(setFileName(id));
            }
        }
    },[dataLoaded,reportLoaded,dataList,dispatch]);
    return((dataLoaded&&reportLoaded)?<ReportView sendMessageToParent={sendMessageToParent}/>:<PageLoading/>);
}