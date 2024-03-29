import { Col, Row,ConfigProvider } from 'antd';
import { useEffect } from 'react';
import { useSelector,useDispatch } from 'react-redux';
import { useParams } from 'react-router-dom';
import zh_CN from 'antd/lib/locale/zh_CN';
import en_US from 'antd/lib/locale/en_US';


import useFrame from '../../hooks/useFrame';
import ListOperationBar from './ListOperationBar';
import ListTable from './ListTable';
import ModelViewList from './ModelViewList';

import PageLoading from './PageLoading';
import SearchBar from './SearchBar';
import StatusBar from './StatusBar';

import {createGetModelConfMessage} from '../../utils/normalOperations';
import {initDataView} from '../../redux/dataSlice';
import NoView from './NoView';
import ColumnSettingDialog from './ColumnSettingDialog';
import useI18n from '../../hooks/useI18n';

import './index.css';

const locales={
    zh_CN:zh_CN,
    en_US:en_US
}

const theme={
    token: {
        colorBgSpotlight: '#000000',
        colorInfoHover:'#000000',
    },
}

export default function View(){
    const dispatch= useDispatch();
    const {locale}=useI18n();
    const {loaded,views} = useSelector(state=>state.definition);
    const {initialized} = useSelector(state=>state.data);
    const {origin,item}=useSelector(state=>state.frame);
    const {showColumnSettingDialog}=useSelector(state=>state.definition);
    const sendMessageToParent=useFrame();
    const {modelID}=useParams();

    useEffect(()=>{
        if(origin&&item){
            if(loaded===false){
                //加载配置
                console.log('get model config ...');
                sendMessageToParent(createGetModelConfMessage({frameType:item.frameType,frameID:item.params.key,origin:origin},modelID,item.params.views));
            } else if (initialized===false) {
                console.log("loaded views :",views);
                dispatch(initDataView({views,currentView:item.params.view,filter:item.params.filter}));
            }
        }
    },[loaded,origin,item,modelID,initialized,dispatch,sendMessageToParent,views]);

    if(loaded&&initialized){
        if(views?.length>0){            
            return (
                <ConfigProvider locale={locales[locale]} theme={theme}>
                    <div className='list_view_main'>
                        <Row>
                            <Col span={0}><ModelViewList/></Col>
                            <Col span={12}><SearchBar/></Col>
                            <Col span={12}><ListOperationBar sendMessageToParent={sendMessageToParent}/></Col>
                        </Row>
                        <Row>
                            <Col ><StatusBar/></Col>
                        </Row>
                        <Row>
                            <Col span={24}><ListTable sendMessageToParent={sendMessageToParent} /></Col>                   
                        </Row>
                    </div>
                    {showColumnSettingDialog===true?<ColumnSettingDialog/>:null}
                </ConfigProvider>
            );
        } else {
            return(<NoView/>);
        }
    } else {
        return(<PageLoading/>);
    }
}