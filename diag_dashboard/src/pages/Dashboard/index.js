import './index.css';

import GridCellWrapper from './GridCellWrapper';
import CarNumber from './CarNumber';
import ProjectCarChart from './ProjectCarChart';
import FaultStaticChart from './FaultStaticChart';
import FaultStatusChart from './FaultStatusChart';
import FaultList from './FaultList';
import { useEffect, useState } from 'react';
import {getDashboard} from '../../api';
import { useDispatch,useSelector } from 'react-redux';

export default function Dashboard(){
    const [refresh,setRefresh]=useState(0);
    const dispatch=useDispatch();
    const {carCount,faultCountByStatus,faultCountByType,faultList,projectcarCount}=useSelector(state=>state.dashboard);

    useEffect(()=>{
        dispatch(getDashboard());
        /*setTimeout(()=>{
            setRefresh(refresh+1);
        },5000);*/
    },[dispatch,refresh]);

    return (
        <div className="dashboard">
            <GridCellWrapper col={1} row={1} colSpan={1} rowSpan={1}>
                <CarNumber total={carCount}/>
            </GridCellWrapper>
            <GridCellWrapper col={2} row={1} colSpan={2} rowSpan={1}>
                <div className='dashboard-title'>智能故障专家库</div>
            </GridCellWrapper>
            <GridCellWrapper col={1} row={2} colSpan={3} rowSpan={2}>
                <FaultList data={faultList}/>                
            </GridCellWrapper>
            <GridCellWrapper col={1} row={4} colSpan={1} rowSpan={1}>
                <FaultStatusChart data={faultCountByStatus}/>
            </GridCellWrapper>
            <GridCellWrapper col={2} row={4} colSpan={1} rowSpan={1}>
                <FaultStaticChart data={faultCountByType}/>
            </GridCellWrapper>
            <GridCellWrapper col={3} row={4} colSpan={1} rowSpan={1}>
                <ProjectCarChart data={projectcarCount}/>
            </GridCellWrapper>
        </div>
    )
}