import {useCallback} from 'react';
import { Space,message,Button } from "antd";
import { useMemo } from "react";
import { useSelector,useDispatch } from "react-redux";

import {FRAME_MESSAGE_TYPE} from '../../../utils/constant';
import OperationButton from '../../../components/OperationButton';
import useI18n from '../../../hooks/useI18n';

import {
    getListOperationPreporcessFunc
} from '../../../utils/functions';

import './index.css';

export default function ListOperationBar({sendMessageToParent}){
    const {currentView} = useSelector(state=>state.data);
    const {getLocaleLabel}=useI18n();
    const {fields,views,modelID,operations}=useSelector(state=>state.definition);
    const {selectedRowKeys,filter,pagination,sorter,selectAll}=useSelector(state=>state.data.views[state.data.currentView].data);

    const {searchFields,filterData,viewFilter}=useMemo(()=>{
        let searchFields=[];
        const viewConf=views.find(item=>item.viewID===currentView);
        if(viewConf&&viewConf.fields){
            viewConf.fields.forEach((fieldItem,index) => {
                const fieldConf=fields.find(item=>item.field===fieldItem.field);
                if(fieldConf){
                    const searchField={
                        field:fieldItem.field,
                        dataType:fieldConf.dataType,
                        fieldType:fieldConf.fieldType,
                        relatedModelID:fieldConf.relatedModelID,
                        relatedField:fieldConf.relatedField,
                        associationModelID:fieldConf.associationModelID,
                        fields:fieldItem.fields
                    }
                    searchFields.push(searchField);
                }
            });
        }

        return {searchFields,filterData:viewConf?.filterData,viewFilter:viewConf.filter}
    },[currentView]);

    const doOperation=useCallback((opItem)=>{
        if(opItem.selectedRows){
            if(opItem.selectedRows?.min>selectedRowKeys.length){
                message.info(getLocaleLabel(opItem.selectedRows.prompt));
                return;
            }

            if(opItem.selectedRows?.max<selectedRowKeys.length){
                message.info(getLocaleLabel(opItem.selectedRows.prompt));
                return;
            }
        }
        
        const operation=operations.find(element=>element.id===opItem.operationID);
        if(operation){
            let queryFilter=filter;
            if(viewFilter&&Object.keys(viewFilter).length>0){
                if(Object.keys(filter).length>0){
                    queryFilter={
                        'Op.and':[filter,viewFilter]
                    };
                } else {
                    queryFilter=viewFilter;
                }
            }

            //由于行的ID可能是一个引用字段，所以这里需要对selectedRowKeys做一个检查和变换
            const selectedRowIDs=selectedRowKeys.map(item=>{
                if(item.value){
                    return item.value;
                }
                return item;
            });

            const input={
                modelID:modelID,
                viewID:currentView,
                selectedRowKeys:selectedRowIDs,
                filter:queryFilter,
                filterData:filterData,
                pagination:pagination,
                sorter:sorter,
                fields:searchFields,
                selectedAll:selectAll
            };

            //对operation做预处理，一般是基于数据行为operaiton增加过滤条件
            if(operation&&opItem.preprocessing){
                //console.log('preprocessing',opItem.preprocessing);
                operation=getListOperationPreporcessFunc(opItem.preprocessing)(operation,input);
                //console.log('preprocessing',operation);
            }

            const message={
                type:FRAME_MESSAGE_TYPE.DO_OPERATION,
                data:{
                    operationItem:{
                        ...operation,
                        input:{...input,...operation.input}
                    }
                }
            };
            
            sendMessageToParent(message);
        }
    },[operations,currentView,modelID,getLocaleLabel,viewFilter,selectAll,selectedRowKeys,filter,filterData,pagination,sorter,searchFields,sendMessageToParent]);

    const buttonControls=useMemo(()=>{
        let buttonControls=[];
        const viewConf=views.find(item=>item.viewID===currentView);
        if(viewConf){
            const listToolbar=viewConf.toolbar?.listToolbar;
            if(listToolbar){
                const {showCount,buttons}=listToolbar;
                if(buttons){
                    for(let i=0;i<buttons.length&&i<showCount;++i){
                        const item=buttons[i];
                        const operation=operations.find(element=>element.id===item.operationID);
                        if(operation){
                            buttonControls.push(
                                <OperationButton key={item.operationID} type='primary' doOperation={doOperation} operation={{name:operation.name,...item}}/>
                            );
                        }
                    }
                }
            }
        }
        return buttonControls;
    },[currentView,operations,doOperation]);

    return (
        <div className="list-operation-bar">
            <Space >
                {buttonControls}
            </Space>
        </div>
    );
}