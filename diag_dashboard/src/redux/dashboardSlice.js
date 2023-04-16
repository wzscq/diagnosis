import { createSlice } from '@reduxjs/toolkit';
import { message } from 'antd';

import {getDashboard} from '../api';

const initialState = {
    carCount:0,
    faultCountByStatus:{
        openCount:0,
        closedCount:0,
    },
    faultCountByType: [],
    faultList: [],
    projectcarCount:[],
    pending:false
}

export const dashboardSlice = createSlice({
    name: 'dashboard',
    initialState,
    reducers: {
       
    },
    extraReducers: (builder) => {
        builder.addCase(getDashboard.pending, (state, action) => {
            state.pending=true;
        });
        builder.addCase(getDashboard.fulfilled, (state, action) => {
            state.carCount=action.payload.carCount;
            state.faultCountByStatus=action.payload.faultCountByStatus;
            state.faultCountByType=action.payload.faultCountByType;
            state.faultList=action.payload.faultList;
            state.projectcarCount=action.payload.projectcarCount;
            state.pending=false;
        });
        builder.addCase(getDashboard.rejected , (state, action) => {
            state.pending=false;
            message.error("获取数据出错！");
        });
    }
});

// Action creators are generated for each case reducer function
//export const { setFilter} = reportListSlice.actions

export default dashboardSlice.reducer