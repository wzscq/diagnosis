import { createSlice } from '@reduxjs/toolkit';
import { message } from 'antd';

import {getReport} from '../api';

// Define the initial state using that type
const initialState = {
    loaded:false,
    filter:{
        vin:"",
        code:"",
        startDate:"",
        endDate:""
    },
    list:{
        page:0,
        count:10,
        total:0,
        data:[]
    },
    pending:false
}

export const reportListSlice = createSlice({
    name: 'reportList',
    initialState,
    reducers: {
        setFilter: (state,action) => {
          state.filter=action.payload;
        },
        setLoad:(state,action) => {
            state.loaded=action.payload;
        },
    },
    extraReducers: (builder) => {
        builder.addCase(getReport.pending, (state, action) => {
            state.pending=true;
            state.loaded=true;
        });
        builder.addCase(getReport.fulfilled, (state, action) => {
            console.log('getReport fulfilled',action);
            state.list.total=action.payload.total;
            state.list.data=action.payload.data;
            state.pending=false;
        });
        builder.addCase(getReport.rejected , (state, action) => {
            state.pending=false;
            message.error("获取报告出错！");
        });
    }
});

// Action creators are generated for each case reducer function
export const { setFilter,setLoad} = reportListSlice.actions

export default reportListSlice.reducer