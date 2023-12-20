import {useMemo,useCallback, useState} from 'react';
import { Input, Space, Button,Tooltip } from 'antd';
import { SettingOutlined } from '@ant-design/icons';
import { useDispatch,useSelector } from "react-redux";
import { refreshData,setFilter } from '../../../redux/dataSlice';
import {setShowColumnSettingDialog} from '../../../redux/definitionSlice';

import './index.css';
import useI18n from '../../../hooks/useI18n';

export default function SearchBar(){
    const {getLocaleLabel}=useI18n();
    const {fields,views}=useSelector(state=>state.definition);
    const {currentView} = useSelector(state=>state.data);
    const [searchText,setSearchText]=useState("");
    const dispatch=useDispatch();

    const {quickSearchFields,showColumnSettings}=useMemo(()=>{
        const viewConf=views.find(item=>item.viewID===currentView);
        const showColumnSettings=viewConf?.options?.showColumnSettings;
        let quickSearchFields=[];
        if(viewConf&&viewConf.fields){
            quickSearchFields= fields.filter(item=>
                item.quickSearch&&
                viewConf.fields.find(viewItem=>viewItem.field===item.field&&viewItem.visible!==false)
            ).map(field=>field.field)
        }
        return {quickSearchFields,showColumnSettings};
    },[fields,currentView]);

    const onSearch=useCallback(()=>{
        if(quickSearchFields.length>0){
            const fieldsFilter=quickSearchFields.map(element => {
                const tempFieldFilter={};
                tempFieldFilter[element]='%'+searchText+'%';
                return tempFieldFilter;
            });
            const op='Op.or';
            dispatch(setFilter({[op]:fieldsFilter}));
        } else {
            console.log('no quick search fields');
        }
    },[quickSearchFields,searchText,dispatch]);

    const reset=()=>{
        dispatch(setFilter({}));
    }

    const refresh=()=>{
        dispatch(refreshData());
    }

    const columnSettings=()=>{
        dispatch(setShowColumnSettingDialog(true));
    }

    return (
        <div className='search-bar'>
            <Space>
                {quickSearchFields.length>0?<>
                    <Input style={{width:"250px"}} value={searchText} placeholder={getLocaleLabel({key:'page.crvlistview.searchInputPlaceholder',default:'input search text'})} onChange={(e)=>{setSearchText(e.target.value)}}/>
                    <Button
                        type="primary"
                        loading={false}
                        onClick={onSearch}
                    >
                        {getLocaleLabel({key:'page.crvlistview.query',default:'查询'})}
                    </Button>
                    <Button
                        type="primary"
                        loading={false}
                        onClick={reset}
                    >
                        {getLocaleLabel({key:'page.crvlistview.resetFilter',default:'重置'})}
                    </Button>
                    <Button
                        type="primary"
                        loading={false}
                        onClick={refresh}
                    >
                        {getLocaleLabel({key:'page.crvlistview.refresh',default:'刷新'})}
                    </Button>    
                </>:null
                }
                {
                showColumnSettings===true?(
                <Tooltip title={getLocaleLabel({key:'page.crvlistview.column',default:'列设置'})}>
                    <Button
                        type="primary"
                        icon={<SettingOutlined />}
                        loading={false}
                        onClick={columnSettings}
                    />
                </Tooltip>):null
                }
            </Space>
        </div>
    )
}