import {Pagination} from 'antd';
import { useDispatch,useSelector } from 'react-redux';

import {setPagination} from '../../../redux/dataSlice';
import useI18n from '../../../hooks/useI18n';

import './TableFooter.css';

export default function TableFooter(){
    const {getLocaleLabel}=useI18n();
    const dispatch=useDispatch()
    const {total,selectedRowKeys,pagination} = useSelector(state=>state.data.views[state.data.currentView].data);
    
    const sumLabel=getLocaleLabel({key:'page.crvlistview.total',default:'共 '})+
                    total+
                    getLocaleLabel({key:'page.crvlistview.item',default:' 条'})+'，'+
                    getLocaleLabel({key:'page.crvlistview.selected',default:'选中 '})+
                    selectedRowKeys.length+
                    getLocaleLabel({key:'page.crvlistview.item',default:' 条'});

    const onPaginationChange=(page, pageSize)=>{
        if(pageSize!==pagination.pageSize){
            page=1;
        }
        dispatch(setPagination({...pagination,current:page,pageSize:pageSize}));
    }

    const itemRender = (_, type, originalElement) => {
        console.log('itemRender:',type);
        return originalElement;
    };

    return (
        <div className="list-table-footer">
            <div className='sum-label'>
                {sumLabel}
            </div>
            <div className="list-table-pagination">
                <Pagination onChange={onPaginationChange} size="small" itemRender={itemRender} {...pagination} total={total} showSizeChanger />
            </div>
        </div>
    );
}